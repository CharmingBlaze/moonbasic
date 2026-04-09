// Package lineprof defines the optional VM line profiler interface without importing the VM or runtime.
package lineprof

// LineProfiler receives per-source-line instruction counts during VM execution.
type LineProfiler interface {
	RecordLine(line int)
}
