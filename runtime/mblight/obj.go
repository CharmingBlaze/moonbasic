package mblight

import "moonbasic/vm/heap"

// lightObj is a lightweight CPU-side light description for PBR + shadow mapping.
// No native allocations — safe to free repeatedly (idempotent).
type lightObj struct {
	release heap.ReleaseOnce
	self    heap.Handle

	kind      string
	r, g, b   float32
	colA      float32 // multiplies diffuse RGB (0–1 after normalization; default 1)
	intensity float32
	dirX      float32
	dirY      float32
	dirZ      float32
	shadow    bool

	// Shadow ortho camera looks at this world point (directional shadow frustum).
	targetX float32
	targetY float32
	targetZ float32

	// Multiplier for depth bias in the PBR shadow term (default 1).
	shadowBiasK float32

	// Point / spot parameters (stored for API completeness; main PBR path uses directional sun).
	posX float32
	posY float32
	posZ float32

	innerConeDeg float32
	outerConeDeg float32
	rangeDist    float32

	enabled bool

	// Optional scene-graph parent (entity#); not yet applied by the light system — reserved for follow / attachment.
	parentEntID int64
}

func (o *lightObj) TypeName() string { return "Light" }

func (o *lightObj) TypeTag() uint16 { return heap.TagLight }

func (o *lightObj) Free() {
	unregisterPointFollow(o.self)
	o.release.Do(func() {
		shadowMu.Lock()
		if shadowCasterHandle == o.self {
			shadowCasterHandle = 0
		}
		shadowMu.Unlock()
	})
}
