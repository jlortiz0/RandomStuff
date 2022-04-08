package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func tokenizer(in rune) bool {
	return in == ' ' || in == '\t' || in == ',' || in == '(' || in == ')'
}

func getRegister(line string) uint8 {
	if len(line) == 0 {
		return 255
	}
	if line[0] == 'x' || line[0] == '$' {
		line = line[1:]
	}
	reg, err := strconv.ParseUint(line, 10, 5)
	if err == nil {
		return uint8(reg)
	}
	if len(line) == 1 {
		return 255
	}
	switch line {
	case "zero":
		return 0
	case "ra":
		return 1
	case "sp":
		return 2
	case "gp":
		return 3
	case "tp":
		return 4
	case "fp":
		return 8
	}
	reg, err = strconv.ParseUint(line[1:], 10, 4)
	if err != nil {
		return 255
	}
	switch line[0] {
	case 't':
		if reg < 3 {
			reg += 5
		} else {
			reg += 25
			if reg > 31 {
				return 255
			}
		}
		return uint8(reg)
	case 'f':
		return 255
	case 's':
		if reg < 2 {
			reg += 8
		} else {
			reg += 16
			if reg > 27 {
				return 255
			}
		}
		return uint8(reg)
	case 'a':
		reg += 10
		if reg > 17 {
			return 255
		}
		return uint8(reg)
	}
	return 255
}

func assembleInstr(line string) (uint32, error) {
	tokens := strings.FieldsFunc(line, tokenizer)
	if len(tokens) == 0 {
		return 0, nil
	}
	if len(tokens) > 4 {
		return 0, fmt.Errorf("wrong number of args for instruction %s", tokens[0])
	}
	if len(tokens[0]) < 2 {
		return 0, fmt.Errorf("unknown instruction %s", tokens[0])
	}
	data := new(InstrData)
	if tokens[0] == "lui" {
		data.format = ITYPE_U
		data.opcode = 0x37
		if len(tokens) != 3 {
			return 0, fmt.Errorf("wrong number of args for instruction %s", tokens[0])
		}
		data.rd = getRegister(tokens[1])
		if data.rd == 255 {
			return 0, fmt.Errorf("couldn't understand rd %s", tokens[1])
		}
		imm, err := strconv.ParseInt(tokens[2], 0, 20)
		if err != nil {
			return 0, fmt.Errorf("couldn't understand imm %s", tokens[2])
		}
		data.imm = uint32(imm)
	} else if tokens[0] == "auipc" {
		data.format = ITYPE_U
		data.opcode = 0x17
		if len(tokens) != 2 {
			return 0, fmt.Errorf("wrong number of args for instruction %s", tokens[0])
		}
		imm, err := strconv.ParseInt(tokens[1], 0, 20)
		if err != nil {
			return 0, fmt.Errorf("couldn't understand imm %s", tokens[1])
		}
		data.imm = uint32(imm)
	} else if tokens[0][0] == 'j' {
		if tokens[0][:3] != "jal" || len(tokens[0]) > 4 {
			return 0, fmt.Errorf("unknown instruction %s", tokens[0])
		}
		if len(tokens) < 3 {
			return 0, fmt.Errorf("wrong number of args for instruction %s", tokens[0])
		}
		data.rd = getRegister(tokens[1])
		if data.rd == 255 {
			return 0, fmt.Errorf("couldn't understand rd %s", tokens[1])
		}
		imm, err := strconv.ParseInt(tokens[2], 0, 20)
		if err != nil {
			return 0, fmt.Errorf("couldn't understand imm %s", tokens[2])
		}
		data.imm = uint32(imm)
		if len(tokens[0]) == 3 {
			data.format = ITYPE_UJ
			data.opcode = 0x6f
		} else {
			data.opcode = 0x67
			data.format = ITYPE_I
			if len(tokens) != 4 {
				return 0, fmt.Errorf("wrong number of args for instruction %s", tokens[0])
			}
			data.rs1 = getRegister(tokens[3])
			if data.rs1 == 255 {
				return 0, fmt.Errorf("couldn't understand rs1 %s", tokens[3])
			}
		}
	} else if tokens[0][0] == 's' {
		if len(tokens[0]) != 2 {
			return 0, fmt.Errorf("unknown instruction %s", tokens[0])
		}
		data.opcode = 0x23
		data.format = ITYPE_S
		switch tokens[0][1] {
		case 'b':
		case 'h':
			data.funct3 = 1
		case 'w':
			data.funct3 = 2
		case 'd':
			data.funct3 = 3
		default:
			return 0, fmt.Errorf("unknown instruction %s", tokens[0])
		}
		data.rs1 = getRegister(tokens[3])
		if data.rs1 == 255 {
			return 0, fmt.Errorf("couldn't understand rs1 %s", tokens[3])
		}
		data.rs2 = getRegister(tokens[1])
		if data.rs2 == 255 {
			return 0, fmt.Errorf("couldn't understand rs2 %s", tokens[1])
		}
		imm, err := strconv.ParseInt(tokens[2], 0, 12)
		if err != nil {
			return 0, fmt.Errorf("couldn't understand imm %s", tokens[2])
		}
		data.imm = uint32(imm)
	} else if tokens[0][0] == 'l' {
		data.format = ITYPE_I
		data.opcode = 3
		if len(tokens[0])&0xfe != 2 {
			return 0, fmt.Errorf("unknown instruction %s", tokens[0])
		}
		switch tokens[0][1] {
		case 'b':
			break
		case 'h':
			data.funct3 = 1
		case 'w':
			data.funct3 = 2
		case 'd':
			data.funct3 = 3
		default:
			return 0, fmt.Errorf("unknown instruction %s", tokens[0])
		}
		if len(tokens[0]) == 3 {
			if tokens[0][1] == 'd' {
				return 0, fmt.Errorf("unknown instruction %s", tokens[0])
			}
			data.funct3 |= 4
		}
		data.rs1 = getRegister(tokens[3])
		if data.rs1 == 255 {
			return 0, fmt.Errorf("couldn't understand rs1 %s", tokens[3])
		}
		data.rd = getRegister(tokens[1])
		if data.rd == 255 {
			return 0, fmt.Errorf("couldn't understand rd %s", tokens[1])
		}
		imm, err := strconv.ParseInt(tokens[2], 0, 12)
		if err != nil {
			return 0, fmt.Errorf("couldn't understand imm %s", tokens[2])
		}
		data.imm = uint32(imm)
	} else if tokens[0][0] == 'b' {
		data.format = ITYPE_SB
		data.opcode = 0x63
		switch tokens[0] {
		case "beq":
		case "bne":
			data.funct3 = 1
		case "blt":
			data.funct3 = 4
		case "bge":
			data.funct3 = 5
		case "bgeu":
			data.funct3 = 7
		case "bltu":
			data.funct3 = 6
		default:
			return 0, fmt.Errorf("unknown instruction %s", tokens[0])
		}
		data.rs1 = getRegister(tokens[2])
		if data.rs1 == 255 {
			return 0, fmt.Errorf("couldn't understand rs1 %s", tokens[2])
		}
		data.rs2 = getRegister(tokens[1])
		if data.rs2 == 255 {
			return 0, fmt.Errorf("couldn't understand rs2 %s", tokens[1])
		}
		imm, err := strconv.ParseInt(tokens[3], 0, 12)
		if err != nil {
			return 0, fmt.Errorf("couldn't understand imm %s", tokens[3])
		}
		data.imm = uint32(imm)
	} else {
		orig := tokens[0]
		if tokens[0][len(tokens[0])-1] == 'w' {
			data.opcode = 8
			tokens[0] = tokens[0][:len(tokens[0])-1]
		}
		if tokens[0][len(tokens[0])-1] == 'i' {
			data.format = ITYPE_I
			data.opcode |= 0x13
			tokens[0] = tokens[0][:len(tokens[0])-1]
		} else {
			data.format = ITYPE_R
			data.opcode |= 0x33
		}
		switch tokens[0] {
		case "add":
		case "sub":
			if data.format == ITYPE_I {
				return 0, errors.New("unknown instruction subi")
			}
			data.funct7 = 0x10
		case "sll":
			data.funct3 = 1
		case "slt":
			data.funct3 = 2
		case "sltiu":
			data.format = ITYPE_I
			data.opcode ^= 0x20
			fallthrough
		case "sltu":
			data.funct3 = 3
		case "xor":
			data.funct3 = 4
		case "sra":
			data.funct7 = 0x10
			fallthrough
		case "sri":
			data.funct3 = 5
		case "or":
			data.funct3 = 6
		case "and":
			data.funct3 = 7
		default:
			return 0, fmt.Errorf("unknown instruction %s", orig)
		}
		if data.opcode&8 != 0 && data.funct3 != 0 && data.funct3 != 1 && data.funct3 != 5 {
			return 0, fmt.Errorf("unknown instruction %s", orig)
		}
		if len(tokens) != 4 {
			return 0, fmt.Errorf("wrong number of args for instruction %s", tokens[0])
		}
		data.rs1 = getRegister(tokens[2])
		if data.rs1 == 255 {
			return 0, fmt.Errorf("couldn't understand rs1 %s", tokens[3])
		}
		data.rd = getRegister(tokens[1])
		if data.rd == 255 {
			return 0, fmt.Errorf("couldn't understand rd %s", tokens[1])
		}
		if data.format == ITYPE_R {
			data.rs2 = getRegister(tokens[3])
			if data.rs2 == 255 {
				return 0, fmt.Errorf("couldn't understand rs2 %s", tokens[1])
			}
		} else {
			imm, err := strconv.ParseInt(tokens[3], 0, 12)
			if err != nil {
				return 0, fmt.Errorf("couldn't understand imm %s", tokens[3])
			}
			data.imm = uint32(imm)
		}
	}
	return encodeInstr(data), nil
}
