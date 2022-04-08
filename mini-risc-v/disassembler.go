package main

import (
	"fmt"
)

func disassembleInstr(instr uint32) string {
	data := decodeInstr(instr)
	if data == nil {
		return "unknown"
	}
	if data.opcode == 3 {
		if data.funct3 == 7 {
			return "illegal load"
		}
		var read uint8
		switch data.funct3 & 3 {
		case 0:
			read = 'b'
		case 1:
			read = 'h'
		case 2:
			read = 'w'
		case 3:
			read = 'd'
		}
		if data.funct3&4 != 0 {
			return fmt.Sprintf("l%cu x%d, %d(x%d)", read, data.rd, signExtend(data.imm, 12), data.rs1)
		}
		return fmt.Sprintf("l%c x%d, %d(x%d)", read, data.rd, signExtend(data.imm, 12), data.rs1)
	} else if data.format == ITYPE_U {
		if data.opcode&0x20 == 0 {
			return fmt.Sprintf("auipc %d", signExtend(data.imm<<12, 32))
		} else {
			return fmt.Sprintf("lui x%d, %d", data.rd, signExtend(data.imm<<12, 32))
		}
	} else if data.opcode&0x77 == 0x67 {
		if data.opcode&8 != 0 {
			return fmt.Sprintf("jal x%d, %d", data.rd, signExtend(data.imm, 12))
		}
		return fmt.Sprintf("jalr x%d, %d(x%d)", data.rd, signExtend(data.imm, 12), data.rs1)
	} else if data.format == ITYPE_SB {
		var read string
		switch data.funct3 {
		case 0:
			read = "eq"
		case 1:
			read = "ne"
		case 2:
			fallthrough
		case 3:
			return "illegal branch"
		case 4:
			read = "lt"
		case 5:
			read = "ge"
		case 6:
			read = "ltu"
		case 7:
			read = "geu"
		}
		return fmt.Sprintf("b%s x%d, x%d, %d", read, data.rs1, data.rs2, data.imm)
	} else if data.format == ITYPE_S {
		var read uint8
		switch data.funct3 {
		case 0:
			read = 'b'
		case 1:
			read = 'h'
		case 2:
			read = 'w'
		case 3:
			read = 'd'
		default:
			return "illegal store"
		}
		return fmt.Sprintf("s%c x%d, %d(x%d)", read, data.rs2, signExtend(data.imm, 12), data.rs1)
	} else if data.opcode&0x57 == 0x13 {
		word := data.opcode&0x1f == 0x1b
		reg := data.opcode&0x20 != 0
		unsigned := false
		var read string
		switch data.funct3 {
		case 0:
			if data.funct7 != 0 {
				read = "sub"
			} else {
				read = "add"
			}
		case 1:
			read = "sll"
		case 2:
			read = "slt"
		case 3:
			read = "sltu"
		case 4:
			read = "xor"
		case 5:
			if data.funct7 != 0 || data.imm&0x400 != 0 {
				read = "srl"
			} else {
				read = "sra"
			}
		case 6:
			read = "or"
		case 7:
			read = "and"
		}
		if !reg {
			read += "i"
		}
		if unsigned {
			read += "u"
		}
		if word {
			read += "w"
		}
		if reg {
			return fmt.Sprintf("%s x%d, x%d, x%d", read, data.rd, data.rs1, data.rs2)
		} else {
			return fmt.Sprintf("%s x%d, x%d, %d", read, data.rd, data.rs1, signExtend(data.imm, 12))
		}
	}
	return "illegal instruction"
}
