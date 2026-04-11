//go:build cgo || (windows && !cgo)

package mbfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type fileObj struct {
	f       *os.File
	rd      *bufio.Reader
	wr      *bufio.Writer
	name    string
	release heap.ReleaseOnce
}

func (o *fileObj) TypeName() string { return "File" }
func (o *fileObj) TypeTag() uint16  { return heap.TagFile }
func (o *fileObj) Free() {
	o.release.Do(func() {
		if o.wr != nil {
			o.wr.Flush()
		}
		if o.f != nil {
			o.f.Close()
			o.f = nil
		}
		o.wr = nil
		o.rd = nil
	})
}

type Module struct {
	h *heap.Store
}

// NewModule creates the file I/O module (CGO builds only).
func NewModule() *Module { return &Module{} }

func (m *Module) Names() []string {
	return []string{
		"FILE.OPENREAD", "FILE.OPENWRITE", "FILE.CLOSE", "FILE.READLINE", "FILE.WRITE", "FILE.WRITELN", "FILE.EOF",
		"FILE.SEEK", "FILE.TELL", "FILE.SIZE",
	}
}

func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	for _, n := range m.Names() {
		name := n
		r.Register(n, "file", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return m.Run(rt, name, args...)
		})
	}
	r.Register("FILE.OPEN", "file", m.fileOpen)

	// Flat spec names (manifest) → same behavior as FILE.* (see COMMAND_AUDIT).
	r.Register("OPENFILE", "file", m.fileOpen)
	flatForward := []struct{ flat, canon string }{
		{"CLOSEFILE", "FILE.CLOSE"},
		{"READFILE", "FILE.READLINE"},
		{"WRITEFILE", "FILE.WRITE"},
		{"EOF", "FILE.EOF"},
		{"FILEPOS", "FILE.TELL"},
		{"SEEKFILE", "FILE.SEEK"},
		{"FILESIZE", "FILE.SIZE"},
	}
	for _, a := range flatForward {
		canon := a.canon
		r.Register(a.flat, "file", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return m.Run(rt, canon, args...)
		})
	}
	r.Register("WRITEFILELN", "file", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.Run(rt, "FILE.WRITELN", args...)
	})
	m.registerFileExtras(r)
	m.registerFileBlitz(r)
	m.registerBankTransfer(r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) fileOpen(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("FILE.OPEN expects (path, mode)")
	}
	mode, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	mode = strings.ToLower(strings.TrimSpace(mode))
	switch mode {
	case "r", "read":
		return m.Run(rt, "FILE.OPENREAD", args[0])
	case "w", "write":
		return m.Run(rt, "FILE.OPENWRITE", args[0])
	default:
		return value.Nil, runtime.Errorf("FILE.OPEN: mode must be r or w, got %q", mode)
	}
}

func (m *Module) Run(rt *runtime.Runtime, name string, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("%s: heap not bound", name)
	}
	switch name {
	case "FILE.OPENREAD":
		return m.fileOpenRead(rt, args...)
	case "FILE.OPENWRITE":
		return m.fileOpenWrite(rt, args...)
	case "FILE.CLOSE":
		return m.fileClose(args...)
	case "FILE.READLINE":
		return m.fileReadLine(rt, args...)
	case "FILE.WRITE":
		return m.fileWriteRaw(rt, args...)
	case "FILE.WRITELN":
		return m.fileWriteLine(rt, args...)
	case "FILE.EOF":
		return m.fileEOF(args...)
	case "FILE.SEEK":
		return m.fileSeek(args...)
	case "FILE.TELL":
		return m.fileTell(args...)
	case "FILE.SIZE":
		return m.fileSize(args...)
	}
	return value.Nil, runtime.Errorf("unknown file command: %s", name)
}

func (m *Module) fileOpenRead(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("FILE.OPENREAD expects 1 string argument (path)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return value.Nil, err
	}
	o := &fileObj{f: f, rd: bufio.NewReader(f), name: path}
	id, err := m.h.Alloc(o)
	if err != nil {
		f.Close()
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) fileOpenWrite(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("FILE.OPENWRITE expects 1 string argument (path)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	f, err := os.Create(path)
	if err != nil {
		return value.Nil, err
	}
	o := &fileObj{f: f, wr: bufio.NewWriter(f), name: path}
	id, err := m.h.Alloc(o)
	if err != nil {
		f.Close()
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) fileClose(args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("FILE.CLOSE expects 1 handle argument")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) fileReadLine(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("FILE.READLINE expects 1 handle argument")
	}
	o, err := heap.Cast[*fileObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if o.rd == nil {
		return value.Nil, fmt.Errorf("file not open for reading")
	}
	line, err := o.rd.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		return value.Nil, err
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	return rt.RetString(line), nil
}

func (m *Module) fileEOF(args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("FILE.EOF expects 1 handle argument")
	}
	o, err := heap.Cast[*fileObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if o.rd == nil {
		return value.FromBool(true), nil
	}
	_, err = o.rd.Peek(1)
	return value.FromBool(err != nil), nil
}

func (m *Module) fileSeek(args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("FILE.SEEK expects (handle, offset)")
	}
	o, err := heap.Cast[*fileObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	off, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("FILE.SEEK: offset must be numeric")
	}
	if o.wr != nil {
		o.wr.Flush()
	}
	_, err = o.f.Seek(off, 0)
	if err != nil {
		return value.Nil, err
	}
	if o.rd != nil {
		o.rd.Reset(o.f)
	}
	if o.wr != nil {
		o.wr.Reset(o.f)
	}
	return value.Nil, nil
}

func (m *Module) fileTell(args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("FILE.TELL expects 1 handle argument")
	}
	o, err := heap.Cast[*fileObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	curr, err := o.f.Seek(0, 1)
	if err != nil {
		return value.Nil, err
	}
	// Subtract buffered data if reading
	if o.rd != nil {
		curr -= int64(o.rd.Buffered())
	}
	return value.FromInt(curr), nil
}

func (m *Module) fileSize(args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("FILE.SIZE expects 1 handle argument")
	}
	o, err := heap.Cast[*fileObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	fi, err := o.f.Stat()
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(fi.Size()), nil
}

func (m *Module) Reset() {}

