// Package opt applies bytecode-level optimisation passes to compiled chunks (IR v3).
package opt

import (
	"moonbasic/vm/opcode"
)

// HoistFromLoops is a placeholder for conservative pure-call hoisting (requires sound purity data).
func HoistFromLoops(_ *opcode.Chunk) {}

// OptimizeProgram runs optimisation passes on main and all function chunks.
func OptimizeProgram(p *opcode.Program) {
	if p == nil || p.Main == nil {
		return
	}
	OptimizeChunk(p.Main)
	for _, ch := range p.Functions {
		if ch != nil {
			OptimizeChunk(ch)
		}
	}
}

// OptimizeChunk applies peephole cleanup and jump threading for direct JUMP chains.
func OptimizeChunk(ch *opcode.Chunk) {
	if ch == nil || len(ch.Instructions) == 0 {
		return
	}
	// TODO: implement register-based peephole (e.g. redundant MOVE removal)
	threadJumps(ch)
}

func threadJumps(ch *opcode.Chunk) {
	ins := ch.Instructions
	target := func(ip int32) int32 {
		for int(ip) >= 0 && int(ip) < len(ins) && ins[ip].Op == opcode.OpJump {
			nxt := ins[ip].Operand
			if nxt == ip {
				break
			}
			ip = nxt
		}
		return ip
	}
	for i := range ins {
		switch ins[i].Op {
		case opcode.OpJump:
			ins[i].Operand = target(ins[i].Operand)
		case opcode.OpJumpIfFalse, opcode.OpJumpIfTrue:
			ins[i].Operand = target(ins[i].Operand)
		}
	}
}
