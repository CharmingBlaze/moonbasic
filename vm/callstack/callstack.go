// Package callstack implements the VM execution context for function calls.
package callstack

import (
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

// Frame is one activation record on the call stack.
type Frame struct {
	Chunk     *opcode.Chunk // The code segment being executed
	IP        int           // Instruction pointer
	Registers [256]value.Value // Register file (R0-R255)
	ReturnReg uint8         // Target register in the CALLER frame to receive results
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
func (s *Stack) Push(c *opcode.Chunk, ip int, retReg uint8) {
	s.Frames = append(s.Frames, Frame{
		Chunk:     c,
		IP:        ip,
		ReturnReg: retReg,
	})
}

// Pop removes the top frame from the call stack.
func (s *Stack) Pop() Frame {
	if len(s.Frames) == 0 {
		return Frame{}
	}
	f := s.Frames[len(s.Frames)-1]
	s.Frames = s.Frames[:len(s.Frames)-1]
	return f
}

// Range iterates over all frames in the stack, starting from the oldest.
func (s *Stack) Range(fn func(*Frame)) {
	for i := range s.Frames {
		fn(&s.Frames[i])
	}
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
