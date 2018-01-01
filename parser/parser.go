package parser

/*
#cgo LDFLAGS: -lexpat
#include <stdlib.h>
#include <expat.h>
#include "parser.h"

extern int Create(XML_Char *encoding, int namespace);
extern int Feed(int id, XML_Char *chunk, int len, int finish);
extern void SetHandlers(int id, int start, int end);
extern void Free(int id);
extern int GetCurrentLineNumber(int id);
extern char* GetError(int id, int code);
extern int GetCurrentAttributeCount(int id);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type StartElementHandler func(tag string, attrib map[string]string)
type EndElementHandler func(tag string)

type XMLParser struct {
	_parser
}

type _parser struct {
	id    int
	start StartElementHandler
	end   EndElementHandler
}

var pool *XMLParser

func Create(encoding string, namespace bool) *XMLParser {
	p := XMLParser{}
	var cnamespace C.int
	if namespace {
		cnamespace = C.int(1)
	} else {
		cnamespace = C.int(0)
	}
	cencoding := C.CString(encoding)
	defer C.free(unsafe.Pointer(cencoding))
	id := C.Create((*C.XML_Char)(cencoding), cnamespace)
	p.id = int(id)

	//register to pool
	pool = &p

	return pool
}

func (xp *XMLParser) Parse(data string) error {
	cdata := (*C.XML_Char)(C.CString(data))
	defer C.free(unsafe.Pointer(cdata))
	// C.Feed(C.int(xp.id), cdata, C.int(len(data)), C.int(1))
	cerr := C.Feed(C.int(xp.id), cdata, C.int(len(data)), C.int(1))
	errCode := int(cerr)
	if errCode != 0 {
		cerrMsg := C.GetError(C.int(xp.id), cerr)
		defer C.free(unsafe.Pointer(cerrMsg))
		return fmt.Errorf("parsing error with code %d: %s", errCode, C.GoString(cerrMsg))
	}
	return nil
}

func (xp *XMLParser) Free() {
	C.Free(C.int(xp.id))
	pool = nil // reset id/remove from pool
}

func (xp *XMLParser) SetElementHandler(s StartElementHandler, e EndElementHandler) {
	xp.start = s
	xp.end = e
}

//export GStartElementHandler
func GStartElementHandler(id C.int, el *C.XML_Char, attr **C.XML_Char) {
	// get parser by id
	p := pool
	tag := C.GoString((*C.char)(el))
	max := int(C.GetCurrentAttributeCount(id))
	if max > 0 && attr == nil {
		fmt.Printf("attr null: %t (count=%d)\n", attr == nil, max)
	}
	if max == 0 || attr == nil {
		p.start(tag, nil)
		return
	}

	// collect attribute data
	attrib := make(map[string]string)
	gattr := (*[1 << 30]*C.XML_Char)(unsafe.Pointer(attr))[:max:max]
	// _ = (*[1 << 30]*C.XML_Char)(unsafe.Pointer(attr))[:max-1 : max-1]
	for i := 0; i < len(gattr); i += 2 {
		goname := C.GoString((*C.char)(gattr[i]))
		val := C.GoString((*C.char)(gattr[i+1]))
		attrib[goname] = val
	}
	p.start(tag, attrib)
}

//export GEndElementHandler
func GEndElementHandler(id C.int, el *C.XML_Char) {
	// get parser by id
	p := pool
	// invoke corresponding handler
	tag := C.GoString((*C.char)(el))
	p.end(tag)
}
