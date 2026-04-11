//go:build cgo || (windows && !cgo)

package mbtilemap

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type tilemapObj struct {
	tw, th       int32 // logical map size in tiles
	tileW, tileH int32 // TMX tile pixel size
	drawW, drawH int32 // display size (SetTileSize)

	tilesetFirstGID int32
	tex             rl.Texture2D
	tsCols          int32 // tiles per row in tileset image

	tileLayers [][][]int32 // [layer][y][x] gid
	layerNames []string   // parallel to tileLayers (TMX layer name)

	// collision[y][x]: 0 = walkable; non-zero = category / system id (bitmasking is script-side).
	collision [][]uint8

	tmxDir string

	release heap.ReleaseOnce
}

func (o *tilemapObj) TypeName() string { return "Tilemap" }

func (o *tilemapObj) TypeTag() uint16 { return heap.TagTilemap }

func (o *tilemapObj) Free() {
	o.release.Do(func() { rl.UnloadTexture(o.tex) })
}

type tmxMap struct {
	Width      int `xml:"width,attr"`
	Height     int `xml:"height,attr"`
	TileWidth  int `xml:"tilewidth,attr"`
	TileHeight int `xml:"tileheight,attr"`
	Tilesets   []tmxTileset `xml:"tileset"`
	Layers     []tmxLayer   `xml:"layer"`
	ObjGroups  []tmxOG      `xml:"objectgroup"`
}

type tmxTileset struct {
	FirstGID int `xml:"firstgid,attr"`
	Image    struct {
		Source string `xml:"source,attr"`
		Width  int    `xml:"width,attr"`
		Height int    `xml:"height,attr"`
	} `xml:"image"`
	TileWidth  int `xml:"tilewidth,attr"`
	TileHeight int `xml:"tileheight,attr"`
}

type tmxLayer struct {
	Name  string `xml:"name,attr"`
	Width int    `xml:"width,attr"`
	Data  struct {
		Encoding string `xml:"encoding,attr"`
		Chars    string `xml:",chardata"`
	} `xml:"data"`
	Properties []tmxProp `xml:"properties>property"`
}

type tmxProp struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type tmxOG struct {
	Name    string     `xml:"name,attr"`
	Objects []tmxObj   `xml:"object"`
}

type tmxObj struct {
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Width  float64 `xml:"width,attr"`
	Height float64 `xml:"height,attr"`
}

func parseCSVInts(s string) ([]int32, error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, ",")
	out := make([]int32, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseInt(p, 10, 32)
		if err != nil {
			return nil, err
		}
		out = append(out, int32(n))
	}
	return out, nil
}

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("TILEMAP.*: heap not bound")
	}
	return nil
}

func argF(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func argInt32(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("TILEMAP.LOAD", "tilemap", m.tmLoad)
	reg.Register("TILEMAP.FREE", "tilemap", runtime.AdaptLegacy(m.tmFree))
	reg.Register("TILEMAP.SETTILESIZE", "tilemap", runtime.AdaptLegacy(m.tmSetTileSize))
	reg.Register("TILEMAP.DRAW", "tilemap", runtime.AdaptLegacy(m.tmDraw))
	reg.Register("TILEMAP.GETTILE", "tilemap", runtime.AdaptLegacy(m.tmGetTile))
	reg.Register("TILEMAP.SETTILE", "tilemap", runtime.AdaptLegacy(m.tmSetTile))
	reg.Register("TILEMAP.ISSOLID", "tilemap", runtime.AdaptLegacy(m.tmIsSolid))
	reg.Register("TILEMAP.WIDTH", "tilemap", runtime.AdaptLegacy(m.tmWidth))
	reg.Register("TILEMAP.HEIGHT", "tilemap", runtime.AdaptLegacy(m.tmHeight))
	reg.Register("TILEMAP.LAYERCOUNT", "tilemap", runtime.AdaptLegacy(m.tmLayerCount))
	reg.Register("TILEMAP.LAYERNAME", "tilemap", m.tmLayerName)
	reg.Register("TILEMAP.DRAWLAYER", "tilemap", runtime.AdaptLegacy(m.tmDrawLayer))
	reg.Register("TILEMAP.COLLISIONAT", "tilemap", runtime.AdaptLegacy(m.tmCollisionAt))
	reg.Register("TILEMAP.SETCOLLISION", "tilemap", runtime.AdaptLegacy(m.tmSetCollision))
	reg.Register("TILEMAP.MERGECOLLISIONLAYER", "tilemap", runtime.AdaptLegacy(m.tmMergeCollisionLayer))
	reg.Register("TILEMAP.ISSOLIDCATEGORY", "tilemap", runtime.AdaptLegacy(m.tmIsSolidCategory))
}

func (m *Module) Shutdown() {}

func (m *Module) tmLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TILEMAP.LOAD expects path")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return value.Nil, fmt.Errorf("TILEMAP.LOAD: %w", err)
	}
	var raw tmxMap
	if err := xml.Unmarshal(data, &raw); err != nil {
		return value.Nil, fmt.Errorf("TILEMAP.LOAD: parse TMX: %w", err)
	}
	if len(raw.Tilesets) < 1 {
		return value.Nil, fmt.Errorf("TILEMAP.LOAD: no tileset")
	}
	ts := raw.Tilesets[0]
	imgPath := ts.Image.Source
	if !filepath.IsAbs(imgPath) {
		imgPath = filepath.Join(filepath.Dir(path), imgPath)
	}
	tex := rl.LoadTexture(imgPath)
	tw := int32(ts.TileWidth)
	if tw <= 0 {
		tw = int32(raw.TileWidth)
	}
	th := int32(ts.TileHeight)
	if th <= 0 {
		th = int32(raw.TileHeight)
	}
	cols := int32(1)
	if tw > 0 && ts.Image.Width > 0 {
		cols = int32(ts.Image.Width) / tw
		if cols < 1 {
			cols = 1
		}
	}

	o := &tilemapObj{
		tw:              int32(raw.Width),
		th:              int32(raw.Height),
		tileW:           int32(raw.TileWidth),
		tileH:           int32(raw.TileHeight),
		drawW:           int32(raw.TileWidth),
		drawH:           int32(raw.TileHeight),
		tilesetFirstGID: int32(ts.FirstGID),
		tex:             tex,
		tsCols:          cols,
		tmxDir:          filepath.Dir(path),
		collision:       make([][]uint8, raw.Height),
	}
	for y := range o.collision {
		o.collision[y] = make([]uint8, raw.Width)
	}

	for _, L := range raw.Layers {
		if L.Data.Encoding != "" && strings.ToLower(L.Data.Encoding) != "csv" {
			continue
		}
		if strings.EqualFold(L.Name, "collision") || layerIsSolidMeta(&L) {
			cells, err := parseCSVInts(L.Data.Chars)
			if err != nil || len(cells) < raw.Width*raw.Height {
				continue
			}
			i := 0
			for y := 0; y < raw.Height; y++ {
				for x := 0; x < raw.Width; x++ {
					if i < len(cells) && cells[i] != 0 {
						o.collision[y][x] = 1
					}
					i++
				}
			}
			continue
		}
		cells, err := parseCSVInts(L.Data.Chars)
		if err != nil {
			rl.UnloadTexture(tex)
			return value.Nil, fmt.Errorf("TILEMAP.LOAD: layer %q csv: %w", L.Name, err)
		}
		if len(cells) < raw.Width*raw.Height {
			rl.UnloadTexture(tex)
			return value.Nil, fmt.Errorf("TILEMAP.LOAD: layer %q wrong cell count", L.Name)
		}
		grid := make([][]int32, raw.Height)
		i := 0
		for y := 0; y < raw.Height; y++ {
			row := make([]int32, raw.Width)
			for x := 0; x < raw.Width; x++ {
				row[x] = cells[i]
				i++
			}
			grid[y] = row
		}
		o.tileLayers = append(o.tileLayers, grid)
		o.layerNames = append(o.layerNames, L.Name)
	}
	if len(o.tileLayers) == 0 {
		rl.UnloadTexture(tex)
		return value.Nil, fmt.Errorf("TILEMAP.LOAD: no tile layer found")
	}

	// Object group collision (grid-aligned rects)
	for _, og := range raw.ObjGroups {
		if !strings.EqualFold(og.Name, "collision") {
			continue
		}
		for _, ob := range og.Objects {
			x0 := int(ob.X / float64(o.tileW))
			y0 := int(ob.Y / float64(o.tileH))
			x1 := int((ob.X + ob.Width) / float64(o.tileW))
			y1 := int((ob.Y + ob.Height) / float64(o.tileH))
			for yy := y0; yy <= y1; yy++ {
				for xx := x0; xx <= x1; xx++ {
					if yy >= 0 && yy < int(o.th) && xx >= 0 && xx < int(o.tw) {
						o.collision[yy][xx] = 1
					}
				}
			}
		}
	}

	id, err := m.h.Alloc(o)
	if err != nil {
		rl.UnloadTexture(tex)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func layerIsSolidMeta(L *tmxLayer) bool {
	for _, p := range L.Properties {
		if strings.EqualFold(p.Name, "solid") && (p.Value == "1" || strings.EqualFold(p.Value, "true")) {
			return true
		}
	}
	return false
}

func (m *Module) getTM(args []value.Value, ix int, op string) (*tilemapObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected tilemap handle", op)
	}
	return heap.Cast[*tilemapObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) tmFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TILEMAP.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) tmSetTileSize(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TILEMAP.SETTILESIZE expects (map, drawW, drawH)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.SETTILESIZE")
	if err != nil {
		return value.Nil, err
	}
	w, ok1 := argInt32(args[1])
	h, ok2 := argInt32(args[2])
	if !ok1 || !ok2 || w < 1 || h < 1 {
		return value.Nil, fmt.Errorf("TILEMAP.SETTILESIZE: sizes must be positive integers")
	}
	o.drawW, o.drawH = w, h
	return value.Nil, nil
}

func gidToSrc(o *tilemapObj, gid int32) (rl.Rectangle, bool) {
	if gid < o.tilesetFirstGID {
		return rl.Rectangle{}, false
	}
	lid := gid - o.tilesetFirstGID
	if lid < 0 {
		return rl.Rectangle{}, false
	}
	tx := lid % o.tsCols
	ty := lid / o.tsCols
	return rl.Rectangle{
		X:      float32(tx * o.tileW),
		Y:      float32(ty * o.tileH),
		Width:  float32(o.tileW),
		Height: float32(o.tileH),
	}, true
}

func (m *Module) tmDraw(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TILEMAP.DRAW expects tilemap handle")
	}
	o, err := m.getTM(args, 0, "TILEMAP.DRAW")
	if err != nil {
		return value.Nil, err
	}
	for li := range o.tileLayers {
		tmDrawLayerGrid(o, o.tileLayers[li])
	}
	return value.Nil, nil
}

func tmDrawLayerGrid(o *tilemapObj, layer [][]int32) {
	tint := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for y := int32(0); y < o.th; y++ {
		for x := int32(0); x < o.tw; x++ {
			gid := layer[y][x]
			if gid == 0 {
				continue
			}
			src, ok := gidToSrc(o, gid)
			if !ok {
				continue
			}
			dx := float32(x * o.drawW)
			dy := float32(y * o.drawH)
			rl.DrawTexturePro(o.tex, src, rl.Rectangle{X: dx, Y: dy, Width: float32(o.drawW), Height: float32(o.drawH)}, rl.Vector2{}, 0, tint)
		}
	}
}

func (m *Module) tmDrawLayer(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TILEMAP.DRAWLAYER expects (map, layerIndex)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.DRAWLAYER")
	if err != nil {
		return value.Nil, err
	}
	li, ok := argInt32(args[1])
	if !ok || int(li) < 0 || int(li) >= len(o.tileLayers) {
		return value.Nil, fmt.Errorf("TILEMAP.DRAWLAYER: invalid layer index")
	}
	tmDrawLayerGrid(o, o.tileLayers[li])
	return value.Nil, nil
}

func (m *Module) tmLayerCount(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TILEMAP.LAYERCOUNT expects map handle")
	}
	o, err := m.getTM(args, 0, "TILEMAP.LAYERCOUNT")
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(len(o.tileLayers))), nil
}

func (m *Module) tmLayerName(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TILEMAP.LAYERNAME expects (map, layerIndex)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.LAYERNAME")
	if err != nil {
		return value.Nil, err
	}
	li, ok := argInt32(args[1])
	if !ok || int(li) < 0 || int(li) >= len(o.layerNames) {
		return rt.RetString(""), nil
	}
	return rt.RetString(o.layerNames[li]), nil
}

func (m *Module) tmGetTile(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TILEMAP.GETTILE expects (map, layer, tx, ty)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.GETTILE")
	if err != nil {
		return value.Nil, err
	}
	li, okL := argInt32(args[1])
	x, okx := argInt32(args[2])
	y, oky := argInt32(args[3])
	if !okL || !okx || !oky {
		return value.Nil, fmt.Errorf("TILEMAP.GETTILE: layer, x, y must be numeric")
	}
	if int(li) < 0 || int(li) >= len(o.tileLayers) || y < 0 || x < 0 || y >= o.th || x >= o.tw {
		return value.FromInt(0), nil
	}
	return value.FromInt(int64(o.tileLayers[li][y][x])), nil
}

func (m *Module) tmSetTile(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("TILEMAP.SETTILE expects (map, layer, tx, ty, gid)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.SETTILE")
	if err != nil {
		return value.Nil, err
	}
	li, okL := argInt32(args[1])
	x, okx := argInt32(args[2])
	y, oky := argInt32(args[3])
	gid, okg := argInt32(args[4])
	if !okL || !okx || !oky || !okg {
		return value.Nil, fmt.Errorf("TILEMAP.SETTILE: numeric args required")
	}
	if int(li) < 0 || int(li) >= len(o.tileLayers) || y < 0 || x < 0 || y >= o.th || x >= o.tw {
		return value.Nil, nil
	}
	o.tileLayers[li][y][x] = gid
	return value.Nil, nil
}

func (m *Module) tmIsSolid(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TILEMAP.ISSOLID expects (map, tx, ty)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.ISSOLID")
	if err != nil {
		return value.Nil, err
	}
	x, okx := argInt32(args[1])
	y, oky := argInt32(args[2])
	if !okx || !oky {
		return value.Nil, fmt.Errorf("TILEMAP.ISSOLID: x, y must be numeric")
	}
	if y < 0 || x < 0 || y >= o.th || x >= o.tw {
		return value.FromBool(false), nil
	}
	return value.FromBool(o.collision[y][x] != 0), nil
}

func (m *Module) tmCollisionAt(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TILEMAP.COLLISIONAT expects (map, tx, ty)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.COLLISIONAT")
	if err != nil {
		return value.Nil, err
	}
	x, okx := argInt32(args[1])
	y, oky := argInt32(args[2])
	if !okx || !oky {
		return value.Nil, fmt.Errorf("TILEMAP.COLLISIONAT: x, y must be numeric")
	}
	if y < 0 || x < 0 || y >= o.th || x >= o.tw {
		return value.FromInt(0), nil
	}
	return value.FromInt(int64(o.collision[y][x])), nil
}

func (m *Module) tmSetCollision(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TILEMAP.SETCOLLISION expects (map, tx, ty, category)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.SETCOLLISION")
	if err != nil {
		return value.Nil, err
	}
	x, okx := argInt32(args[1])
	y, oky := argInt32(args[2])
	cat, okc := argInt32(args[3])
	if !okx || !oky || !okc {
		return value.Nil, fmt.Errorf("TILEMAP.SETCOLLISION: numeric args required")
	}
	if y < 0 || x < 0 || y >= o.th || x >= o.tw {
		return value.Nil, nil
	}
	if cat < 0 {
		cat = 0
	}
	if cat > 255 {
		cat = 255
	}
	o.collision[y][x] = uint8(cat)
	return value.Nil, nil
}

func (m *Module) tmMergeCollisionLayer(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TILEMAP.MERGECOLLISIONLAYER expects (map, layerIndex, category)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.MERGECOLLISIONLAYER")
	if err != nil {
		return value.Nil, err
	}
	li, okL := argInt32(args[1])
	cat, okc := argInt32(args[2])
	if !okL || !okc {
		return value.Nil, fmt.Errorf("TILEMAP.MERGECOLLISIONLAYER: layer and category must be numeric")
	}
	if cat < 1 || cat > 255 {
		return value.Nil, fmt.Errorf("TILEMAP.MERGECOLLISIONLAYER: category must be 1..255")
	}
	if int(li) < 0 || int(li) >= len(o.tileLayers) {
		return value.Nil, fmt.Errorf("TILEMAP.MERGECOLLISIONLAYER: invalid layer index")
	}
	layer := o.tileLayers[li]
	u8 := uint8(cat)
	for y := int32(0); y < o.th; y++ {
		for x := int32(0); x < o.tw; x++ {
			if layer[y][x] != 0 {
				o.collision[y][x] = u8
			}
		}
	}
	return value.Nil, nil
}

func (m *Module) tmIsSolidCategory(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TILEMAP.ISSOLIDCATEGORY expects (map, tx, ty, category)")
	}
	o, err := m.getTM(args, 0, "TILEMAP.ISSOLIDCATEGORY")
	if err != nil {
		return value.Nil, err
	}
	x, okx := argInt32(args[1])
	y, oky := argInt32(args[2])
	cat, okc := argInt32(args[3])
	if !okx || !oky || !okc {
		return value.Nil, fmt.Errorf("TILEMAP.ISSOLIDCATEGORY: numeric args required")
	}
	if y < 0 || x < 0 || y >= o.th || x >= o.tw {
		return value.FromBool(false), nil
	}
	c := o.collision[y][x]
	if cat == 0 {
		return value.FromBool(c != 0), nil
	}
	return value.FromBool(int32(c) == cat), nil
}

func (m *Module) tmWidth(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TILEMAP.WIDTH expects map handle")
	}
	o, err := m.getTM(args, 0, "TILEMAP.WIDTH")
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(o.tw)), nil
}

func (m *Module) tmHeight(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TILEMAP.HEIGHT expects map handle")
	}
	o, err := m.getTM(args, 0, "TILEMAP.HEIGHT")
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(o.th)), nil
}
