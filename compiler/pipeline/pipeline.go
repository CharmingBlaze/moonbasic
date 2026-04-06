// Package pipeline orchestrates the moonBASIC compiler and VM execution stages.
// It serves as the primary library entry point for moonBASIC host applications.
package pipeline

import (
	"fmt"
	"io"
	"os"
	goruntime "runtime"
	"strings"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/builtinmanifest"
	"moonbasic/compiler/codegen"
	"moonbasic/compiler/include"
	"moonbasic/compiler/parser"
	"moonbasic/compiler/semantic"
	"moonbasic/runtime"
	mbaudio "moonbasic/runtime/audio"
	mbiome "moonbasic/runtime/biome"
	"moonbasic/runtime/bitwise"
	mbcamera "moonbasic/runtime/camera"
	mbcharcontroller "moonbasic/runtime/charcontroller"
	mcloud "moonbasic/runtime/cloudmod"
	mbcsv "moonbasic/runtime/csvmod"
	mbdb "moonbasic/runtime/dbmod"
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
	mbnet "moonbasic/runtime/net"
	mnoise "moonbasic/runtime/noisemod"
	mbphysics2d "moonbasic/runtime/physics2d"
	mbphysics3d "moonbasic/runtime/physics3d"
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
	"moonbasic/vm/moon"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

// Options carries configuration for the VM execution.
type Options struct {
	Debug bool      // If true, print disassembly before execution
	Trace bool      // If true, print VM state after each opcode
	Out   io.Writer // Output stream for trace and errors (default os.Stderr)

	// ProfileRecorder when non-nil accumulates per-source-line instruction counts during Execute.
	ProfileRecorder *vm.ProfileRecorder

	// HostArgs is argv for ARGC / COMMAND$; nil leaves Registry.HostArgs nil so those builtins use os.Args.
	HostArgs []string
}

// CompileSource parses, analyzes, and generates code from a string.
func CompileSource(name, src string) (*opcode.Program, error) {
	SyncPackageIncludeRoots()
	lines := parser.SplitLines(src)
	ar := arena.NewArena()
	defer ar.Reset()

	// 1. Parsing
	prog, err := parser.ParseSourceWithArena(name, src, ar)
	if err != nil {
		return nil, err
	}
	prog, err = include.ExpandWithArena(name, prog, ar)
	if err != nil {
		return nil, err
	}

	// 2. Semantic Analysis
	an := semantic.DefaultAnalyzer(name, lines)
	if err := an.Run(prog); err != nil {
		return nil, err
	}

	// 3. Code Generation
	g := codegen.New(name, lines)
	bc, err := g.Compile(prog)
	if err != nil {
		return nil, fmt.Errorf("[moonBASIC] CodeGen Error: %v", err)
	}

	return bc, nil
}

// CompileFile reads a file from disk and compiles it.
func CompileFile(path string) (*opcode.Program, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return CompileSource(path, string(data))
}

// CheckFile reads a file from disk and performs only semantic analysis.
func CheckFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return CheckSource(path, string(data))
}

// CheckSource performs parsing and semantic analysis only.
func CheckSource(name, src string) error {
	SyncPackageIncludeRoots()
	ar := arena.NewArena()
	defer ar.Reset()
	prog, err := parser.ParseSourceWithArena(name, src, ar)
	if err != nil {
		return err
	}
	prog, err = include.ExpandWithArena(name, prog, ar)
	if err != nil {
		return err
	}
	an := semantic.DefaultAnalyzer(name, parser.SplitLines(src))
	return an.Run(prog)
}

// ListBuiltins returns all registered native command keys.
func ListBuiltins() []string {
	h := heap.New()
	reg := runtime.NewRegistry(h)
	// Use empty options for listing; doesn't matter for registration
	setupRegistry(reg, h, Options{})
	return reg.CommandKeys()
}

// RunProgram initializes the runtime and executes the program in the VM.
func RunProgram(prog *opcode.Program, opts Options) error {
	if opts.Out == nil {
		opts.Out = os.Stderr
	}

	if opts.Debug {
		fmt.Fprintln(opts.Out, prog.Main.Disassemble())
	}

	goruntime.LockOSThread()

	// 1. Initialize Runtime
	h := heap.New()
	reg := runtime.NewRegistry(h)
	setupRegistry(reg, h, opts)

	// 2. Setup VM
	machine := vm.New(reg, h)
	// Wire up modules that need to call back into the VM (using the machine just created)
	wireRegistryCallbacks(reg, machine)

	machine.Trace = opts.Trace
	machine.TraceOut = opts.Out
	machine.StackHygieneDebug = opts.Debug
	machine.Profiler = opts.ProfileRecorder

	defer reg.Shutdown() // Raylib + heap cleanup on success or VM error

	// 3. Execution
	if err := machine.Execute(prog); err != nil {
		return err
	}

	return nil
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
	reg.RegisterModule(mbnet.NewModule())

	// Physics / character: register before manifest so natives win; char before physics3d
	// so Shutdown frees CharacterVirtual instances before the Jolt world is torn down.
	reg.RegisterModule(mbcharcontroller.NewModule())
	reg.RegisterModule(mbphysics2d.NewModule())
	reg.RegisterModule(mbphysics3d.NewModule())
	reg.RegisterModule(mbcollision.NewModule())
	reg.RegisterModule(mnoise.NewModule())
	reg.RegisterModule(mbgame.NewModule())

	// Stubs for manifest entries not yet implemented natively
	reg.RegisterFromManifest(builtinmanifest.Default())
}

func wireRegistryCallbacks(reg *runtime.Registry, machine *vm.VM) {
	// Find modules that need callbacks
	var sceneMod *mbscene.Module
	var poolMod *mbpool.Module
	var tweenMod *mbtween.Module
	var eventMod *mbevent.Module
	var p3 *mbphysics3d.Module
	var navMod *mbnav.Module
	var netMod *mbnet.Module

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
		case *mbphysics3d.Module:
			p3 = mod
		case *mbnav.Module:
			navMod = mod
		case *mbnet.Module:
			netMod = mod
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
	if p3 != nil {
		p3.SetUserInvoker(machine.CallUserFunction)
	}
	if navMod != nil {
		navMod.SetUserInvoker(machine.CallUserFunction)
	}
	if netMod != nil {
		netMod.SetUserInvoker(machine.CallUserFunction)
	}

	h := reg.Heap
	runtime.SeedInputKeyGlobals(machine.Globals)
	runtime.SeedBlendModeGlobals(machine.Globals)
	window.SeedWindowFlagGlobals(machine.Globals)
	input.SeedGestureGlobals(machine.Globals)
	mbmodel3d.SeedMaterialMapGlobals(machine.Globals)
	mbmatrix.SeedColorGlobals(h, machine.Globals)
	mbnet.SeedMultiplayerGlobals(machine.Globals)
	texture.SeedTextureGlobals(machine.Globals)
	mbgui.SeedGUIGlobals(machine.Globals)
}


// EncodeMOON serializes a compiled program to MOON container bytes (.mbc).
func EncodeMOON(prog *opcode.Program) ([]byte, error) {
	return moon.Encode(prog)
}

// DecodeMOON loads a program from MOON bytes after validating magic and version.
func DecodeMOON(data []byte) (*opcode.Program, error) {
	return moon.Decode(data)
}
