/*
@author: sk
@date: 2024/5/26
*/
package main

// https://www.jmeiners.com/lc3-vm/  a simple vm

func main() {
	//num := int16(-2233)
	//uNum := uint16(num)
	//fmt.Println(num, uNum, uint16(num), int16(uNum))
	//fmt.Println(2233 >> 1)
	//fmt.Println(strconv.FormatUint(10, 2), strconv.FormatUint(10, 2))
	//var d uint8 = 2
	//fmt.Printf("%08b\n", d)  // 00000010
	//fmt.Printf("%08b\n", ^d) // 11111101
	vm := NewVM()
	vm.LoadFile("2048.obj")
	vm.Run()
}
