/*
@author: sk
@date: 2024/5/26
*/
package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type VM struct {
	Mem [MemSize]uint16 // 内存
	// 寄存器 这里虽然设置 uint16 有符号需求直接强转就行了
	Rs [8]uint16 // 使用下标访问
	// pc 寄存器
	Pc uint16
	// 条件寄存器
	Flag uint16
	// 计数器
	Count uint16
	Run0  bool
}

func NewVM() *VM {
	return &VM{}
}

func (v *VM) Run() {
	v.Pc = PcStart
	v.Run0 = true
	for v.Run0 {
		cmd := v.Read()
		opCode := (cmd >> 12) & 0xF
		switch opCode { // 最开始的4位位指令位，后面是参数
		case OpAdd:
			v.add(cmd)
		case OpBr:
			v.br(cmd)
		case OpLd:
			v.ld(cmd)
		case OpSt:
			v.st(cmd)
		case OpJsr:
			v.jsr(cmd)
		case OpAnd:
			v.and(cmd)
		case OpLdr:
			v.ldr(cmd)
		case OpStr:
			v.str(cmd)
		case OpNot:
			v.not(cmd)
		case OpLdi:
			v.ldi(cmd)
		case OpSti:
			v.sti(cmd)
		case OpJmp:
			v.jmp(cmd)
		case OpLea:
			v.lea(cmd)
		case OpTrap:
			v.trap(cmd)
		case OpRti, OpRes:
			panic(fmt.Sprintf("unsuport opCode %d", opCode))
		default:
			panic(fmt.Sprintf("unknown opCode %d", opCode))
		}
	}
}

func (v *VM) Read() uint16 {
	v.Pc++
	return v.Mem[v.Pc-1]
}

func (v *VM) UpdateFlag(val uint16) {
	if val == 0 {
		v.Flag = FlagZro
	} else if (val >> 15) > 0 {
		v.Flag = FlagNeg
	} else {
		v.Flag = FlagPos
	}
}

func (v *VM) add(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	sr1 := SubVal(cmd, 6, 8)
	if SubBool(cmd, 5, 5) { // 使用立即数
		num := SubVal(cmd, 0, 4)
		v.Rs[dr] = v.Rs[sr1] + num
	} else { // 使用寄存器
		sr2 := SubVal(cmd, 0, 2)
		v.Rs[dr] = v.Rs[sr1] + v.Rs[sr2]
	}
	v.UpdateFlag(v.Rs[dr]) // 更新 flag
}

func (v *VM) ldi(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	pcOffset := SubVal(cmd, 0, 8)
	// 获取的是地址
	addr := v.Mem[v.Pc+pcOffset]
	v.Rs[dr] = v.Mem[addr]
	v.UpdateFlag(v.Rs[dr]) // 更新 flag
}

func (v *VM) and(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	sr1 := SubVal(cmd, 6, 8)
	if SubBool(cmd, 5, 5) { // 使用立即数
		num := SubVal(cmd, 0, 4)
		v.Rs[dr] = v.Rs[sr1] & num
	} else { // 使用寄存器
		sr2 := SubVal(cmd, 0, 2)
		v.Rs[dr] = v.Rs[sr1] & v.Rs[sr2]
	}
	v.UpdateFlag(v.Rs[dr]) // 更新 flag
}

func (v *VM) not(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	sr := SubVal(cmd, 6, 8)
	v.Rs[dr] = ^v.Rs[sr]
	v.UpdateFlag(v.Rs[dr])
}

func (v *VM) br(cmd uint16) {
	pcOffset := SubVal(cmd, 0, 8)
	flag := SubVal(cmd, 9, 11)
	if (v.Flag & flag) > 0 {
		v.Pc += pcOffset
	}
}

func (v *VM) jmp(cmd uint16) {
	sr := SubVal(cmd, 6, 8)
	v.Pc = v.Rs[sr]
}

func (v *VM) jsr(cmd uint16) {
	if SubBool(cmd, 11, 11) {
		pcOffset := SubVal(cmd, 0, 10)
		v.Pc += pcOffset
	} else {
		sr := SubVal(cmd, 6, 8)
		v.Pc = v.Rs[sr]
	}
}

func (v *VM) ld(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	pcOffset := SubVal(cmd, 0, 8)
	v.Rs[dr] = v.Mem[v.Pc+pcOffset]
	v.UpdateFlag(v.Rs[dr])
}

func (v *VM) ldr(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	sr := SubVal(cmd, 6, 8)
	addrOffset := SubVal(cmd, 0, 5)
	v.Rs[dr] = v.Mem[v.Rs[sr]+addrOffset]
	v.UpdateFlag(v.Rs[dr])
}

func (v *VM) lea(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	addrOffset := SubVal(cmd, 0, 8)
	v.Rs[dr] = v.Pc + addrOffset
	v.UpdateFlag(v.Rs[dr])
}

func (v *VM) st(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	addrOffset := SubVal(cmd, 0, 8)
	v.Mem[v.Pc+addrOffset] = v.Rs[dr]
}

func (v *VM) sti(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	pcOffset := SubVal(cmd, 0, 8)
	v.Mem[v.Mem[v.Pc+pcOffset]] = v.Rs[dr]
}

func (v *VM) str(cmd uint16) {
	dr := SubVal(cmd, 9, 11)
	sr := SubVal(cmd, 6, 8)
	addrOffset := SubVal(cmd, 0, 5)
	v.Mem[v.Rs[sr]+addrOffset] = v.Rs[dr]
}

func (v *VM) trap(cmd uint16) {
	trap := SubVal(cmd, 0, 7)
	switch trap {
	case TrapGetC:
		v.getc()
	case TrapOut:
		v.out()
	case TrapPutS:
		v.puts()
	case TrapIn:
		v.in()
	case TrapPutSp:
		v.putsp()
	case TrapHalt:
		v.halt()
	}
}

func (v *VM) puts() {
	fmt.Println("VM.puts")
}

func (v *VM) getc() {
	v.Rs[R0] = 2233
	v.UpdateFlag(v.Rs[R0])
}

func (v *VM) out() {
	fmt.Println("VM.out")
}

func (v *VM) in() {
	v.Rs[R0] = 2233
	v.UpdateFlag(v.Rs[R0])
}

func (v *VM) putsp() {
	fmt.Println("VM.putsp")
}

func (v *VM) halt() {
	v.Run0 = false
	fmt.Println("VM.halt")
}

func (v *VM) LoadFile(file string) {
	bs, err := os.ReadFile(file)
	HandleErr(err)
	offset := binary.LittleEndian.Uint16(bs[:2])
	index := 3
	for int(offset) < len(v.Mem) && index < len(bs) {
		v.Mem[offset] = binary.LittleEndian.Uint16(bs[index-1 : index+1])
	}
}
