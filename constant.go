/*
@author: sk
@date: 2024/5/26
*/
package main

const (
	MemSize = 1 << 16
)

const (
	OpBr   uint16 = iota // 符合要求偏移 pc
	OpAdd                // add
	OpLd                 // load
	OpSt                 // store
	OpJsr                // 跳转指针
	OpAnd                // and
	OpLdr                // load 寄存器
	OpStr                // store 寄存器
	OpRti                // 没有使用 返回值？
	OpNot                // not
	OpLdi                // load 内存中的数到寄存器
	OpSti                // store 直接数
	OpJmp                // jump
	OpRes                // 没有使用 反转？
	OpLea                // 加载有效的地址
	OpTrap               // 调用系统例程，陷入内核态
)

const (
	FlagPos uint16 = 1 << iota
	FlagZro
	FlagNeg
)

const (
	PcStart = 0x3000
)

const (
	R0 = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
)

const (
	TrapGetC  = 0x20
	TrapOut   = 0x21
	TrapPutS  = 0x22
	TrapIn    = 0x23
	TrapPutSp = 0x24
	TrapHalt  = 0x25
)
