package test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type Hello struct {
	KeyA string
	KeyB string
	KeyC *int
}

// 增加String方法, 类似于java里的toString(),
// 在打印或其它地方会返回该方法的结果
func (hello *Hello) String() string {
	// 如果为空，直接返回
	if hello == nil {
		return "nil"
	}
	typ := reflect.TypeOf(hello).Elem()
	obj := reflect.ValueOf(hello).Elem()
	numField := typ.NumField()
	buffer := new(bytes.Buffer)
	for i := 0; i < numField; i++ {
		key := typ.Field(i).Name
		value := obj.Field(i)
		fmt.Fprintf(buffer, "%v:\t%v\n", key, value)
		// fmt.Println(key, value)
	}
	return buffer.String()
}
func TestFmt(t *testing.T) {
	var test *Hello = &Hello{
		KeyA: "ValueA",
		KeyB: "ValueB",
	}
	var integer int = 1
	var intPtr *int = &integer
	fmt.Println(intPtr)
	// 该方法返回自定义字符串
	fmt.Println(test)
	// 该方法返回系统默认字符串, 因为只绑定了struct 指针, 但是这里传递的是struct 值
	fmt.Printf("%+v \n", *test)
}

func TestBuffer(t *testing.T) {
	// buffer.String() 返回未被读取的内容组成的字符串
	var buffer bytes.Buffer
	buffer.WriteString("prefix")
	// 这里打印的是 prefix
	fmt.Println(buffer.String())
	buffer.WriteString("-add content-")
	// 这里打印的是 prefix-add content-
	fmt.Println(buffer.String())
	// 读取buffer内容, 直到第一次出现'-', 返回的内容包括'-'
	line, err := buffer.ReadString('-')
	// 这里打印的是 prefix-
	fmt.Println(line, err)
	buffer.WriteString("suffix")
	// 这里打印的是 add content-suffix
	fmt.Println(buffer.String())

}
