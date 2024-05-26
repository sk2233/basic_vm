/*
@author: sk
@date: 2024/5/26
*/
package main

func SubVal(val uint16, start, end int) uint16 {
	val <<= 15 - end
	return val >> (15 - (end - start))
}

func SubBool(val uint16, start, end int) bool {
	return SubVal(val, start, end) > 0
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}
