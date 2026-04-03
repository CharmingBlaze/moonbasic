// Package callstack implements the VM execution context for function calls.
package callstack

import (
	"moonbasic/vm/opcode"
)

// Frame is one activation record on the call stack.
type Frame struct {
	Chunk     *opcode.Chunk // The code segment being executed
	IP        int           // Instruction pointer
	StackBase int           // Value stack offset where this frame's locals start
}

// Stack is the caller's execution environment.
type Stack struct {
	Frames []Frame
}

// New creates a new call stack.
func New() *Stack {
	return &Stack{
		Frames: make([]Frame, 0, 64), // Default depth limit of 64
	}
}

// Push adds a new frame to the call stack.
func (s *Stack) Push(c *opcode.Chunk, ip, base int) {
	s.Frames = append(s.Frames, Frame{
		Chunk:     c,
		IP:        ip,
		StackBase: base,
	})
}

// Pop removes the innermost frame.
func (s *Stack) Pop() Frame {
	if len(s.Frames) == 0 {
		return Frame{}
	}
	f := s.Frames[len(s.Frames)-1]
	s.Frames = s.Frames[:len(s.Frames)-1]
	return f
}

// Top returns the current execution frame.
func (s *Stack) Top() *Frame {
	if len(s.Frames) == 0 {
		return nil
	}
	return &s.Frames[len(s.Frames)-1]
}

// Depth returns the current call stack depth.
func (s *Stack) Depth() int {
	return len(s.Frames)
}

// FramesCopy returns a shallow copy of frames (outermost first) for stack traces.
func (s *Stack) FramesCopy() []Frame {
	out := make([]Frame, len(s.Frames))
	copy(out, s.Frames)
	return out
}
