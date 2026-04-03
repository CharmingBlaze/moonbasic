// Package opt applies bytecode-level optimisation passes to compiled chunks (IR v2).
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
// (Full unreachable-block DCE would require jump target relabeling; not done here.)
func OptimizeChunk(ch *opcode.Chunk) {
	if ch == nil || len(ch.Instructions) == 0 {
		return
	}
	peephole(ch)
	threadJumps(ch)
}

func peephole(ch *opcode.Chunk) {
	ins := ch.Instructions
	lines := ch.SourceLines
	outI := make([]opcode.Instruction, 0, len(ins))
	outL := make([]int32, 0, len(lines))
	for i := 0; i < len(ins); i++ {
		if i+1 < len(ins) {
			a, b := ins[i], ins[i+1]
			// PUSH immediately POP
			switch a.Op {
			case opcode.OpPushInt, opcode.OpPushFloat, opcode.OpPushString, opcode.OpPushBool, opcode.OpPushNull:
				if b.Op == opcode.OpPop {
					i++
					continue
				}
			}
		}
		outI = append(outI, ins[i])
		if i < len(lines) {
			outL = append(outL, lines[i])
		} else {
			outL = append(outL, 0)
		}
	}
	ch.Instructions = outI
	ch.SourceLines = outL
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
