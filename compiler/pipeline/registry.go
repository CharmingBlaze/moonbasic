//go:build fullruntime

package pipeline

import (
	"strings"

	"moonbasic/compiler/builtinmanifest"
	"moonbasic/internal/driver"
	"moonbasic/runtime"
	mbaudio "moonbasic/runtime/audio"
	mbiome "moonbasic/runtime/biome"
	"moonbasic/runtime/bitwise"
	mbcamera "moonbasic/runtime/camera"
	mcloud "moonbasic/runtime/cloudmod"
	mbcsv "moonbasic/runtime/csvmod"
	mbdb "moonbasic/runtime/dbmod"
	mbblitz "moonbasic/runtime/blitzengine"
	mbdraw "moonbasic/runtime/draw"
	mbfile "moonbasic/runtime/file"
	mbfont "moonbasic/runtime/font"
	"moonbasic/runtime/input"
	mbjson "moonbasic/runtime/jsonmod"
	"moonbasic/runtime/mathmod"
	mbarray "moonbasic/runtime/mbarray"
	mbcollision "moonbasic/runtime/mbcollision"
	mbdata "moonbasic/runtime/mbdata"
	mbdebug "moonbasic/runtime/mbdebug"
	mbevent "moonbasic/runtime/mbevent"
	mbgame "moonbasic/runtime/mbgame"
	mbentity "moonbasic/runtime/mbentity"
	mbgui "moonbasic/runtime/mbgui"
	"moonbasic/runtime/mbimage"
	mblight "moonbasic/runtime/mblight"
	mblight2d "moonbasic/runtime/mblight2d"
	mbmatrix "moonbasic/runtime/mbmatrix"
	mbmem "moonbasic/runtime/mbmem"
	"moonbasic/runtime/mbmodel3d"
	mbnav "moonbasic/runtime/mbnav"
	mbparticles "moonbasic/runtime/mbparticles"
	mbpool "moonbasic/runtime/mbpool"
	mbrand "moonbasic/runtime/mbrand"
	mbscene "moonbasic/runtime/mbscene"
	mbtilemap "moonbasic/runtime/mbtilemap"
	mbtransition "moonbasic/runtime/mbtransition"
	mbtween "moonbasic/runtime/mbtween"
	mbutil "moonbasic/runtime/mbutil"
	mnoise "moonbasic/runtime/noisemod"
	mscatter "moonbasic/runtime/scatter"
	msky "moonbasic/runtime/sky"
	mbsprite "moonbasic/runtime/sprite"
	"moonbasic/runtime/strmod"
	mbsystem "moonbasic/runtime/system"
	mbtable "moonbasic/runtime/tablemod"
	terrain "moonbasic/runtime/terrain"
	"moonbasic/runtime/texture"
	mbtime "moonbasic/runtime/time"
	mwater "moonbasic/runtime/water"
	mweather "moonbasic/runtime/weathermod"
	"moonbasic/runtime/window"
	worldmgr "moonbasic/runtime/worldmgr"
	"moonbasic/vm"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// ListBuiltins returns all registered native command keys.
func ListBuiltins() []string {
	h := heap.New()
	reg := runtime.NewRegistry(h)
	// Use empty options for listing; doesn't matter for registration
	setupRegistry(reg, h, Options{})
	return reg.CommandKeys()
}

func setupRegistry(reg *runtime.Registry, h *heap.Store, opts Options) {
	reg.DebugMode = opts.Debug
	reg.DiagOut = opts.Out
	if opts.HostArgs != nil {
		reg.HostArgs = opts.HostArgs
	}
	reg.InitCore() // Register core built-ins (PRINT, etc)
	reg.RegisterModule(bitwise.NewModule())
	reg.RegisterModule(strmod.NewModule())

	debugMod := mbdebug.NewModule()

	// Native WINDOW / minimal RENDER (Raylib when CGO enabled; stubs otherwise)
	winMod := window.NewModule()
	winMod.BindDriverSelection(driver.GetDefaultDriver())
	mblight2d.RegisterFrameHook(winMod)
	mbtransition.RegisterFrameHook(winMod)
	winMod.SetFrameEndHook(debugMod.DrawFrameOverlay)
	winMod.SetDiagnostics(opts.Out, opts.Debug)
	audMod := mbaudio.NewModule()
	winMod.SetAudioHooks(audMod.OnWindowOpen, audMod.OnWindowClose)
	reg.RegisterModule(winMod)
	reg.RegisterModule(input.NewModule())
	reg.RegisterModule(mathmod.NewModule())
	reg.RegisterModule(mbmatrix.NewModule())
	reg.RegisterModule(mbtime.NewModule())
	reg.RegisterModule(debugMod)
	reg.RegisterModule(mblight.NewModule())
	reg.RegisterModule(mbsystem.NewModule())
	reg.RegisterModule(mbfile.NewModule())
	reg.RegisterModule(mbmem.NewModule())
	reg.RegisterModule(mbarray.NewModule())
	reg.RegisterModule(mbdata.NewModule())
	reg.RegisterModule(mbrand.NewModule())
	reg.RegisterModule(mbutil.NewModule())
	reg.RegisterModule(mbdraw.NewModule())
	reg.RegisterModule(texture.NewModule())
	reg.RegisterModule(mbimage.NewModule())
	reg.RegisterModule(mbmodel3d.NewModule())
	reg.RegisterModule(mbparticles.NewModule())
	reg.RegisterModule(mbcamera.NewModule())
	reg.RegisterModule(mbsprite.NewModule())
	reg.RegisterModule(mbtilemap.NewModule())
	reg.RegisterModule(mbscene.NewModule())
	reg.RegisterModule(mbpool.NewModule())
	reg.RegisterModule(mbtween.NewModule())
	reg.RegisterModule(mbevent.NewModule())
	reg.RegisterModule(mbnav.NewModule())
	reg.RegisterModule(mblight2d.NewModule())
	reg.RegisterModule(mbtransition.NewModule())
	reg.RegisterModule(mbfont.NewModule())
	reg.RegisterModule(mbgui.NewModule())
	reg.RegisterModule(audMod)
	reg.RegisterModule(mbjson.NewModule())
	reg.RegisterModule(mbcsv.NewModule())
	reg.RegisterModule(mbdb.NewModule())
	reg.RegisterModule(mbtable.NewModule())
	terrMod := terrain.NewModule()
	worldMod := worldmgr.NewModule(terrMod)
	reg.RegisterModule(terrMod)
	reg.RegisterModule(worldMod)
	reg.RegisterModule(mwater.NewModule())
	reg.RegisterModule(msky.NewModule())
	reg.RegisterModule(mcloud.NewModule())
	reg.RegisterModule(mweather.NewModule())
	reg.RegisterModule(mscatter.NewModule())
	reg.RegisterModule(mbiome.NewModule())

	// Call separate registration functions for physics and networking
	// (These functions are defined in other files with build tags)
	registerPhysicsModules(reg)
	registerNetModules(reg)

	reg.RegisterModule(mbcollision.NewModule())
	reg.RegisterModule(mnoise.NewModule())
	reg.RegisterModule(mbentity.NewModule())
	reg.RegisterModule(mbgame.NewModule())
	reg.RegisterModule(mbblitz.NewModule())

	// Stubs for manifest entries not yet implemented natively
	reg.RegisterFromManifest(builtinmanifest.Default())
}

func wireRegistryCallbacks(reg *runtime.Registry, machine *vm.VM) {
	// Find modules that need callbacks
	var sceneMod *mbscene.Module
	var poolMod *mbpool.Module
	var tweenMod *mbtween.Module
	var eventMod *mbevent.Module
	var navMod *mbnav.Module

	for _, m := range reg.Modules {
		switch mod := m.(type) {
		case *mbscene.Module:
			sceneMod = mod
		case *mbpool.Module:
			poolMod = mod
		case *mbtween.Module:
			tweenMod = mod
		case *mbevent.Module:
			eventMod = mod
		case *mbnav.Module:
			navMod = mod
		}
	}

	if sceneMod != nil {
		sceneMod.SetUserInvoker(machine.CallUserFunction)
	}
	if poolMod != nil {
		poolMod.SetUserInvoker(machine.CallUserFunction)
	}
	if tweenMod != nil {
		tweenMod.SetUserInvoker(machine.CallUserFunction)
		tweenMod.SetGlobalAccessor(
			func(k string) (value.Value, bool) {
				k = strings.ToUpper(strings.TrimSpace(k))
				v, ok := machine.Globals[k]
				return v, ok
			},
			func(k string, v value.Value) {
				machine.Globals[strings.ToUpper(strings.TrimSpace(k))] = v
			},
		)
	}
	if eventMod != nil {
		eventMod.SetUserInvoker(machine.CallUserFunction)
	}
	if navMod != nil {
		navMod.SetUserInvoker(machine.CallUserFunction)
	}

	// Dynamic wiring for optional modules
	wirePhysicsCallbacks(reg, machine)
	wireNetCallbacks(reg, machine)

	h := reg.Heap
	runtime.SeedInputKeyGlobals(machine.Globals)
	runtime.SeedBlendModeGlobals(machine.Globals)
	window.SeedWindowFlagGlobals(machine.Globals)
	input.SeedGestureGlobals(machine.Globals)
	mbmodel3d.SeedMaterialMapGlobals(machine.Globals)
	mbmatrix.SeedColorGlobals(h, machine.Globals)
	texture.SeedTextureGlobals(machine.Globals)
	mbgui.SeedGUIGlobals(machine.Globals)
}
