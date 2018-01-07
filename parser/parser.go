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
extern int GetCurrentColumnNumber(int id);
extern char* GetError(int id, int code);
extern int GetCurrentAttributeCount(int id);
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

type ParseError struct {
	Desc   string
	Code   int
	Line   int
	Column int
}

type XMLParser struct {
	id      int
	version string
	target  *TreeBuilder
	names   map[string]string
	entity  map[string]string
}

var pool *XMLParser

// CreateParser initialize new XMLParser instance
func CreateParser(encoding string, namespace bool, target *TreeBuilder) *XMLParser {
	p := XMLParser{
		names:  make(map[string]string),
		entity: make(map[string]string),
	}
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

	if target != nil {
		p.target = target
	} else {
		t := CreateBuilder()
		p.target = t
	}

	p.version = "Expat version 2.2.5" // TODO: determine actual Expat version

	// TODO: register main callbacks: start_element, end_element, character_data
	// start, end, handler, data
	// TODO: register miscellaneous callbacks: comment, pi

	//register to pool
	pool = &p

	return pool
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("Error [%d] at line %d column %d: %s",
		pe.Code, pe.Line, pe.Column, pe.Desc,
	)
}

// fixName expands qname
func (xp *XMLParser) fixName(key string) (name string) {
	if val, ok := xp.names[key]; ok {
		name = val
	} else {
		name = key
		if strings.Contains(name, "}") {
			name = "{" + name
		}
		xp.names[key] = name
	}
	return name
}

// handler is a default handler for expat events
func (xp *XMLParser) handler(text string) {
	prefix := text[:1]
	if prefix == "&" {
		entityRef := ""
		if val, ok := xp.entity[text[1:len(text)-1]]; ok {
			entityRef = val
		} else {
			// TODO: notify caller about the error
			errMsg := "undefined entity " + text
			cline := C.GetCurrentLineNumber(C.int(xp.id))
			ccol := C.GetCurrentColumnNumber(C.int(xp.id))
			fmt.Printf(ParseError{
				Desc:   errMsg,
				Code:   11, // XML_ERROR_UNDEFINED_ENTITY
				Line:   int(cline),
				Column: int(ccol),
			}.Error())
		}
		xp.target.Data(entityRef)
	}
}

// start is a handler for expat's StartElementHandler. Since ordered_attributes
// is set, the attributes are reported as a list of alternating
// attribute name,value.
func (xp *XMLParser) start(tag string, attrib map[string]string) {
	tag = xp.fixName(tag)
	xp.target.Start(tag, attrib)
}

// end is a handler for expat's EndElementHandler
func (xp *XMLParser) end(tag string) {
	tag = xp.fixName(tag)
	xp.target.End(tag)
}

// data is a handler for expat's CharacterDataHandler
func (xp *XMLParser) data(text string) {
	xp.target.Data(text)
}

// Feed feeds chunk of XML data to be parsed
func (xp *XMLParser) Feed(data string) error {
	cdata := (*C.XML_Char)(C.CString(data))
	defer C.free(unsafe.Pointer(cdata))
	cerr := C.Feed(C.int(xp.id), cdata, C.int(len(data)), C.int(0))
	errCode := int(cerr)
	if errCode != 0 {
		cerrMsg := C.GetError(C.int(xp.id), cerr)
		defer C.free(unsafe.Pointer(cerrMsg))
		cline := C.GetCurrentLineNumber(C.int(xp.id))
		ccol := C.GetCurrentColumnNumber(C.int(xp.id))
		return ParseError{
			Desc:   C.GoString(cerrMsg),
			Code:   errCode,
			Line:   int(cline),
			Column: int(ccol),
		}
	}
	return nil
}

// Close finishes feeding data to parser and return element structure
func (xp *XMLParser) Close() (*Element, error) {
	cdata := (*C.XML_Char)(C.CString(""))
	cerr := C.Feed(C.int(xp.id), cdata, C.int(1), C.int(1))
	errCode := int(cerr)
	if errCode != 0 {
		cerrMsg := C.GetError(C.int(xp.id), cerr)
		defer C.free(unsafe.Pointer(cerrMsg))
		cline := C.GetCurrentLineNumber(C.int(xp.id))
		ccol := C.GetCurrentColumnNumber(C.int(xp.id))
		return nil, ParseError{
			Desc:   C.GoString(cerrMsg),
			Code:   errCode,
			Line:   int(cline),
			Column: int(ccol),
		}
	}
	return xp.target.Close()
}

// ParseWhole parses entire XML document and return the root, if success
func (xp *XMLParser) ParseWhole(data string) (*Element, error) {
	cdata := (*C.XML_Char)(C.CString(data))
	defer C.free(unsafe.Pointer(cdata))
	cerr := C.Feed(C.int(xp.id), cdata, C.int(len(data)), C.int(1))
	errCode := int(cerr)
	if errCode != 0 {
		cerrMsg := C.GetError(C.int(xp.id), cerr)
		defer C.free(unsafe.Pointer(cerrMsg))
		cline := C.GetCurrentLineNumber(C.int(xp.id))
		ccol := C.GetCurrentColumnNumber(C.int(xp.id))
		return nil, ParseError{
			Desc:   C.GoString(cerrMsg),
			Code:   errCode,
			Line:   int(cline),
			Column: int(ccol),
		}
	}
	return xp.target.Close()
}

func (xp *XMLParser) Free() {
	C.Free(C.int(xp.id))
	pool = nil // reset id/remove from pool
}

// GStartElementHandler ....
//export GStartElementHandler
func GStartElementHandler(id C.int, el *C.XML_Char, attr **C.XML_Char) {
	defer C.free(unsafe.Pointer(el))

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
	for i := 0; i < len(gattr); i += 2 {
		goname := C.GoString((*C.char)(gattr[i]))
		val := C.GoString((*C.char)(gattr[i+1]))
		attrib[goname] = val
	}
	p.start(tag, attrib)
	return
}

// GEndElementHandler is
//export GEndElementHandler
func GEndElementHandler(id C.int, el *C.XML_Char) {
	defer C.free(unsafe.Pointer(el))

	// get parser by id
	p := pool

	// invoke corresponding handler
	tag := C.GoString((*C.char)(el))
	p.end(tag)
	return
}

// GDefaultHandler is
//export GDefaultHandler
func GDefaultHandler(id C.int, s *C.XML_Char, length C.int) {
	defer C.free(unsafe.Pointer(s))

	// get parser by id
	p := pool

	text := C.GoStringN((*C.char)(s), length)
	p.handler(text)
}

// GCharDataHandler is
//export GCharDataHandler
func GCharDataHandler(id C.int, s *C.XML_Char, length C.int) {
	defer C.free(unsafe.Pointer(s))

	// get parser by id
	p := pool

	text := C.GoStringN((*C.char)(s), length)
	p.data(text)
}
