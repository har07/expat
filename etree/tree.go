package etree

import (
	"fmt"
	"io/ioutil"

	"github.com/har07/expat/etree/builder"
	"github.com/har07/expat/etree/element"
	"github.com/har07/expat/etree/parser"
)

// ElementTree is an XML element hierarchy.
// This class also provides support for serialization to and from
// standard XML.
// *element* is an optional root element node,
// *file* is an optional file handle or file name of an XML file whose
// contents will be used to initialize the tree with.
type ElementTree struct {
	root *element.E
}

// Root returns root element of this tree.
func (t *ElementTree) Root() *element.E {
	return t.root
}

// Parse Load external XML document into element tree.
// *source* is a file name or file object, *xmlpar* is an optional parser
// instance that defaults to parser.XML.
// ParseError is raised if the parser fails to parse the document.
// Returns the root element of the given source document.
func (t *ElementTree) Parse(source string, _xmlpar ...parser.XML) error {
	content, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	var xmlpar parser.XML
	if len(_xmlpar) > 0 {
		xmlpar = _xmlpar[0]
	} else {
		xmlpar = parser.NewExpat("UTF-8", true, builder.New())
	}
	root, err := xmlpar.ParseWhole(string(content))
	if err != nil {
		return err
	}
	t.root = root
	return nil
}

// FromString parses XML document from string constant.
// This function can be used to embed "XML Literals" in Python code.
// *text* is a string containing XML data, *parser* is an
// optional parser instance, defaulting to the standard parser.XML.
// Returns an element.E instance.
func FromString(text string, _xmlpar ...parser.XML) (*element.E, error) {
	var xmlpar parser.XML
	if len(_xmlpar) > 0 {
		xmlpar = _xmlpar[0]
	} else {
		xmlpar = parser.NewExpat("UTF-8", true, builder.New())
	}
	defer xmlpar.Free()
	if err := xmlpar.Feed(text); err != nil {
		return nil, fmt.Errorf("Parsing finished with error: %s", err.Error())
	}
	return xmlpar.Close()
}

// CreateElement return new element.E
func CreateElement(tag string, attrib map[string]string) element.E {
	return element.E{
		Tag:    tag,
		Attrib: attrib,
	}
}

// CreateTree returns new ElementTree
func CreateTree(element *element.E, file ...string) (ElementTree, error) {
	tree := ElementTree{root: element}
	if len(file) > 0 {
		if err := tree.Parse(file[0]); err != nil {
			return tree, err
		}
	}
	return tree, nil
}
