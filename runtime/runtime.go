// Package runtime implements the moonBASIC native command layer.
// It bridges the VM to external libraries (Raylib, Jolt, Box2D, ENet).
package runtime

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"moonbasic/compiler/builtinmanifest"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

// Errorf formats a runtime error with the moonBASIC prefix (for use from subpackages).
func Errorf(format string, a ...any) error {
	return fmt.Errorf("[moonBASIC] Runtime Error: "+format, a...)
}

// Runtime is an alias for Registry, used by native functions to access runtime services.
type Runtime = Registry

// BuiltinFn is the Go signature for every moonBASIC native command.
// It receives a pointer to the active runtime and a slice of argument values.
type BuiltinFn func(rt *Runtime, args ...value.Value) (value.Value, error)

// AdaptLegacy wraps handlers that only need the flat argument slice (pre-runtime-aware shape).
func AdaptLegacy(fn func(args []value.Value) (value.Value, error)) BuiltinFn {
	return func(rt *Runtime, args ...value.Value) (value.Value, error) {
		return fn(args)
	}
}

// Registrar provides an interface for modules to register their commands.
type Registrar interface {
	Register(name, namespace string, fn BuiltinFn)
}

// Module is implemented by every moonBASIC runtime module (Window, Render, etc).
type Module interface {
	Register(r Registrar)
	Shutdown()
}

// HeapAware modules receive the VM heap before Register (e.g. FILE.* allocates handles).
type HeapAware interface {
	BindHeap(h *heap.Store)
}

// Registry manages the global dispatch table and handle heap.
type Registry struct {
	mu       sync.RWMutex
	Commands map[string]BuiltinFn
	Heap     *heap.Store
	Modules  []Module
	// Prog is the bytecode program currently executing (set by vm.VM.Execute); used for string pool resolution.
	Prog *opcode.Program
	// DiagOut receives DEBUG.* and similar diagnostics (pipeline sets this to Options.Out).
	DiagOut io.Writer
	// StackTraceFn is set by vm.VM.Execute while running; natives can call it for DEBUG.STACKTRACE.
	StackTraceFn func() string
	// TerminateVM is set by vm.VM.Execute; QUIT/STOP call it to end the program before normal main return.
	TerminateVM func()
	// HostArgs is process argv for ARGC / COMMAND$; nil means use os.Args. A non-nil empty slice is a deliberate empty argv.
	HostArgs []string
}

// NewRegistry initializes the runtime environment.
func NewRegistry(h *heap.Store) *Registry {
	return &Registry{
		Commands: make(map[string]BuiltinFn),
		Heap:     h,
		Modules:  []Module{},
	}
}

// Register registers a native Go function to a command name.
func (r *Registry) Register(name, namespace string, fn BuiltinFn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Commands == nil {
		r.Commands = make(map[string]BuiltinFn)
	}
	r.Commands[strings.ToUpper(name)] = fn
}

// Bind is a legacy method for registration. Use Register instead.
func (r *Registry) Bind(name string, fn func(args []value.Value) (value.Value, error)) {
	r.Register(name, "legacy", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		return fn(args)
	})
}

// RegisterFromManifest automatically stubs out any command listed in the manifest
// that doesn't already have a native implementation.
func (r *Registry) RegisterFromManifest(table *builtinmanifest.Table) {
	if table == nil || table.Commands == nil {
		return
	}
	seen := make(map[string]bool)
	for _, overloads := range table.Commands {
		for _, cmd := range overloads {
			key := cmd.Key
			if seen[key] {
				continue
			}
			seen[key] = true
			r.mu.RLock()
			_, exists := r.Commands[key]
			r.mu.RUnlock()

			if !exists {
				name := key // Capture for closure
				r.Register(key, cmd.Namespace, func(rt *Runtime, args ...value.Value) (value.Value, error) {
					return value.Value{}, Errorf("command %s is not yet implemented", name)
				})
			}
		}
	}
}

// CommandKeys returns a snapshot of registered built-in keys (for diagnostics / did-you-mean).
func (r *Registry) CommandKeys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.Commands))
	for k := range r.Commands {
		out = append(out, k)
	}
	return out
}

// Call executes a command by its fully qualified name.
func (r *Registry) Call(name string, args []value.Value) (value.Value, error) {
	key := strings.ToUpper(name)
	r.mu.RLock()
	fn, ok := r.Commands[key]
	r.mu.RUnlock()
	if !ok {
		return value.Value{}, FormatUnknownRegistryCommand(key, r.CommandKeys())
	}
	exit := enterCall(r)
	defer exit()
	return fn(r, args...)
}

// Shutdown releases all module-level resources.
func (r *Registry) Shutdown() {
	for _, m := range r.Modules {
		m.Shutdown()
	}
	r.Heap.FreeAll()
}

// RegisterModule adds a module to the registry and performs its registration.
func (r *Registry) RegisterModule(m Module) {
	if ha, ok := m.(HeapAware); ok {
		ha.BindHeap(r.Heap)
	}
	m.Register(r)
	r.Modules = append(r.Modules, m)
}
