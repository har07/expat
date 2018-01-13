package main

/*
#include <stdlib.h>

extern char* GetError();
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	cstr := C.GetError()
	defer C.free(unsafe.Pointer(cstr))
	msg := C.GoString(cstr)
	fmt.Println(msg)
}
