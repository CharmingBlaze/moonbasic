//go:build (cgo || (windows && !cgo)) && (!windows || !gopls_stub)

package terrain

import (
	"runtime"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type chunkMeshJob struct {
	cx, cz int
	prep   *heightmapPrep
}

func (t *TerrainObject) initMeshJobChannel() {
	if t.meshJobs == nil {
		t.meshJobs = make(chan chunkMeshJob, 256)
	}
}

func waitForPendingMeshWorkers(t *TerrainObject) {
	deadline := time.Now().Add(8 * time.Second)
	for t.meshJobsInflight.Load() > 0 && time.Now().Before(deadline) {
		t.drainMeshBuildJobs()
		runtime.Gosched()
	}
	t.drainMeshBuildJobs()
}

// drainMeshBuildJobs applies completed CPU heightmap jobs on the main thread (GenMeshHeightmap + GPU).
// Call from WORLD.UPDATE / streaming tick before scheduling new work.
func (t *TerrainObject) drainMeshBuildJobs() {
	if t.meshJobs == nil {
		return
	}
	for {
		select {
		case job := <-t.meshJobs:
			if t.freed {
				if idx := idx2(t, job.cx, job.cz); idx >= 0 && idx < len(t.Chunks) {
					t.Chunks[idx].PendingAsync = false
				}
				continue
			}
			idx := idx2(t, job.cx, job.cz)
			if idx < 0 || idx >= len(t.Chunks) {
				continue
			}
			ch := &t.Chunks[idx]
			if !ch.PendingAsync {
				continue
			}
			if job.prep == nil {
				ch.PendingAsync = false
				ch.Dirty = true
				rl.PollInputEvents()
				continue
			}
			t.applyHeightmapPrep(job.cx, job.cz, job.prep)
			ch.PendingAsync = false
			rl.PollInputEvents()
		default:
			return
		}
	}
}

func (t *TerrainObject) trySendMeshJob(job chunkMeshJob) {
	idx := idx2(t, job.cx, job.cz)
	if idx < 0 || idx >= len(t.Chunks) {
		return
	}
	ch := &t.Chunks[idx]
	if t.freed || t.meshJobs == nil {
		ch.PendingAsync = false
		ch.Dirty = true
		return
	}
	select {
	case t.meshJobs <- job:
	default:
		ch.PendingAsync = false
		ch.Dirty = true
	}
}

// shutdownMeshJobQueue waits for prep goroutines, drains and closes meshJobs. Call after t.freed = true.
func (t *TerrainObject) shutdownMeshJobQueue() {
	if t.meshJobs == nil {
		return
	}
	waitForPendingMeshWorkers(t)
	for {
		select {
		case job := <-t.meshJobs:
			if idx := idx2(t, job.cx, job.cz); idx >= 0 && idx < len(t.Chunks) {
				t.Chunks[idx].PendingAsync = false
			}
			_ = job
		default:
			close(t.meshJobs)
			t.meshJobs = nil
			return
		}
	}
}

func (t *TerrainObject) enqueueAsyncChunkMeshBuild(cx, cz int) {
	if t.freed {
		return
	}
	snap, ok := t.snapshotChunkHeights(cx, cz)
	if !ok {
		return
	}
	t.initMeshJobChannel()
	idx := idx2(t, cx, cz)
	if idx < 0 || idx >= len(t.Chunks) {
		return
	}
	ch := &t.Chunks[idx]
	if ch.PendingAsync {
		return
	}
	ch.PendingAsync = true
	go func(s *chunkHeightSnapshot, icx, icz int) {
		t.meshJobsInflight.Add(1)
		defer t.meshJobsInflight.Add(-1)
		prep, ok := buildHeightmapPrepFromSnapshot(s)
		if !ok {
			t.trySendMeshJob(chunkMeshJob{cx: icx, cz: icz, prep: nil})
			return
		}
		t.trySendMeshJob(chunkMeshJob{cx: icx, cz: icz, prep: prep})
	}(snap, cx, cz)
}
