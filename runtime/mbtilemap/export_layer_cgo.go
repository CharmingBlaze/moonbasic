//go:build cgo || (windows && !cgo)

package mbtilemap

import (
	"fmt"

	"moonbasic/vm/heap"
)

// TileLayerGIDs exports one tile layer for terrain / editor bridges (read-only view of internal grid).
func TileLayerGIDs(store *heap.Store, h heap.Handle, layerIndex int) (tw, th int32, drawW, drawH int32, gids [][]int32, err error) {
	o, err := heap.Cast[*tilemapObj](store, h)
	if err != nil {
		return 0, 0, 0, 0, nil, err
	}
	if layerIndex < 0 || layerIndex >= len(o.tileLayers) {
		return 0, 0, 0, 0, nil, fmt.Errorf("tilemap: invalid layer index %d", layerIndex)
	}
	return o.tw, o.th, o.drawW, o.drawH, o.tileLayers[layerIndex], nil
}
