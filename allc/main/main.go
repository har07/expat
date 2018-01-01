package main

/*
#cgo LDFLAGS: -lexpat
#include "../demo.c"
*/
import "C"
import "fmt"
import "unsafe"

func main() {
	raw := `<root>
	<child><grandchild/></child>
	<child></child>
	<child/>
</root>`
	cs := C.CString(raw)
	defer C.free(unsafe.Pointer(cs))
	xc := C.demo(cs, C.int(len(raw)))
	fmt.Printf("exit code: %d\n", int(xc))
}
