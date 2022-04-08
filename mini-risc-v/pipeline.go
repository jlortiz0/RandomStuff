package main

type Pipeline struct {
	fetch     uint32
	decode    *InstrStage2
	execute   *InstrStage3
	mem       *InstrStage4
	writeback *InstrStage5
}

func pushPipeline(pipe *Pipeline, instr uint32, registers *[33]int64) bool {
	pipe.writeback = execStage5(pipe.mem, registers)
	registers[0] = 0
	pipe.mem = execStage4(pipe.execute)
	if pipe.mem.opcode == 3 && pipe.mem.rd != 0 {
		if pipe.mem.rd == pipe.execute.rs1 {
			pipe.execute.rs1Out = pipe.mem.memOut
		}
		if pipe.mem.rd == pipe.execute.rs2 {
			pipe.execute.rs2Out = pipe.mem.memOut
		}
		if pipe.mem.rd == pipe.decode.rs1 {
			pipe.decode.rs1Out = pipe.mem.memOut
		}
		if pipe.mem.rd == pipe.decode.rs2 {
			pipe.decode.rs2Out = pipe.mem.memOut
		}
	}
	pipe.execute = execStage3(pipe.decode)
	if pipe.execute.format == ITYPE_SB {
		yes := pipe.execute.ALUOut == 0
		if pipe.execute.funct3&1 == 1 {
			yes = !yes
		}
		if yes {
			pipe.decode = nil
			pipe.fetch = 0
			registers[REG_PC] += signExtend(pipe.execute.imm, 12)
			pipe.execute = nil
			return false
		}
	}
	if pipe.execute.rd != 0 && pipe.execute.rd != REG_PC && pipe.execute.opcode != 3 {
		if pipe.execute.rd == pipe.decode.rs1 {
			pipe.decode.rs1Out = pipe.mem.memOut
		}
		if pipe.execute.rd == pipe.decode.rs2 {
			pipe.decode.rs2Out = pipe.mem.memOut
		}
	}
	pipe.decode = execStage2(pipe.fetch, registers)
	pipe.fetch = instr
	return true
}
