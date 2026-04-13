package entityspatial

// MaxEntitySpatialIndex is the exclusive upper bound for numeric entity indices used with
// ENTITY.X/Y/Z/P/W/YAW/R spatial macros (SoA fast path). Must stay in sync with
// runtime.MaxEntitySpatialIndex in ../../runtime/runtime.go and vm checks in vm_dispatch.go.
const MaxEntitySpatialIndex int64 = 1 << 24 // 16_777_216
