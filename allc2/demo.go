package allc2

/*
#cgo LDFLAGS: -lexpat
#include <stdlib.h>
extern int parse(char *data, int len);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func Parse(data string) {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))
	xc := C.parse(cs, C.int(len(data)))
	fmt.Printf("exit code: %d\n", int(xc))
}
