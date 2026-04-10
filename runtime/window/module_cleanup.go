package window

// EnqueueCleanup adds a function to be run on the main thread at the start of the next frame.
// Safe to call from any goroutine (e.g. Go finalizers).
func (m *Module) EnqueueCleanup(fn func()) {
	m.cleanupMu.Lock()
	m.cleanupQueue = append(m.cleanupQueue, fn)
	m.cleanupMu.Unlock()
}

func (m *Module) drainCleanupQueue() {
	m.cleanupMu.Lock()
	if len(m.cleanupQueue) == 0 {
		m.cleanupMu.Unlock()
		return
	}
	q := m.cleanupQueue
	m.cleanupQueue = nil
	m.cleanupMu.Unlock()
	for _, fn := range q {
		if fn != nil {
			fn()
		}
	}
}
