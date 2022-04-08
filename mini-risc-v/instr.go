package main

import "fmt"

const (
	ITYPE_R = iota
	ITYPE_I
	ITYPE_S
	ITYPE_SB
	ITYPE_U
	ITYPE_UJ
)

type InstrType uint8

type InstrData struct {
	format                               InstrType
	imm                                  uint32
	rd, rs1, rs2, opcode, funct3, funct7 uint8
}

func MASK_LOWER(x uint) uint32 {
	return (1 << x) - 1
}

func get_bit(data uint32, n uint) uint32 {
	return (data & (1 << n))
}

func encodeInstr(data *InstrData) uint32 {
	if data == nil {
		return 0
	}
	out := uint32(data.opcode)
	var format InstrType
	if (out&0x64) == 0 || out == 0x73 || out == 0x67 {
		format = ITYPE_I
	} else if out == 0x23 {
		format = ITYPE_S
	} else if out&0x44 == 0 {
		format = ITYPE_R
	} else if out == 0x6f {
		format = ITYPE_UJ
	} else if out == 0x63 {
		format = ITYPE_SB
	} else {
		format = ITYPE_U
	}
	out |= (uint32(data.rd) << 7)
	out |= (uint32(data.funct3) << 12)
	out |= (uint32(data.rs1) << 15)
	out |= (uint32(data.rs2) << 20)
	imm := data.imm
	switch format {
	case ITYPE_I:
		out &= MASK_LOWER(19)
		out |= (imm << 20)
	case ITYPE_R:
		out |= (uint32(data.funct7) << 25)
	case ITYPE_UJ:
		imm2 := imm
		imm = ((imm2 >> 12) & MASK_LOWER(8))
		if get_bit(imm2, 11) != 0 {
			imm |= 0x100
		}
		imm |= ((imm2 & (MASK_LOWER(11) - 1)) << 8)
		imm |= get_bit(imm2, 20) >> 1
		fallthrough
	case ITYPE_U:
		out &= MASK_LOWER(11)
		out |= (imm << 12)
	case ITYPE_SB:
		if get_bit(imm, 11) != 0 {
			imm |= 1
			imm ^= (1 << 11)
		}
		if get_bit(imm, 12) != 0 {
			imm |= (1 << 11)
		}
		fallthrough
	case ITYPE_S:
		out &= 0x01fff07f
		out |= ((imm & MASK_LOWER(5)) << 7)
		out |= ((imm & 0xfe0) << 20)
	default:
		return 0
	}
	return out
}

func decodeInstr(instr uint32) *InstrData {
	data := new(InstrData)
	opcode := instr & MASK_LOWER(7)
	data.opcode = uint8(opcode)
	var format InstrType
	if opcode&0x64 == 0 || opcode == 0x73 || opcode == 0x67 {
		format = ITYPE_I
	} else if opcode == 0x23 {
		format = ITYPE_S
	} else if opcode&0x44 == 0 {
		format = ITYPE_R
	} else if opcode == 0x6f {
		format = ITYPE_UJ
	} else if opcode == 0x63 {
		format = ITYPE_SB
	} else {
		format = ITYPE_U
	}
	data.format = format
	data.rd = uint8((instr >> 7) & MASK_LOWER(5))
	data.funct3 = uint8((instr >> 12) & MASK_LOWER(3))
	data.rs1 = uint8((instr >> 15) & MASK_LOWER(5))
	data.rs2 = uint8((instr >> 20) & MASK_LOWER(5))
	switch format {
	case ITYPE_R:
		data.funct7 = uint8(instr >> 25)
	case ITYPE_I:
		data.rs2 = 0
		data.imm = (instr >> 20)
	case ITYPE_S:
		fallthrough
	case ITYPE_SB:
		data.imm = uint32(data.rd)
		data.rd = 0
		data.imm |= (instr >> 20) & 0xfe0
	case ITYPE_U:
		fallthrough
	case ITYPE_UJ:
		data.rs1 = 0
		data.rs2 = 0
		data.funct3 = 0
		data.imm = (instr >> 12)
	default:
		return nil
	}
	if format == ITYPE_SB {
		if get_bit(data.imm, 11) != 0 {
			data.imm ^= (3 << 11)
		}
		if data.imm&1 != 0 {
			data.imm |= (1 << 11)
			data.imm ^= 1
		}
	} else if format == ITYPE_UJ {
		imm := ((data.imm >> 8) & (MASK_LOWER(11) - 1))
		if get_bit(data.imm, 8) != 0 {
			imm |= 0x800
		}
		imm |= ((data.imm & MASK_LOWER(8)) << 12)
		imm |= get_bit(data.imm, 19) << 1
		data.imm = imm
	}
	return data
}

func printInstrData(data *InstrData) {
	if data == nil {
		fmt.Println("Attempt to print NULL InstrData!")
		return
	}
	fmt.Printf("Opcode: 0x%x funct3: %d\nRD: x%d R1: x%d R2: x%d\nfunct7: %d imm: 0x%03x\n",
		data.opcode, data.funct3, data.rd, data.rs1, data.rs2, data.funct7,
		data.imm)
	switch data.format {
	case ITYPE_I:
		fmt.Println("I-type")
	case ITYPE_R:
		fmt.Println("R-type")
	case ITYPE_S:
		fmt.Println("S-type")
	case ITYPE_SB:
		fmt.Println("SB-type")
	case ITYPE_U:
		fmt.Println("U-type")
	case ITYPE_UJ:
		fmt.Println("UJ-type")
	default:
		fmt.Println("Unknown type")
	}
}
