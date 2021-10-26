package test

import (
	"fmt"
	"testing"
)

func TestSliceMemory(t *testing.T) {

	// list := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	list := []byte("0123456789")
	slice1 := list[:]
	slice2 := slice1[2:9]
	slice3 := slice2[3:5]

	slice3[1] = byte(255)
	fmt.Println(list)
	fmt.Println(slice1)
	fmt.Println(slice2)
	fmt.Println(slice3)
	fmt.Println("slice切片并不会分配新的内存来复制原来的值")
	/*
		[48 49 50 51 52 53 255 55 56 57]
		[48 49 50 51 52 53 255 55 56 57]
		[50 51 52 53 255 55 56]
		[53 255]
	*/
}

func TestAppendMemory(t *testing.T) {

	ints0 := make([]int, 1, 10)
	ints1 := append(ints0, 1, 2, 3, 4)
	ints2 := append(ints1, 5, 6, 7, 8)
	ints3 := append(ints2, 9, 10, 11, 12)

	ints0[0] = 666
	fmt.Println(ints0)
	fmt.Println(ints1)
	fmt.Println(ints2)
	fmt.Println(ints3)
	/*
		当append结果不超过容量时，不需要额外分配内存
		将输出以下结果：
			[666]
			[666 1 2 3 4]
			[666 1 2 3 4 5 6 7 8]
			[0 1 2 3 4 5 6 7 8 9 10 11 12]
	*/
}
