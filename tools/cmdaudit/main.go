// Cmdaudit compares compiler/builtinmanifest/commands.json to docs coverage and
// writes docs/COMMAND_AUDIT.md. Run from repo root: go run ./tools/cmdaudit
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type cmdRow struct {
	Key         string   `json:"key"`
	Args        []string `json:"args"`
	Returns     string   `json:"returns,omitempty"`
	Description string   `json:"description,omitempty"`
	Namespace   string   `json:"namespace,omitempty"`
}

type root struct {
	Commands []cmdRow `json:"commands"`
}

// namespaceRef maps dotted namespace (first segment) to primary reference path under docs/.
// (global) is documented across COMMANDS.md topic sections and LANGUAGE.md.
var namespaceRef = map[string]struct {
	Doc   string
	Blurb string
}{
	"(global)": {Doc: "COMMANDS.md", Blurb: "Console I/O, math, strings, arrays, files, bitwise, dates, and globals — see topic sections in COMMANDS.md and LANGUAGE.md."},
	"ANGLE":    {Doc: "reference/GAME_MATH_HELPERS.md", Blurb: "Angle wrap / difference helpers."},
	"ANIM":     {Doc: "reference/ENTITY.md", Blurb: "Entity animation clips, parameters, and transitions."},
	"ATLAS":    {Doc: "reference/ATLAS.md", Blurb: "Texture atlases for packed sprites."},
	"AUDIO":    {Doc: "reference/AUDIO.md", Blurb: "Audio device init, sounds, and music playback."},
	"AUDIOSTREAM": {Doc: "reference/AUDIOSTREAM.md", Blurb: "Low-level streaming audio buffers."},
	"BBOX":     {Doc: "reference/MODEL.md", Blurb: "Axis-aligned bounding boxes for models."},
	"BODY2D":   {Doc: "reference/PHYSICS2D.md", Blurb: "2D rigid bodies (Box2D)."},
	"BODY3D":   {Doc: "reference/PHYSICS3D.md", Blurb: "3D rigid bodies (Jolt where available)."},
	"BOX2D":    {Doc: "reference/PHYSICS2D.md", Blurb: "2D physics helpers (world-scale, queries)."},
	"BSPHERE":  {Doc: "reference/MODEL.md", Blurb: "Bounding spheres for culling and tests."},
	"BTREE":    {Doc: "reference/NAV_AI.md", Blurb: "Behavior-tree style AI helpers."},
	"CAMERA":   {Doc: "reference/CAMERA.md", Blurb: "3D perspective cameras and view/projection."},
	"CAMERA2D": {Doc: "reference/CAMERA.md", Blurb: "2D orthographic camera for scrolling."},
	"CHECK":    {Doc: "reference/COLLISION.md", Blurb: "Frustum / in-view checks (`CHECK.INVIEW`, …)."},
	"CHARCONTROLLER": {Doc: "reference/CHARCONTROLLER.md", Blurb: "Character controller vs 3D physics (Linux/Jolt)."},
	"CHUNK":    {Doc: "reference/TERRAIN.md", Blurb: "Terrain chunk streaming queries (`CHUNK.*`); see also WORLD.md."},
	"CLIENT":   {Doc: "reference/NETWORK.md", Blurb: "Network client connection API."},
	"CLIPBOARD": {Doc: "reference/IMAGE.md", Blurb: "Clipboard image/text (where supported)."},
	"CONTROLLER": {Doc: "reference/INPUT.md", Blurb: "Gamepad / controller device queries."},
	"COLOR":    {Doc: "reference/DRAW2D.md", Blurb: "Color constants and conversions for drawing."},
	"COMPUTESHADER": {Doc: "reference/SHADER.md", Blurb: "Compute shader dispatch and uniforms."},
	"CURSOR":   {Doc: "reference/INPUT.md", Blurb: "Mouse cursor shape and visibility."},
	"DATA":     {Doc: "reference/DATA.md", Blurb: "Structured game data assets."},
	"DEBUG":    {Doc: "reference/DEBUG.md", Blurb: "Runtime watches, logging, and overlays."},
	"DECAL":    {Doc: "reference/RENDER.md", Blurb: "Deferred decals on surfaces."},
	"DRAW":     {Doc: "reference/DRAW2D.md", Blurb: "2D draw aliases and helpers."},
	"DRAW3D":   {Doc: "reference/DRAW3D.md", Blurb: "3D primitives, billboards, and debug draws."},
	"DRAWPRIM2D": {Doc: "reference/DRAW2D.md", Blurb: "Low-level 2D primitive batch helpers."},
	"DRAWPRIM3D": {Doc: "reference/DRAW3D.md", Blurb: "Low-level 3D primitive batch helpers."},
	"FOG":      {Doc: "reference/WEATHER.md", Blurb: "Distance fog parameters (weather / atmosphere)."},
	"EFFECT":   {Doc: "reference/RENDER.md", Blurb: "Screen-space and post effects."},
	"ENEMY":    {Doc: "reference/NAV_AI.md", Blurb: "Enemy / wave helpers (where registered)."},
	"ENET":     {Doc: "reference/NETWORK.md", Blurb: "Legacy ENet-style stubs; prefer NET.*."},
	"EVENT":    {Doc: "reference/NETWORK.md", Blurb: "Multiplayer event and lobby routing."},
	"ENTITY":   {Doc: "reference/ENTITY.md", Blurb: "3D entities: create, transform, draw, animation, physics hooks."},
	"FILE":     {Doc: "reference/FILE.md", Blurb: "File and directory I/O beyond legacy globals."},
	"FONT":     {Doc: "reference/FONT.md", Blurb: "TTF/OTF font loading for text drawing."},
	"FREE":     {Doc: "reference/ENTITY.md", Blurb: "Legacy `FREE` / entity free aliases."},
	"GAME":     {Doc: "reference/GAMEHELPERS.md", Blurb: "Screen size, delta time, and game shortcuts."},
	"GESTURE":  {Doc: "reference/INPUT.md", Blurb: "Touch and gesture recognition."},
	"GRID":     {Doc: "reference/GRID.md", Blurb: "Grid overlays and helpers (`GRID.*`)."},
	"GUI":      {Doc: "reference/GUI.md", Blurb: "Immediate-mode UI (raygui or purego subset)."},
	"IMAGE":    {Doc: "reference/IMAGE.md", Blurb: "CPU images, pixels, and image file I/O."},
	"INPUT":    {Doc: "reference/INPUT.md", Blurb: "Keyboard, mouse, and gamepad input."},
	"INSTANCE": {Doc: "reference/INSTANCE.md", Blurb: "GPU instanced draws and instance buffers."},
	"JOINT2D":  {Doc: "reference/PHYSICS2D.md", Blurb: "2D physics joints (Box2D)."},
	"JOINT3D":  {Doc: "reference/PHYSICS3D.md", Blurb: "3D physics joints (Jolt where available)."},
	"JOLT":     {Doc: "reference/PHYSICS3D.md", Blurb: "Low-level Jolt queries and settings (Linux/CGO)."},
	"JSON":     {Doc: "reference/JSON.md", Blurb: "JSON parse, stringify, and DOM-style access."},
	"KEY":      {Doc: "reference/INPUT.md", Blurb: "Virtual key constants (`KEY.*`)."},
	"LEVEL":    {Doc: "reference/LEVEL.md", Blurb: "Level / scene load from glTF and markers."},
	"LIGHT":    {Doc: "reference/LIGHT.md", Blurb: "3D lights and shadows integration."},
	"LIGHT2D":  {Doc: "reference/RENDER.md", Blurb: "2D lighting passes and layers."},
	"LOBBY":    {Doc: "reference/NETWORK.md", Blurb: "Lobby and matchmaking helpers."},
	"MAT4":     {Doc: "reference/MAT4.md", Blurb: "4×4 matrices for 3D transforms."},
	"MATERIAL": {Doc: "reference/MODEL.md", Blurb: "PBR and default materials for meshes."},
	"MATH":     {Doc: "reference/MATH.md", Blurb: "Extended math beyond global builtins."},
	"MATRIX":   {Doc: "reference/TRANSFORM.md", Blurb: "Legacy matrix handle (see also MAT4)."},
	"MEM":      {Doc: "reference/MEM.md", Blurb: "Raw memory views and binary packing."},
	"MUSIC":    {Doc: "reference/AUDIO.md", Blurb: "Streaming music playback helpers."},
	"MESH":     {Doc: "reference/MESH.md", Blurb: "3D mesh geometry builders."},
	"MODEL":    {Doc: "reference/MODEL.md", Blurb: "3D models, drawing, and animation."},
	"MOUSE":    {Doc: "reference/INPUT.md", Blurb: "Legacy mouse aliases (`MOUSE.*`)."},
	"MOVE":     {Doc: "reference/MOVEMENT.md", Blurb: "Movement helpers (`MOVE.*`)."},
	"MOVER":    {Doc: "reference/MOVEMENT.md", Blurb: "Mover / kinematic helpers."},
	"NAV":      {Doc: "reference/NAV_AI.md", Blurb: "Navigation mesh baking and queries."},
	"NAVAGENT": {Doc: "reference/NAV_AI.md", Blurb: "Agents moving on nav meshes."},
	"NET":      {Doc: "reference/NETWORK.md", Blurb: "Networking (ENet-based) connections and channels."},
	"NOISE":    {Doc: "reference/NOISE.md", Blurb: "Procedural noise generators."},
	"PACKET":   {Doc: "reference/NETWORK.md", Blurb: "Binary packet read/write helpers."},
	"PARTICLE": {Doc: "reference/PARTICLE.md", Blurb: "3D GPU particle systems."},
	"PARTICLE2D": {Doc: "reference/PARTICLES.md", Blurb: "2D particle presets and emitters."},
	"PARTICLE3D": {Doc: "reference/PARTICLE.md", Blurb: "3D particle variants and batches."},
	"PARTICLES":  {Doc: "reference/PARTICLES.md", Blurb: "Legacy `PARTICLES.*` namespace."},
	"PATH":     {Doc: "reference/NAV_AI.md", Blurb: "Pathfinding paths and waypoints."},
	"PEER":     {Doc: "reference/NETWORK.md", Blurb: "Connected peer handles and send/receive."},
	"PICK":     {Doc: "reference/RAYCAST.md", Blurb: "3D pick / intersection queries (`PICK.*`)."},
	"PLAYER":   {Doc: "reference/PLAYER.md", Blurb: "First-person / player controller helpers."},
	"PLAYER2D": {Doc: "reference/PHYSICS2D.md", Blurb: "2D player / platformer helpers."},
	"PHYSICS":  {Doc: "reference/PHYSICS3D.md", Blurb: "Legacy `PHYSICS.*` world helpers (see PHYSICS2D/PHYSICS3D)."},
	"PHYSICS2D": {Doc: "reference/PHYSICS2D.md", Blurb: "2D physics world (Box2D): gravity, step, tuning."},
	"PHYSICS3D": {Doc: "reference/PHYSICS3D.md", Blurb: "3D physics world (Jolt on Linux/CGO)."},
	"POOL":     {Doc: "reference/POOL.md", Blurb: "Object pools for hot paths."},
	"POST":     {Doc: "reference/RENDER.md", Blurb: "Post-process passes."},
	"QUAT":     {Doc: "reference/VEC_QUAT.md", Blurb: "Quaternions for rotation."},
	"RAND":     {Doc: "reference/MATH.md", Blurb: "Random streams and distributions."},
	"RAY":      {Doc: "reference/DRAW3D.md", Blurb: "Picking, rays, and 3D intersection helpers."},
	"RAY2D":    {Doc: "reference/PHYSICS2D.md", Blurb: "2D ray casts against Box2D world."},
	"RAYLIB":   {Doc: "reference/RAYLIB_EXTRAS.md", Blurb: "Raylib misc utilities exposed to scripts."},
	"RENDER":   {Doc: "reference/RENDER.md", Blurb: "Frame lifecycle, clear, present, pipeline modes."},
	"RENDERTARGET": {Doc: "reference/RENDER.md", Blurb: "Render-to-texture targets."},
	"RES":      {Doc: "reference/CONFIG.md", Blurb: "Resource path helpers (`RES.PATH`, `RES.EXISTS`)."},
	"RPC":      {Doc: "reference/NETWORK.md", Blurb: "Remote procedure calls over network sessions."},
	"SAVE":     {Doc: "reference/FILE.md", Blurb: "Save-game / persistence helpers (`SAVE.*`)."},
	"SCENE":    {Doc: "reference/SCENE.md", Blurb: "Scene graph and entity helpers."},
	"SKY":      {Doc: "reference/SKY.md", Blurb: "Skybox / procedural sky rendering."},
	"SERVER":   {Doc: "reference/NETWORK.md", Blurb: "Game server bind and poll loop."},
	"SHADER":   {Doc: "reference/SHADER.md", Blurb: "Vertex/fragment shaders and uniforms."},
	"SOUND":    {Doc: "reference/AUDIO.md", Blurb: "Sound handle helpers."},
	"SPRITE":   {Doc: "reference/SPRITE.md", Blurb: "2D sprites, animation, and batches."},
	"SPRITEBATCH": {Doc: "reference/SPRITE.md", Blurb: "Batched sprite draws."},
	"SPRITEGROUP": {Doc: "reference/SPRITE.md", Blurb: "Sprite grouping for sorting."},
	"SPRITELAYER": {Doc: "reference/SPRITE.md", Blurb: "Layered sprite ordering."},
	"SPRITEUI": {Doc: "reference/SPRITE.md", Blurb: "Screen-space sprite UI helpers."},
	"STEER":    {Doc: "reference/NAV_AI.md", Blurb: "Steering behaviors for agents."},
	"STOPWATCH": {Doc: "reference/TIME.md", Blurb: "High-resolution timers (`STOPWATCH.*`)."},
	"STRING":   {Doc: "reference/STRING.md", Blurb: "String heap helpers (`STRING.*`)."},
	"SYSTEM":   {Doc: "reference/SYSTEM.md", Blurb: "Host OS, process, and environment."},
	"TABLE":    {Doc: "reference/TABLE.md", Blurb: "Associative table / map data structures."},
	"TERRAIN":  {Doc: "reference/TERRAIN.md", Blurb: "Heightfield terrain chunks."},
	"TEXTDRAW": {Doc: "reference/DRAW2D.md", Blurb: "Legacy text-draw aliases (`TEXTDRAW.*`)."},
	"TEXTURE":  {Doc: "reference/TEXTURE.md", Blurb: "GPU textures and procedural gen."},
	"TILEMAP":  {Doc: "reference/TILEMAP.md", Blurb: "2D tile maps and layers."},
	"TIME":     {Doc: "reference/TIME.md", Blurb: "Frame delta and timers."},
	"TRANSFORM": {Doc: "reference/TRANSFORM.md", Blurb: "TRS and hierarchy transforms."},
	"TRANSITION": {Doc: "reference/TRANSITION.md", Blurb: "Screen transitions."},
	"TRIGGER":  {Doc: "reference/PHYSICS3D.md", Blurb: "Physics trigger volumes (`TRIGGER.*`)."},
	"TWEEN":    {Doc: "reference/TWEEN.md", Blurb: "Interpolation and easing helpers."},
	"UI":       {Doc: "reference/GUI.md", Blurb: "UI layout helpers (`UI.*`); see also GUI."},
	"UTIL":     {Doc: "reference/UTIL.md", Blurb: "Miscellaneous utilities."},
	"VEC2":     {Doc: "reference/VEC_QUAT.md", Blurb: "Two-component vectors."},
	"VEC3":     {Doc: "reference/VEC_QUAT.md", Blurb: "Three-component vectors."},
	"WAVE":     {Doc: "reference/WAVE.md", Blurb: "In-memory wave samples."},
	"WATER":    {Doc: "reference/WATER.md", Blurb: "Water surfaces and rendering."},
	"WEATHER":  {Doc: "reference/WEATHER.md", Blurb: "Weather state and atmosphere."},
	"WIND":     {Doc: "reference/WEATHER.md", Blurb: "Wind vectors for effects and foliage."},
	"WINDOW":   {Doc: "reference/WINDOW.md", Blurb: "Window, OpenGL context, and platform."},
	"WORLD":    {Doc: "reference/WORLD.md", Blurb: "Open-world streaming center, preload, WORLD.UPDATE, WORLD.ISREADY."},
}

func main() {
	repoRoot := findRepoRoot()
	jsonPath := filepath.Join(repoRoot, "compiler", "builtinmanifest", "commands.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", jsonPath, err)
		os.Exit(1)
	}
	var r root
	if err := json.Unmarshal(data, &r); err != nil {
		fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		os.Exit(1)
	}

	byNS := make(map[string]int)
	descRows := 0
	for _, c := range r.Commands {
		k := c.Key
		var ns string
		if idx := strings.IndexByte(k, '.'); idx >= 0 {
			if c.Namespace != "" {
				ns = strings.ToUpper(c.Namespace)
			} else {
				ns = k[:idx]
			}
		} else {
			ns = "(global)"
		}
		byNS[ns]++
		if strings.TrimSpace(c.Description) != "" {
			descRows++
		}
	}

	var nss []string
	for ns := range byNS {
		nss = append(nss, ns)
	}
	sort.Strings(nss)

	var b strings.Builder
	b.WriteString("# Command documentation audit\n\n")
	b.WriteString("Generated by `go run ./tools/cmdaudit` from the repository root.\n\n")
	b.WriteString("This file cross-checks `compiler/builtinmanifest/commands.json` against reference docs under `docs/`.\n\n")
	fmt.Fprintf(&b, "- **Manifest rows (overloads):** %d\n", len(r.Commands))
	fmt.Fprintf(&b, "- **Dotted namespaces:** %d (plus `(global)` builtins)\n", len(byNS)-1)
	fmt.Fprintf(&b, "- **Rows with inline `description` in JSON:** %d (optional prose; most docs live in topic pages)\n\n", descRows)

	b.WriteString("## Namespace → reference\n\n")
	b.WriteString("| Namespace | Overloads | Primary doc | Blurb |\n")
	b.WriteString("|-----------|----------:|---------------|-------|\n")
	var unmapped []string
	var missingDoc []string
	for _, ns := range nss {
		ref, ok := namespaceRef[ns]
		if !ok {
			unmapped = append(unmapped, ns)
			fmt.Fprintf(&b, "| `%s` | %d | *(unmapped — add to tools/cmdaudit)* |  |\n", ns, byNS[ns])
			continue
		}
		docPath := filepath.Join(repoRoot, "docs", filepath.FromSlash(ref.Doc))
		exists := "yes"
		if _, err := os.Stat(docPath); err != nil {
			exists = "**missing file**"
			missingDoc = append(missingDoc, ref.Doc)
		}
		blurb := strings.ReplaceAll(ref.Blurb, "|", "\\|")
		fmt.Fprintf(&b, "| `%s` | %d | [%s](%s) (%s) | %s |\n", ns, byNS[ns], ref.Doc, ref.Doc, exists, blurb)
	}
	b.WriteString("\n")

	if len(unmapped) > 0 {
		b.WriteString("## Unmapped namespaces\n\n")
		b.WriteString("Add each to `namespaceRef` in `tools/cmdaudit/main.go`.\n\n")
		for _, ns := range unmapped {
			fmt.Fprintf(&b, "- `%s` (%d overloads)\n", ns, byNS[ns])
		}
		b.WriteString("\n")
	}
	if len(missingDoc) > 0 {
		b.WriteString("## Missing documentation files\n\n")
		sort.Strings(missingDoc)
		for _, p := range missingDoc {
			fmt.Fprintf(&b, "- `docs/%s`\n", p)
		}
		b.WriteString("\n")
	}

	b.WriteString("## How documentation is organized\n\n")
	b.WriteString("1. **`docs/COMMANDS.md`** — Topic index, human-written explanations, DONE/PARTIAL status, and tables for major modules.\n")
	b.WriteString("2. **`docs/reference/*.md`** — Deep dives per subsystem (window, render, physics, …).\n")
	b.WriteString("3. **`docs/API_CONSISTENCY.md`** — Generated list of every manifest overload and arity (`go run ./tools/apidoc`).\n")
	b.WriteString("4. **`docs/reference/API_CONVENTIONS.md`** — Recommended verbs (`LOAD`, `SETPOS`, …) across object types.\n")
	b.WriteString("5. **Optional `description` fields** in `commands.json` — Extra inline help for tools; not required for every row.\n\n")

	out := filepath.Join(repoRoot, "docs", "COMMAND_AUDIT.md")
	if err := os.WriteFile(out, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", out)
	if len(unmapped) > 0 {
		fmt.Fprintf(os.Stderr, "warning: %d unmapped namespaces\n", len(unmapped))
		os.Exit(2)
	}
}

func findRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for d := dir; d != filepath.VolumeName(d)+string(filepath.Separator); d = filepath.Dir(d) {
		if st, err := os.Stat(filepath.Join(d, "go.mod")); err == nil && !st.IsDir() {
			return d
		}
	}
	return dir
}
