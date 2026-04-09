package mbphysics3d

// afterPhysicsMatrixSync runs after Jolt bodies are written into the shared matrix buffer
// (see jolt_linux syncSharedBuffers). mbentity registers this to pull poses into scene entities.
var afterPhysicsMatrixSync func()

// SetAfterPhysicsMatrixSync registers a callback invoked after each physics step's matrix
// sync when Jolt is active. Pass nil to clear.
func SetAfterPhysicsMatrixSync(fn func()) {
	afterPhysicsMatrixSync = fn
}
