package main

/*
extern char** get_strarr();
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	max := 2 * 4
	attrib := make(map[string]string)

	attr := C.get_strarr()
	gattr := (*[1 << 30]*C.char)(unsafe.Pointer(attr))[:max:max]
	for i := 0; i < len(gattr); i += 2 {
		fmt.Printf("i: %d\n", i)
		goname := C.GoString((*C.char)(gattr[i]))
		val := C.GoString((*C.char)(gattr[i+1]))
		attrib[goname] = val
	}
	fmt.Printf("%+v\n", attrib)
}
