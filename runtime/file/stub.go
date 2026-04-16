//go:build !cgo && !windows

package mbfile

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Module is a no-op when CGO is disabled; FILE.* builtins return a clear error.
type Module struct{}

// NewModule creates the stub module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(*heap.Store) {}

var stubNames = []string{
	"FILE.OPEN",
	"FILE.OPENREAD",
	"FILE.OPENWRITE",
	"FILE.CLOSE",
	"FILE.READLINE",
	"FILE.WRITE",
	"FILE.WRITELN",
	"FILE.EOF",
	"FILE.SEEK",
	"FILE.TELL",
	"FILE.SIZE",
	"FILE.EXISTS", "FILE.READALLTEXT", "FILE.WRITEALLTEXT", "FILE.READSTRING",
	// Flat aliases (same CGO requirement as FILE.*)
	"OPENFILE",
	"CLOSEFILE",
	"READFILE",
	"READFILE$",
	"READSTRING",
	"READSTRING$",
	"WRITEFILE",
	"WRITEFILELN",
	"EOF",
	"FILEPOS",
	"SEEKFILE",
	"FILESIZE",
	"WriteFile", "ReadFile", "CloseFile", "WriteLine", "ReadLine$",
	"WriteInt", "ReadInt", "WriteFloat", "ReadFloat",
}

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	for _, n := range stubNames {
		name := n
		r.Register(n, "file", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, runtime.Errorf("%s requires a CGO-enabled build (file I/O)", name)
		})
	}
	m.registerFileBlitz(r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
