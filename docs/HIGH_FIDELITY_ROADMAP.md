# moonBASIC High-Fidelity Engine Architecture Roadmap

## Current Architecture Analysis

### ✅ What's Already Implemented (Strong Foundation)

1. **Register-Based VM**: The codebase already uses a register-based bytecode VM (not stack-based)
   - 8-byte fixed-width instructions (IR v3)
   - Register allocation in `codegen.go` with `nextReg`/`baseReg`
   - Frame-based call stack with registers
   - See: `vm/vm.go`, `vm/opcode/opcode.go`

2. **Symbol Table System**: `compiler/symtable/symtable.go`
   - Tracks locals, globals, params, statics
   - Predeclaration pass for functions/types
   - Slot-based allocation for locals

3. **Modular Compiler Pipeline**: `compiler/pipeline/compile.go`
   - Lexer → Parser → Semantic Analysis → CodeGen
   - Arena allocation support for AST
   - Include file expansion

4. **Modern Value System**: `vm/value/value.go`
   - `float64`, `int64`, `string`, `bool`, `handle` kinds
   - Union-based storage for efficiency

---

## Implementation Phases

### Phase 1: Modern Syntax (Implicit Declaration) - 2 weeks

**Goal**: Remove `VAR` requirement, implement first-assignment declaration

#### 1.1 Two-Pass Symbol Table Builder

Create `compiler/symtable/symtable_builder.go`:

```go
// FirstPass collects all variable assignments before code generation
type FirstPass struct {
    Symbols *symtable.Table
    ScopeStack []map[string]bool
}

func (fp *FirstPass) VisitProgram(prog *ast.Program) {
    // Pass 1: Collect all top-level assignments as globals
    for _, stmt := range prog.Stmts {
        fp.collectGlobals(stmt)
    }
    
    // Pass 2: Collect locals in each function
    for _, fn := range prog.Functions {
        fp.collectFunctionLocals(fn)
    }
}

func (fp *FirstPass) collectGlobals(stmt ast.Stmt) {
    switch s := stmt.(type) {
    case *ast.AssignNode:
        // First assignment declares the variable
        if !fp.Symbols.IsDefined(s.Name) {
            fp.Symbols.DefineGlobalImplicit(s.Name, inferType(s.Expr))
        }
    case *ast.DimNode:
        // Arrays declared via DIM
        fp.Symbols.DefineGlobalImplicit(s.Name, TypeArray)
    }
}
```

#### 1.2 Type Inference Engine

Create `compiler/types/inference.go`:

```go
// TypeTag represents inferred types
type TypeTag byte
const (
    TypeUnknown TypeTag = iota
    TypeInt
    TypeFloat  
    TypeString
    TypeBool
    TypeArray
    TypeHandle
)

// InferType deduces type from expression
func InferType(expr ast.Expr) TypeTag {
    switch e := expr.(type) {
    case *ast.IntLitNode:
        return TypeInt
    case *ast.FloatLitNode:
        return TypeFloat
    case *ast.StringLitNode:
        return TypeString
    case *ast.BoolLitNode:
        return TypeBool
    case *ast.BinopNode:
        // Arithmetic ops → float if any operand is float
        left := InferType(e.Left)
        right := InferType(e.Right)
        return promoteType(left, right)
    case *ast.CallExprNode:
        // Look up function return type from symbol table
        return InferCallReturnType(e.Name)
    case *ast.IdentNode:
        // Variable reference - look up in symbol table
        return GetVariableType(e.Name)
    }
    return TypeUnknown
}
```

#### 1.3 Suffix-Based Type Hinting

The lexer already preserves suffixes (`#`, `$`, `?`). Use them:

```go
func TypeFromSuffix(name string) TypeTag {
    if len(name) == 0 {
        return TypeUnknown
    }
    switch name[len(name)-1] {
    case '#':
        return TypeFloat
    case '$':
        return TypeString  
    case '?':
        return TypeBool
    default:
        return TypeInt  // Default to int for untyped
    }
}
```

---

### Phase 2: High-Performance Rendering - 3 weeks

**Goal**: 4K/144Hz+ with zero-allocation draw loops

#### 2.1 High-DPI & Float64 Coordinates

Update `runtime/window/raylib_cgo.go`:

```go
// WindowState holds high-DPI aware metrics
type WindowState struct {
    // Physical pixels (actual monitor resolution)
    PhysicalWidth, PhysicalHeight int
    // Logical coordinates (game units)
    LogicalWidth, LogicalHeight float64
    // Scale factor for High-DPI
    DPIScale float64
}

// Open creates a High-DPI aware window
func (m *Module) wOpen(args []value.Value) ([]value.Value, error) {
    width := int(args[0].Float64())  // Use float64 for precision
    height := int(args[1].Float64())
    title := args[2].String()
    
    // Request High-DPI window
    raylib.SetConfigFlags(raylib.FlagWindowHighdpi)
    raylib.InitWindow(width, height, title)
    
    // Get actual DPI scale
    state.DPIScale = float64(raylib.GetWindowScaleDPI())
    state.PhysicalWidth = raylib.GetScreenWidth()
    state.PhysicalHeight = raylib.GetScreenHeight()
    
    // Logical size maintains game aspect ratio
    state.LogicalWidth = float64(width)
    state.LogicalHeight = float64(height)
}
```

#### 2.2 Zero-Allocation Draw Loop

Create `runtime/draw/batch_renderer.go`:

```go
// BatchRenderer accumulates draw calls to minimize CGO crossings
type BatchRenderer struct {
    // Pre-allocated C-memory buffer for vertex data
    // Shared between Go and Raylib
    vertexBuffer unsafe.Pointer
    vertexCount  int
    maxVertices  int
    
    // Command buffer for deferred rendering
    commands []DrawCommand
}

type DrawCommand struct {
    Type   CommandType
    Color  Color
    X, Y, W, H float64  // Use float64 for precision
    // ... other fields
}

// Flush executes all batched commands in one CGO call
func (br *BatchRenderer) Flush() {
    if len(br.commands) == 0 {
        return
    }
    
    // Single CGO crossing for all draw commands
    C.batch_render(
        br.vertexBuffer,
        C.int(br.vertexCount),
        (*C.DrawCommand)(unsafe.Pointer(&br.commands[0])),
        C.int(len(br.commands))
    )
    
    br.commands = br.commands[:0]  // Reset slice (no allocation)
    br.vertexCount = 0
}
```

#### 2.3 Modern Post-Processing

Update `runtime/window/render_extra_cgo.go`:

```go
// PostProcessConfig holds effect settings
type PostProcessConfig struct {
    MSAAEnabled     bool
    MSAA samples    int  // 2, 4, 8, 16
    BloomEnabled    bool
    BloomIntensity  float32
    BloomThreshold  float32
    SSAOEnabled     bool
    SSAORadius      float32
    SSAOStrength    float32
}

// SetPostProcess configures post-processing effects
func (m *Module) SetPostProcess(cfg PostProcessConfig) {
    if cfg.MSAAEnabled {
        raylib.SetConfigFlags(raylib.FlagMsaa4xHint)  // Or 8x based on cfg
    }
    
    // Load post-processing shaders
    if cfg.BloomEnabled {
        m.bloomShader = raylib.LoadShader(nil, "resources/bloom.fs")
        raylib.SetShaderValue(m.bloomShader, 
            raylib.GetShaderLocation(m.bloomShader, "intensity"),
            []float32{cfg.BloomIntensity}, raylib.ShaderUniformFloat)
    }
    
    if cfg.SSAOEnabled {
        m.ssaoShader = raylib.LoadShader(nil, "resources/ssao.fs")
        // ... configure SSAO parameters
    }
}
```

#### 2.4 pprof-Based Allocation Audit

Add to `cmd/moonbasic/main.go`:

```go
// ProfileMode enables allocation profiling for optimization
var ProfileMode = flag.Bool("profile", false, "Enable CPU and memory profiling")

func main() {
    if *ProfileMode {
        f, _ := os.Create("cpu.prof")
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
        
        go func() {
            for {
                time.Sleep(5 * time.Second)
                var m runtime.MemStats
                runtime.ReadMemStats(&m)
                log.Printf("Alloc: %d KB, GC cycles: %d", m.Alloc/1024, m.NumGC)
            }
        }()
    }
}
```

**Target Metrics**:
- Draw loop: 0 heap allocations per frame
- Target: 144Hz at 4K (6.94ms frame budget)
- GC pauses: <1ms during gameplay

---

### Phase 3: CGO Fast-Path & Physics - 3 weeks

**Goal**: Minimize CGO bridge crossings, 1000+ bodies at 144Hz

#### 3.1 Shared Memory Buffers

Create `runtime/physics3d/shared_buffer.go`:

```go
package physics3d

/*
#include <stdlib.h>
#include <string.h>

typedef struct {
    float x, y, z;      // Position
    float qx, qy, qz, qw; // Quaternion rotation
    float vx, vy, vz;   // Linear velocity
    float ax, ay, az;   // Angular velocity
    uint32_t id;        // Body ID
    uint32_t flags;     // Active, sleeping, etc.
} BodyState;

// C-accessible functions
void batch_update_bodies(void* joltWorld, BodyState* states, int count);
void batch_get_bodies(void* joltWorld, BodyState* states, int count);
*/
import "C"
import "unsafe"

// SharedBuffer is a C-allocated memory block for batch physics updates
type SharedBuffer struct {
    ptr    unsafe.Pointer
    count  int
    stride int  // sizeof(BodyState)
}

func NewSharedBuffer(maxBodies int) *SharedBuffer {
    size := maxBodies * int(unsafe.Sizeof(C.BodyState{}))
    ptr := C.malloc(C.size_t(size))
    
    return &SharedBuffer{
        ptr:    ptr,
        count:  maxBodies,
        stride: int(unsafe.Sizeof(C.BodyState{})),
    }
}

// BatchUpdate updates 1000+ bodies in a single CGO call
func (sb *SharedBuffer) BatchUpdate(world unsafe.Pointer, bodies []BodyState) {
    count := len(bodies)
    if count > sb.count {
        count = sb.count
    }
    
    // Copy Go slice to C memory (one memcpy)
    src := unsafe.Pointer(&bodies[0])
    size := count * sb.stride
    C.memcpy(sb.ptr, src, C.size_t(size))
    
    // Single CGO call for all bodies
    C.batch_update_bodies(world, (*C.BodyState)(sb.ptr), C.int(count))
}

// BatchGet retrieves body states in single CGO call
func (sb *SharedBuffer) BatchGet(world unsafe.Pointer, bodies []BodyState) {
    count := len(bodies)
    C.batch_get_bodies(world, (*C.BodyState)(sb.ptr), C.int(count))
    
    // Copy back to Go
    dst := unsafe.Pointer(&bodies[0])
    size := count * sb.stride
    C.memcpy(dst, sb.ptr, C.size_t(size))
}
```

#### 3.2 Fixed-Step Physics with Interpolation

Create `runtime/mbgame/fixed_timestep.go`:

```go
package mbgame

// FixedTimestep manages physics at fixed rate while rendering at variable rate
type FixedTimestep struct {
    // Physics runs at fixed rate (e.g., 120Hz = 8.33ms)
    PhysicsDT float64
    
    // Accumulated time for physics steps
    accumulator float64
    
    // Previous and current physics states for interpolation
    PrevStates []Transform
    CurrStates []Transform
}

func (ft *FixedTimestep) Update(deltaTime float64, updateFunc func()) {
    ft.accumulator += deltaTime
    
    // Step physics at fixed rate
    for ft.accumulator >= ft.PhysicsDT {
        // Swap state buffers
        ft.PrevStates, ft.CurrStates = ft.CurrStates, ft.PrevStates
        
        // Run physics step
        updateFunc()
        
        // Copy current state
        copy(ft.CurrStates, currentWorldState())
        
        ft.accumulator -= ft.PhysicsDT
    }
}

// Interpolate gets smooth visual position between physics steps
func (ft *FixedTimestep) Interpolate(bodyIdx int) Transform {
    t := ft.accumulator / ft.PhysicsDT  // 0.0 to 1.0
    
    prev := ft.PrevStates[bodyIdx]
    curr := ft.CurrStates[bodyIdx]
    
    return Transform{
        Position: lerp(prev.Position, curr.Position, t),
        Rotation: slerp(prev.Rotation, curr.Rotation, t),
    }
}
```

#### 3.3 Memory Pool for Physics Bodies

```go
// BodyPool reuses C++ body objects to reduce allocation
type BodyPool struct {
    available []BodyID
    active    map[BodyID]*Body
    world     unsafe.Pointer
}

func (bp *BodyPool) Acquire() *Body {
    if len(bp.available) > 0 {
        id := bp.available[len(bp.available)-1]
        bp.available = bp.available[:len(bp.available)-1]
        return bp.active[id]
    }
    
    // Create new body (CGO call)
    body := createBody(bp.world)
    bp.active[body.ID] = body
    return body
}

func (bp *BodyPool) Release(body *Body) {
    // Reset body state instead of destroying
    body.Reset()
    bp.available = append(bp.available, body.ID)
}
```

---

### Phase 4: Resource Safety & Finalizers - 2 weeks

**Goal**: Zero memory leaks, automatic C++ resource cleanup

#### 4.1 Resource Registry with Finalizers

Update `vm/heap/heap.go`:

```go
package heap

import (
    "runtime"
    "sync"
)

// ManagedHandle wraps a C++ resource with Go finalizer
type ManagedHandle struct {
    ID       HandleID
    Type     ResourceType  // Model, Texture, Body, Sound, etc.
    CPtr     unsafe.Pointer
    RefCount int32
    
    // Cleanup function called by finalizer
    OnRelease func(cPtr unsafe.Pointer)
}

func (h *ManagedHandle) AddFinalizer() {
    runtime.SetFinalizer(h, func(mh *ManagedHandle) {
        if mh.OnRelease != nil && mh.CPtr != nil {
            mh.OnRelease(mh.CPtr)
            mh.CPtr = nil
        }
    })
}

// Store manages all heap-allocated resources
type Store struct {
    mu       sync.RWMutex
    handles  map[HandleID]*ManagedHandle
    nextID   HandleID
    
    // Per-type tracking for debugging
    stats    map[ResourceType]int
}

func (s *Store) Alloc(resourceType ResourceType, cPtr unsafe.Pointer, onRelease func(unsafe.Pointer)) HandleID {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    id := s.nextID
    s.nextID++
    
    handle := &ManagedHandle{
        ID:        id,
        Type:      resourceType,
        CPtr:      cPtr,
        OnRelease: onRelease,
        RefCount:  1,
    }
    handle.AddFinalizer()
    
    s.handles[id] = handle
    s.stats[resourceType]++
    
    return id
}

func (s *Store) Free(id HandleID) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    h, ok := s.handles[id]
    if !ok {
        return fmt.Errorf("invalid handle: %d", id)
    }
    
    // Manual release (prevents finalizer from running later)
    runtime.SetFinalizer(h, nil)
    
    if h.OnRelease != nil && h.CPtr != nil {
        h.OnRelease(h.CPtr)
    }
    
    delete(s.handles, id)
    s.stats[h.Type]--
    
    return nil
}
```

#### 4.2 Resource Leak Detection

```go
// LeakDetector tracks unfreed resources at shutdown
type LeakDetector struct {
    enabled bool
    snapshot map[HandleID]ResourceInfo
}

func (ld *LeakDetector) Snapshot() {
    if !ld.enabled {
        return
    }
    
    heapStore := GetGlobalHeap()
    ld.snapshot = heapStore.GetAllResources()
}

func (ld *LeakDetector) CheckLeaks() {
    if !ld.enabled {
        return
    }
    
    heapStore := GetGlobalHeap()
    current := heapStore.GetAllResources()
    
    for id, info := range current {
        if _, existed := ld.snapshot[id]; !existed {
            log.Printf("LEAK DETECTED: Handle %d (%s) not freed", id, info.Type)
        }
    }
    
    if len(current) > 0 {
        log.Printf("WARNING: %d resources still allocated at shutdown", len(current))
    }
}
```

---

### Phase 5: Modern Networking (ENet) - 2 weeks

**Goal**: Channel-based networking with thread isolation

#### 5.1 Channel-Based ENet Architecture

Create `runtime/net/enet_channels.go`:

```go
package net

// ChannelConfig defines reliability per channel
type ChannelConfig struct {
    ID        uint8
    Reliable  bool   // true = reliable ordered, false = unreliable sequenced
    Priority  int    // Bandwidth allocation priority
}

// Recommended channel setup for games
var DefaultChannels = []ChannelConfig{
    {ID: 0, Reliable: true,  Priority: 10},  // Critical: health, scores, game state
    {ID: 1, Reliable: true,  Priority: 5},   // Important: commands, chat
    {ID: 2, Reliable: false, Priority: 3},   // Frequent: position, rotation (interpolated)
    {ID: 3, Reliable: false, Priority: 1}, // Fire-and-forget: particles, sfx
}

// NetworkManager wraps ENet with Go channels
type NetworkManager struct {
    // ENet runs on background thread
    enetCtx    unsafe.Pointer
    
    // Thread-safe communication
    sendChan   chan Packet      // Go → ENet thread
    recvChan   chan Packet      // ENet thread → Go
    eventChan  chan Event       // Connection events
    
    // Synchronization
    shutdown   chan struct{}
    wg         sync.WaitGroup
}

type Packet struct {
    Channel uint8
    Data    []byte
    Peer    PeerID
}

// Start initializes the ENet background thread
func (nm *NetworkManager) Start(address string, port int) error {
    nm.sendChan = make(chan Packet, 256)
    nm.recvChan = make(chan Packet, 256)
    nm.eventChan = make(chan Event, 32)
    nm.shutdown = make(chan struct{})
    
    // Start ENet in background thread
    nm.wg.Add(1)
    go nm.enetThread(address, port)
    
    return nil
}

// enetThread runs ENet in isolation, communicates via channels
func (nm *NetworkManager) enetThread(address string, port int) {
    defer nm.wg.Done()
    
    // Initialize ENet (C calls)
    host := C.enet_host_create(address, port, 4)  // 4 channels
    
    for {
        select {
        case <-nm.shutdown:
            C.enet_host_destroy(host)
            return
            
        case pkt := <-nm.sendChan:
            // Send packet via ENet
            peer := C.enet_host_get_peer(host, C.int(pkt.Peer))
            flags := C.int(0)
            if pkt.Channel < 2 {  // Channels 0-1 are reliable
                flags = C.ENET_PACKET_FLAG_RELIABLE
            }
            C.enet_peer_send(peer, C.uint8_t(pkt.Channel), 
                C.CBytes(pkt.Data), C.size_t(len(pkt.Data)), flags)
                
        default:
            // Poll for incoming packets (non-blocking)
            event := C.enet_host_service(host, 0)  // 0 = non-blocking
            if event != nil {
                switch event.type {
                case C.ENET_EVENT_TYPE_RECEIVE:
                    nm.recvChan <- Packet{
                        Channel: uint8(event.channelID),
                        Data:    C.GoBytes(event.packet.data, event.packet.dataLength),
                        Peer:    PeerID(event.peer),
                    }
                case C.ENET_EVENT_TYPE_CONNECT:
                    nm.eventChan <- Event{Type: EventConnect, Peer: PeerID(event.peer)}
                case C.ENET_EVENT_TYPE_DISCONNECT:
                    nm.eventChan <- Event{Type: EventDisconnect, Peer: PeerID(event.peer)}
                }
            }
        }
    }
}

// Send queues a packet for transmission (thread-safe)
func (nm *NetworkManager) Send(channel uint8, data []byte, peer PeerID) {
    select {
    case nm.sendChan <- Packet{Channel: channel, Data: data, Peer: peer}:
    default:
        log.Printf("Network send buffer full, dropping packet")
    }
}

// Receive gets next incoming packet (non-blocking)
func (nm *NetworkManager) Receive() (Packet, bool) {
    select {
    case pkt := <-nm.recvChan:
        return pkt, true
    default:
        return Packet{}, false
    }
}
```

#### 5.2 Network Commands for MoonBASIC

```go
// Register network commands
func (m *NetModule) Register(r *registry.Registry) {
    // NET.HOST(port, maxPeers) - create server
    r.Register("NET.HOST", "net", m.host)
    
    // NET.CONNECT(address, port) - connect to server
    r.Register("NET.CONNECT", "net", m.connect)
    
    // NET.SEND(channel, peer, message$) - send packet
    r.Register("NET.SEND", "net", m.send)
    
    // NET.RECEIVE$() - get pending message or ""
    r.Register("NET.RECEIVE$", "net", m.receive)
    
    // NET.PEERCOUNT() - number of connected peers
    r.Register("NET.PEERCOUNT", "net", m.peerCount)
}
```

---

### Phase 6: Crash-Resilient VM - 1 week

**Goal**: Catch errors gracefully, developer console integration

#### 6.1 Defensive VM Execution

Update `vm/vm.go`:

```go
// Execute with panic recovery and error reporting
func (v *VM) ExecuteSafe(prog *opcode.Program) (err error) {
    // Set up recovery
    defer func() {
        if r := recover(); r != nil {
            // Capture stack trace
            stack := debug.Stack()
            
            // Create detailed error report
            err = &VMError{
                Type:    ErrorTypePanic,
                Message: fmt.Sprintf("VM panic: %v", r),
                Stack:   v.FormatCallStack(),
                GoStack: string(stack),
                Line:    v.currentLine(),
            }
            
            // Send to developer console if available
            if v.Console != nil {
                v.Console.LogError(err.(*VMError))
            }
        }
    }()
    
    // Normal execution
    return v.Execute(prog)
}

type VMError struct {
    Type    ErrorType
    Message string
    Stack   string
    GoStack string
    Line    int
}

func (e *VMError) Error() string {
    return fmt.Sprintf("[%s at line %d] %s\nStack:\n%s", 
        e.Type, e.Line, e.Message, e.Stack)
}
```

#### 6.2 Developer Console Integration

Create `runtime/console/developer_console.go`:

```go
package console

// DeveloperConsole provides in-game debugging UI
type DeveloperConsole struct {
    visible     bool
    history     []LogEntry
    commands    map[string]Command
    vm          *vm.VM
    
    // UI state
    scrollPos   int
    inputBuffer string
}

type LogEntry struct {
    Time    time.Time
    Level   LogLevel
    Message string
    Source  string  // "VM", "Engine", "Script"
}

func (dc *DeveloperConsole) LogError(err *vm.VMError) {
    dc.history = append(dc.history, LogEntry{
        Time:    time.Now(),
        Level:   LevelError,
        Message: err.Message,
        Source:  "VM",
    })
    
    // Auto-show console on critical errors
    if err.Type == vm.ErrorTypePanic {
        dc.visible = true
    }
}

func (dc *DeveloperConsole) Draw() {
    if !dc.visible {
        return
    }
    
    // Draw console overlay using raygui
    GUI.WINDOWBOX(10, 10, 780, 580, "Developer Console")
    
    // Show recent log entries
    for i, entry := range dc.getVisibleEntries() {
        color := logLevelColor(entry.Level)
        GUI.LABEL(20, 40+i*20, 760, 18, entry.Source+": "+entry.Message, color)
    }
}

func (dc *DeveloperConsole) Toggle() {
    dc.visible = !dc.visible
}
```

---

### Phase 7: LSP-Ready Symbol Table - 1 week

**Goal**: Export symbol information for IDE integration

#### 7.1 Symbol Table Export

Create `compiler/lsp/symbol_export.go`:

```go
package lsp

import (
    "encoding/json"
    "moonbasic/compiler/symtable"
)

// LSPSymbol represents a symbol for Language Server Protocol
type LSPSymbol struct {
    Name       string       `json:"name"`
    Type       string       `json:"type"`       // "function", "variable", "type"
    DataType   string       `json:"dataType"`   // "int", "float", "string", "bool"
    Location   SourceRange  `json:"location"`
    Scope      string       `json:"scope"`      // "global", "local", "param"
    Signature  string       `json:"signature,omitempty"`  // For functions
    Parameters []Parameter  `json:"parameters,omitempty"`
    Fields     []Field      `json:"fields,omitempty"`     // For types
}

type SourceRange struct {
    Start LinePos `json:"start"`
    End   LinePos `json:"end"`
}

type LinePos struct {
    Line int `json:"line"`
    Char int `json:"char"`
}

// ExportSymbolTable converts internal symbol table to LSP format
func ExportSymbolTable(st *symtable.Table) []LSPSymbol {
    var symbols []LSPSymbol
    
    // Export functions
    for name, fn := range st.Functions {
        symbols = append(symbols, LSPSymbol{
            Name:      name,
            Type:      "function",
            Location:  SourceRange{Start: LinePos{Line: fn.Line}},
            Scope:     "global",
            Signature: fn.Signature,
            Parameters: convertParams(fn.Params),
        })
    }
    
    // Export variables
    for name, sym := range st.Symbols {
        symbols = append(symbols, LSPSymbol{
            Name:     name,
            Type:     "variable",
            DataType: typeToString(sym.Type),
            Location: SourceRange{Start: LinePos{Line: sym.Line}},
            Scope:    scopeToString(sym.Kind),
        })
    }
    
    // Export types
    for name, td := range st.Types {
        symbols = append(symbols, LSPSymbol{
            Name:   name,
            Type:   "type",
            Scope:  "global",
            Fields: convertFields(td.Fields),
        })
    }
    
    return symbols
}

// WriteSymbolFile exports symbols to JSON for LSP consumption
func WriteSymbolFile(st *symtable.Table, path string) error {
    symbols := ExportSymbolTable(st)
    data, err := json.MarshalIndent(symbols, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}
```

#### 7.2 Compilation Database

```go
// CompileDatabase tracks build information for LSP
type CompileDatabase struct {
    Version    string            `json:"version"`
    Entries    []CompileEntry  `json:"entries"`
    Symbols    []LSPSymbol       `json:"symbols"`
}

type CompileEntry struct {
    Source    string   `json:"source"`
    Object    string   `json:"object"`
    Command   string   `json:"command"`
    Includes  []string `json:"includes"`
}

// GenerateCompileDB creates compile_commands.json for clangd/LSP
func GenerateCompileDB(files []string, outputPath string) error {
    db := CompileDatabase{Version: "1.0"}
    
    for _, file := range files {
        // Determine dependencies
        includes := findIncludes(file)
        
        db.Entries = append(db.Entries, CompileEntry{
            Source:   file,
            Object:   strings.Replace(file, ".mb", ".mbc", 1),
            Command:  fmt.Sprintf("moonbasic --compile %s", file),
            Includes: includes,
        })
    }
    
    data, _ := json.MarshalIndent(db, "", "  ")
    return os.WriteFile(outputPath, data, 0644)
}
```

---

## Success Metrics & Testing

### Performance Benchmarks

Create `tests/benchmarks/performance_test.go`:

```go
func BenchmarkDrawLoop(b *testing.B) {
    // Initialize window and renderer
    win := window.New()
    win.Open(3840, 2160, "Benchmark")  // 4K
    
    b.ReportAllocs()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        // Simulate 1000 draw calls
        for j := 0; j < 1000; j++ {
            DRAW.RECTANGLE(
                float64(j*10), float64(j*5), 
                50, 50, 
                255, 255, 255, 255,
            )
        }
        RENDER.FRAME()
    }
    
    // Verify zero allocations
    b.StopTimer()
    if b.Allocs() > 0 {
        b.Fatalf("Expected 0 allocations, got %d", b.Allocs())
    }
}

func BenchmarkPhysics1000Bodies(b *testing.B) {
    // Create 1000 physics bodies
    world := physics3d.NewWorld()
    bodies := make([]physics3d.Body, 1000)
    
    for i := range bodies {
        bodies[i] = world.CreateBody(physics3d.BoxShape{...})
    }
    
    b.ReportAllocs()
    b.ResetTimer()
    
    dt := 1.0 / 144.0  // 144Hz
    
    for i := 0; i < b.N; i++ {
        world.Step(dt)
        
        // Batch update all bodies
        for _, body := range bodies {
            body.Update()
        }
    }
}
```

### Memory Leak Detection

```go
func TestMemoryLeaks(t *testing.T) {
    // Start with clean heap
    runtime.GC()
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    // Run game loop simulation
    for i := 0; i < 10000; i++ {
        // Create and destroy resources
        tex := LOADTEXTURE("test.png")
        model := CREATEMESH()
        FREETEXTURE(tex)
        FREEMESH(model)
    }
    
    // Force GC and check
    runtime.GC()
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    
    // Should have same allocation count
    if m2.HeapObjects > m1.HeapObjects+10 {
        t.Fatalf("Memory leak detected: %d objects remaining", 
            m2.HeapObjects-m1.HeapObjects)
    }
}
```

---

## Integration Checklist

### Compiler Changes
- [ ] Two-pass symbol table builder for implicit declarations
- [ ] Type inference engine
- [ ] Remove VAR keyword requirement
- [ ] LSP symbol export

### VM Changes  
- [ ] Crash recovery with defer/recover
- [ ] Developer console integration
- [ ] Allocation-free hot paths

### Runtime Changes
- [ ] High-DPI window support
- [ ] Batch rendering system
- [ ] Post-processing effects (MSAA, Bloom, SSAO)
- [ ] Shared memory buffers for physics
- [ ] Fixed-timestep physics with interpolation
- [ ] Memory pool for physics bodies
- [ ] Finalizer-based resource management
- [ ] ENet channel system
- [ ] Thread-safe networking

### Tooling
- [ ] pprof integration
- [ ] Memory leak detection
- [ ] Performance benchmarks
- [ ] LSP compile database generation

---

## Estimated Timeline

| Phase | Duration | Key Deliverables |
|-------|----------|------------------|
| 1: Modern Syntax | 2 weeks | Implicit declaration, type inference |
| 2: High-Perf Rendering | 3 weeks | Zero-allocation draw loop, 4K support |
| 3: Physics & CGO | 3 weeks | Shared buffers, fixed timestep |
| 4: Resource Safety | 2 weeks | Finalizers, leak detection |
| 5: Networking | 2 weeks | ENet channels, thread isolation |
| 6: VM Resilience | 1 week | Crash recovery, dev console |
| 7: LSP Support | 1 week | Symbol export, compile DB |
| **Total** | **14 weeks** | High-fidelity engine complete |

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         moonBASIC Engine                         │
├─────────────────────────────────────────────────────────────────┤
│  Compiler Layer                                                  │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐          │
│  │ Lexer    │→│ Parser   │→│ Semantic │→│ CodeGen  │          │
│  │ (tokens) │ │ (AST)    │ │ Analysis │ │ (bytecode│          │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘          │
│       ↓              ↓              ↓              ↓             │
│  ┌──────────────────────────────────────────────────────┐      │
│  │ Symbol Table (2-pass) + Type Inference             │      │
│  └──────────────────────────────────────────────────────┘      │
│       ↓                                                        │
│  ┌──────────────────────────────────────────────────────┐      │
│  │ Register-Based Bytecode (IR v3)                      │      │
│  │ 8-byte instructions: Op/Dst/SrcA/SrcB/Operand        │      │
│  └──────────────────────────────────────────────────────┘      │
├─────────────────────────────────────────────────────────────────┤
│  VM Layer (Crash-Resilient)                                     │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐           │
│  │ ExecuteSafe  │→│ defer/recover│→│ Dev Console  │           │
│  │ (panic guard)│ │ (error catch)│ │ (reporting)  │           │
│  └──────────────┘ └──────────────┘ └──────────────┘           │
│  ┌──────────────┐ ┌──────────────┐                             │
│  │ Call Stack   │ │ Register     │                             │
│  │ (frames)     │ │ Allocation   │                             │
│  └──────────────┘ └──────────────┘                             │
├─────────────────────────────────────────────────────────────────┤
│  Runtime Layer                                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │ Window      │ │ Render      │ │ Physics     │               │
│  │ (High-DPI)  │ │ (Batch)     │ │ (Jolt/Box2D)│               │
│  └─────────────┘ └─────────────┘ └─────────────┘               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │ Heap        │ │ Network     │ │ Audio       │               │
│  │ (Finalizers)│ │ (ENet)      │ │ (Raylib)    │               │
│  └─────────────┘ └─────────────┘ └─────────────┘               │
├─────────────────────────────────────────────────────────────────┤
│  CGO Fast-Path                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │ Shared Memory Buffers (C.malloc)                       │   │
│  │ Pre-allocated vertex buffers, physics state arrays   │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │ Batch Operations (single CGO crossing)                 │   │
│  │ batch_render(), batch_update_bodies(), etc.            │   │
│  └──────────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────────┤
│  External Libraries                                              │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐          │
│  │ Raylib   │ │ Jolt     │ │ Box2D    │ │ ENet     │          │
│  │ 5.5      │ │ Physics  │ │ Physics  │ │ Network  │          │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

This roadmap provides a complete technical blueprint for transforming moonBASIC into a high-fidelity, modern game engine capable of 4K/144Hz+ performance with thousands of physics objects.
