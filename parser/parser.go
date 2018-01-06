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
	"unsafe"
	"strings"
)

type ParseError struct {
	Desc   string
	Code   int
	Line   int
	Column int
}

// StartElementHandler handler function for start element event
type StartElementHandler func(tag string, attrib map[string]string)

// EndElementHandler handler function for end element event
type EndElementHandler func(tag string)

type XMLParser struct {
	id    int
	version string
	start StartElementHandler
	end   EndElementHandler
	target *TreeBuilder
	names map[string]string
	entity map[string]string
}

var pool *XMLParser

func CreateParser(encoding string, namespace bool, target *TreeBuilder) *XMLParser {
	p := XMLParser{
		names: make(map[string]string),
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

// TODO: return root element
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

func (xp *XMLParser) SetElementHandler(s StartElementHandler, e EndElementHandler) {
	xp.start = s
	xp.end = e
}

// GStartElementHandler is a handler for expat's StartElementHandler. Since ordered_attributes
// is set, the attributes are reported as a list of alternating
// attribute name,value.
//export GStartElementHandler
func GStartElementHandler(id C.int, el *C.XML_Char, attr **C.XML_Char) {
	// get parser by id
	p := pool

	tag := C.GoString((*C.char)(el))
	tag = p.fixName(tag)

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
	p.target.Start(tag, attrib)
	return
}

// GEndElementHandler is 
//export GEndElementHandler
func GEndElementHandler(id C.int, el *C.XML_Char) {
	// get parser by id
	p := pool

	// invoke corresponding handler
	tag := C.GoString((*C.char)(el))
	tag = p.fixName(tag)
	p.target.End(tag)
	return
}

// GDefaultHandler is 
//export GDefaultHandler
func GDefaultHandler(id C.int, s *C.XML_Char, length C.int){
	// get parser by id
	p := pool

	text := C.GoString((*C.char)(s))
	prefix := text[:1]
	if prefix == "&" {
		entityRef := ""
		if val, ok := p.entity[text[1:len(text)-1]]; ok {
			entityRef = val
		} else {
			// TODO: notify caller about the error
			errMsg := "undefined entity " + text
			cline := C.GetCurrentLineNumber(id)
			ccol := C.GetCurrentColumnNumber(id)
			fmt.Printf(ParseError{
				Desc:   errMsg,
				Code:   11, // XML_ERROR_UNDEFINED_ENTITY
				Line:   int(cline),
				Column: int(ccol),
			}.Error())
		}
		p.target.Data(text)
	}
}