package main

// #include "memory.h"
import "C"

const (
	REG_ZERO = iota
	REG_RA
	REG_SP
	REG_GP
	REG_TP
	REG_T0
	REG_T1
	REG_T2
	REG_S0
	REG_S1
	REG_A0
	REG_A1
	REG_A2
	REG_A3
	REG_A4
	REG_A5
	REG_A6
	REG_A7
	REG_S2
	REG_S3
	REG_S4
	REG_S5
	REG_S6
	REG_S7
	REG_S8
	REG_S9
	REG_S10
	REG_S11
	REG_T3
	REG_T4
	REG_T5
	REG_T6
	REG_FP = 8
	REG_PC = 32
)

type InstrStage2 struct {
	InstrData

	rs1Out, rs2Out int64
	// pc             int64
}
type InstrStage3 struct {
	InstrStage2

	ALUSrc bool
	ALUOp  string
	ALUOut int64
}
type InstrStage4 struct {
	InstrData

	ALUOut int64
	rs2Out int64
	flags  uint8
	memOut int64
	// pc     int64
}
type InstrStage5 struct {
	InstrData

	memToReg  uint8
	writeData int64
	// pc        int64
}

const (
	MEM_READ  = 1
	MEM_WRITE = 2
	PCSRC_4   = 4
	PCSRC_J   = 8
)

const (
	M2R_ALU = iota
	M2R_MEM
	M2R_PC
	M2R_NONE
)

func initExec() {
	if !C.init() {
		panic("Allocation of memory failed")
	}
}

func cleanupExec() {
	C.cleanup()
}

func signExtend(read uint32, pos uint8) int64 {
	if pos == 0 {
		pos = 8
	}
	out := uint64(read)
	if get_bit(read, uint(pos-1)) != 0 {
		out |= (0xffffffffffffffff << pos)
	}
	return int64(out)
}

func executeInstr(instr uint32, registers *[33]int64) {
	data := decodeInstr(instr)
	if data == nil {
		return
	}
	if data.opcode == 3 {
		if data.funct3 == 7 {
			return
		}
		var read C.int64_t
		addr := signExtend(data.imm, 12) + registers[data.rs1]
		if addr < 0x10000000 || addr > 0x3fffffffff {
			return
		}
		switch data.funct3 & 3 {
		case 0:
			if !C.readByte(C.int64_t(addr), &read) {
				return
			}
		case 1:
			if !C.readHalf(C.int64_t(addr), &read) {
				return
			}
		case 2:
			if !C.readWord(C.int64_t(addr), &read) {
				return
			}
		case 3:
			if !C.readDouble(C.int64_t(addr), &read) {
				return
			}
		}
		if data.funct3&4 != 0 {
			read &= (1 << (8 << (data.funct3 & 3))) - 1
		}
		registers[data.rd] = int64(read)
	} else if data.format == ITYPE_U {
		if data.opcode&0x20 == 0 {
			data.rd = REG_PC
		} else {
			registers[data.rd] = 0
		}
		registers[data.rd] |= signExtend(data.imm<<12, 32)
	} else if data.opcode&0x77 == 0x67 {
		registers[data.rd] = registers[REG_PC] + 4
		if data.opcode&8 != 0 {
			data.rs1 = REG_PC
		}
		registers[REG_PC] = signExtend(data.imm, 12) + registers[data.rs1]
	} else if data.format == ITYPE_SB {
		yes := false
		switch data.funct3 {
		case 0:
			yes = (registers[data.rs1] == registers[data.rs2])
		case 1:
			yes = (registers[data.rs1] != registers[data.rs2])
		case 2:
			fallthrough
		case 3:
			return
		case 4:
			yes = (registers[data.rs1] < registers[data.rs2])
		case 5:
			yes = (registers[data.rs1] >= registers[data.rs2])
		case 6:
			yes = (uint64(registers[data.rs1]) < uint64(registers[data.rs2]))
		case 7:
			yes = (uint64(registers[data.rs1]) >= uint64(registers[data.rs2]))
		}
		if yes {
			registers[REG_PC] += signExtend(data.imm, 13)
		}
	} else if data.format == ITYPE_S {
		addr := signExtend(data.imm, 12) + registers[data.rs1]
		if addr < 0x10000000 || addr > 0x3fffffffff {
			return
		}
		switch data.funct3 {
		case 0:
			if !C.writeByte(C.int64_t(addr), C.int8_t(registers[data.rs2])) {
				return
			}
		case 1:
			if !C.writeHalf(C.int64_t(addr), C.int16_t(registers[data.rs2])) {
				return
			}
		case 2:
			if !C.writeWord(C.int64_t(addr), C.int32_t(registers[data.rs2])) {
				return
			}
		case 3:
			if !C.writeDouble(C.int64_t(addr), C.int64_t(registers[data.rs2])) {
				return
			}
		default:
			return
		}
	} else if data.opcode&0x57 == 0x13 {
		var input1, input2 int64
		input1 = registers[data.rs1]
		if data.opcode&0x1f == 0x1b {
			input1 &= 0xffffffff
		}
		if data.opcode&0x20 != 0 {
			input2 = registers[data.rs2]
		} else {
			input2 = signExtend(data.imm, 12)
		}
		switch data.funct3 {
		case 0:
			if data.funct7 != 0 {
				input2 = 0 - input2
			}
			registers[data.rd] = input1 + input2
		case 1:
			registers[data.rd] = input1 << (input2 & 0x3f)
		case 2:
			if input1 < input2 {
				registers[data.rd] = 1
			} else {
				registers[data.rd] = 0
			}
		case 3:
			if uint64(input1) < uint64(input2) {
				registers[data.rd] = 1
			} else {
				registers[data.rd] = 0
			}
		case 4:
			registers[data.rd] = input1 ^ input2
		case 5:
			if data.funct7 != 0 || data.imm&0x400 != 0 {
				registers[data.rd] = int64(uint64(input1) >> uint64(input2&0x3f))
			} else {
				registers[data.rd] = input1 >> (input2 & 0x3f)
			}
		case 6:
			registers[data.rd] = input1 | input2
		case 7:
			registers[data.rd] = input1 & input2
		}
		if data.opcode&0x1f == 0x1b {
			registers[data.rd] = signExtend(uint32(registers[data.rd]), 32)
		}
	}
}

func execStage2(instr uint32, registers *[33]int64) *InstrStage2 {
	if instr == 0 {
		return nil
	}
	out := &InstrStage2{InstrData: *decodeInstr(instr)}
	if out.format == ITYPE_U {
		out.rs1 = REG_ZERO
		if out.opcode&0x20 == 0 {
			out.rd = REG_PC
		}
	} else if out.opcode&0x77 == 0x67 {
		if out.opcode&8 != 0 {
			out.rs1 = REG_PC
		}
	}
	out.rs1Out = registers[out.rs1]
	out.rs2Out = registers[out.rs2]
	return out
}

func execStage3(instr *InstrStage2) *InstrStage3 {
	if instr == nil {
		return nil
	}
	out := &InstrStage3{InstrStage2: *instr}
	if out.opcode == 3 || out.format == ITYPE_S {
		out.ALUOp = "add"
		out.ALUOut = signExtend(out.imm, 12) + out.rs1Out
		out.ALUSrc = true
	} else if out.format == ITYPE_U {
		out.ALUOp = "or"
		out.ALUSrc = true
		out.ALUOut = signExtend(out.imm<<12, 32) | (out.rs1Out & int64(MASK_LOWER(12)))
	} else if out.opcode&0x77 == 0x67 {
		out.ALUOut = signExtend(out.imm, 12) + out.rs1Out
		out.ALUOp = "add"
		out.ALUSrc = true
	} else if out.format == ITYPE_SB {
		switch out.funct3 & 6 {
		case 0:
			out.ALUOp = "sub"
			out.ALUOut = (out.rs1Out - out.rs2Out)
		case 2:
			return nil
		case 4:
			out.ALUOp = "less"
			if out.rs1Out >= out.rs2Out {
				out.ALUOut = 1
			}
		case 6:
			out.ALUOp = "less unsigned"
			if uint64(out.rs1Out) >= uint64(out.rs2Out) {
				out.ALUOut = 1
			}
		}
	} else if out.opcode&0x57 == 0x13 {
		var input2 int64
		if out.opcode&0x1f == 0x1b {
			out.rs1Out &= 0xffffffff
		}
		if out.opcode&0x20 != 0 {
			input2 = out.rs2Out
		} else {
			input2 = signExtend(out.imm, 12)
			out.ALUSrc = true
		}
		switch out.funct3 {
		case 0:
			if out.funct7 != 0 {
				out.ALUOp = "sub"
				input2 = 0 - input2
			} else {
				out.ALUOp = "add"
			}
			out.ALUOut = out.rs1Out + input2
		case 1:
			out.ALUOp = "sll"
			out.ALUOut = out.rs1Out << (input2 & 0x3f)
		case 2:
			out.ALUOp = "less"
			if out.rs1Out < input2 {
				out.ALUOut = 1
			} else {
				out.ALUOut = 0
			}
		case 3:
			out.ALUOp = "less unsigned"
			if uint64(out.rs1Out) < uint64(input2) {
				out.ALUOut = 1
			} else {
				out.ALUOut = 0
			}
		case 4:
			out.ALUOp = "xor"
			out.ALUOut = out.rs1Out ^ input2
		case 5:
			if out.funct7 != 0 || out.imm&0x400 != 0 {
				out.ALUOp = "srl"
				out.ALUOut = int64(uint64(out.rs1Out) >> uint64(input2&0x3f))
			} else {
				out.ALUOp = "sra"
				out.ALUOut = out.rs1Out >> (input2 & 0x3f)
			}
		case 6:
			out.ALUOp = "or"
			out.ALUOut = out.rs1Out | input2
		case 7:
			out.ALUOp = "and"
			out.ALUOut = out.rs1Out & input2
		}
		if out.opcode&0x1f == 0x1b {
			out.ALUOp += "w"
		}
	}
	return out
}

func execStage4(instr *InstrStage3) *InstrStage4 {
	if instr == nil {
		return nil
	}
	out := &InstrStage4{InstrData: instr.InstrData, rs2Out: instr.rs2Out, ALUOut: instr.ALUOut}
	if out.opcode == 3 {
		if out.funct3 == 7 {
			return nil
		}
		out.ALUOut = instr.ALUOut
		out.flags = MEM_READ
		if out.ALUOut < 0x10000000 || out.ALUOut > 0x3fffffffff {
			return nil
		}
		var read C.int64_t
		switch out.funct3 & 3 {
		case 0:
			if !C.readByte(C.int64_t(out.ALUOut), &read) {
				return nil
			}
		case 1:
			if !C.readHalf(C.int64_t(out.ALUOut), &read) {
				return nil
			}
		case 2:
			if !C.readWord(C.int64_t(out.ALUOut), &read) {
				return nil
			}
		case 3:
			if !C.readDouble(C.int64_t(out.ALUOut), &read) {
				return nil
			}
		}
		out.memOut = int64(read)
		if out.funct3&4 != 0 {
			out.memOut &= (1 << (8 << (out.funct3 & 3))) - 1
		}
	} else if out.format == ITYPE_U {
		if out.rd == REG_PC {
			out.flags = PCSRC_J
		}
	} else if out.opcode&0x77 == 0x67 {
		out.flags = PCSRC_J
	} else if out.format == ITYPE_SB {
		yes := instr.ALUOut == 0
		if out.funct3&1 == 1 {
			yes = !yes
		}
		if yes {
			out.flags = PCSRC_J
		}
	} else if out.format == ITYPE_S {
		if out.ALUOut < 0x10000000 || out.ALUOut > 0x3fffffffff {
			return nil
		}
		out.flags = MEM_WRITE
		switch out.funct3 {
		case 0:
			if !C.writeByte(C.int64_t(out.ALUOut), C.int8_t(out.rs2Out)) {
				return nil
			}
		case 1:
			if !C.writeHalf(C.int64_t(out.ALUOut), C.int16_t(out.rs2Out)) {
				return nil
			}
		case 2:
			if !C.writeWord(C.int64_t(out.ALUOut), C.int32_t(out.rs2Out)) {
				return nil
			}
		case 3:
			if !C.writeDouble(C.int64_t(out.ALUOut), C.int64_t(out.rs2Out)) {
				return nil
			}
		default:
			return nil
		}
	}
	if out.flags&PCSRC_J == 0 {
		out.flags |= PCSRC_4
	}
	return out
}

func execStage5(instr *InstrStage4, registers *[33]int64) *InstrStage5 {
	if instr == nil {
		return nil
	}
	out := &InstrStage5{InstrData: instr.InstrData}
	registers[REG_PC] += 4
	if out.opcode == 3 {
		out.writeData = instr.memOut
		out.memToReg = M2R_MEM
	} else if out.format == ITYPE_U {
		out.memToReg = M2R_NONE
	} else if out.opcode&0x77 == 0x67 {
		out.writeData = registers[REG_PC]
		out.memToReg = M2R_PC
	} else if out.format == ITYPE_SB {
		out.memToReg = M2R_NONE
		instr.ALUOut = registers[REG_PC] + signExtend(instr.imm, 12)
	} else if out.format == ITYPE_S {
		out.memToReg = M2R_NONE
	} else {
		out.writeData = instr.ALUOut
		out.memToReg = M2R_ALU
	}
	if instr.flags&PCSRC_J != 0 {
		registers[REG_PC] = instr.ALUOut
	}
	if out.memToReg != M2R_NONE {
		registers[out.rd] = out.writeData
	}
	return out
}
