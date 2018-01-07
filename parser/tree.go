package parser

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type XMLParser interface {
	ParseWhole(data string) (*Element, error)
	Close() (*Element, error)
	Feed(data string) error
	Free()
}

// ParseError contains error info on parsing tree
type ParseError struct {
	Desc   string
	Code   int
	Line   int
	Column int
}

// Element is an XML element.
// This class is the reference implementation of the Element interface.
// An element's length is its number of subelements.  That means if you
// want to check if an element is truly empty, you should check BOTH
// its length AND its text attribute.
// The element tag, attribute names, and attribute values can be either
// bytes or strings.
// *tag* is the element name.  *attrib* is an optional dictionary containing
// element attributes. *extra* are additional element attributes given as
// keyword arguments.
// Example form:
// 	<tag attrib>text<child/>...</tag>tail
type Element struct {
	Tag      string
	Tail     *string
	Text     *string
	Children []*Element
	Attrib   map[string]string
}

// FromString parses XML document from string constant.
// This function can be used to embed "XML Literals" in Python code.
// *text* is a string containing XML data, *parser* is an
// optional parser instance, defaulting to the standard XMLParser.
// Returns an Element instance.
func FromString(text string, _xmlpar ...XMLParser) (*Element, error) {
	var xmlpar XMLParser
	if len(_xmlpar) > 0 {
		xmlpar = _xmlpar[0]
	} else {
		xmlpar = NewExpatParser("UTF-8", true, NewBuilder())
	}
	defer xmlpar.Free()
	if err := xmlpar.Feed(text); err != nil {
		return nil, fmt.Errorf("Parsing finished with error: %s", err.Error())
	}
	return xmlpar.Close()
}

// CreateElement return new Element
func CreateElement(tag string, attrib map[string]string) Element {
	return Element{
		Tag:    tag,
		Attrib: attrib,
	}
}

// String ...
func (e *Element) String() string {
	return fmt.Sprintf("<Element %s attrib=%+v>", e.Tag, e.Attrib)
}

// Append adds *subelement* to the end of this element.
// The new element will appear in document order after the last existing
// subelement (or directly after the text, if it's the first subelement),
// but before the end tag for this element.
func (e *Element) Append(subelement *Element) {
	e.Children = append(e.Children, subelement)
}

// Extend appends subelements from a sequence.
// *elements* is a sequence with zero or more elements.
func (e *Element) Extend(elements []*Element) {
	e.Children = append(e.Children, elements...)
}

// Insert inserts *subelement* at position *index*.
func (e *Element) Insert(subelement *Element, index int) {
	i := index
	e.Children = append(e.Children, nil)
	copy(e.Children[i+1:], e.Children[i:])
	e.Children[i] = subelement
}

// Remove Remove matching subelement.
// Unlike the find methods, this method compares elements based on
// identity, NOT ON tag value or contents.
func (e *Element) Remove(subelement *Element) {
	for i, v := range e.Children {
		if v == subelement {
			copy(e.Children[i:], e.Children[i+1:])
			e.Children[len(e.Children)-1] = nil
			e.Children = e.Children[:len(e.Children)-1]
		}
	}
}

// Clear resets element.
// This function removes all subelements, clears all attributes, and sets
// the text and tail attributes to nil.
func (e *Element) Clear() {
	e.Attrib = make(map[string]string)
	e.Children = []*Element{}
	e.Text, e.Tail = nil, nil
}

// Get returns element attribute.
// *key* is what attribute to look for, and *defaulted* is what to return
// if the attribute was not found.
// Returns a string containing the attribute value, or the default if
// attribute was not found.
func (e *Element) Get(key string, defaulted ...string) string {
	if val, ok := e.Attrib[key]; ok {
		return val
	}
	if len(defaulted) > 0 {
		return defaulted[0]
	}
	return ""
}

// Set element attribute.
// Equivalent to Attrib[key] = value, but some implementations may handle
// this a bit more efficiently.  *key* is what attribute to set, and
// *value* is the attribute value to set it to.
func (e *Element) Set(key string, value string) {
	e.Attrib[key] = value
}

// Keys returns list of attribute names.
// Names are returned in an arbitrary order, just like an ordinary
// Golang map.  Equivalent to attrib.keys()
func (e *Element) Keys() []string {
	var keys []string
	for k := range e.Attrib {
		keys = append(keys, k)
	}
	return keys
}

// TODO: implement XPath like query methods
// Find()
// FindText()
// FindAll()
// IterFind()

// Question: should I port iterators?
// Iter()
// IterText()

// ElementTree is an XML element hierarchy.
// This class also provides support for serialization to and from
// standard XML.
// *element* is an optional root element node,
// *file* is an optional file handle or file name of an XML file whose
// contents will be used to initialize the tree with.
type ElementTree struct {
	root *Element
}

// CreateTree returns new ElementTree
func CreateTree(element *Element, file ...string) (ElementTree, error) {
	tree := ElementTree{root: element}
	if len(file) > 0 {
		if err := tree.Parse(file[0]); err != nil {
			return tree, err
		}
	}
	return tree, nil
}

// Root returns root element of this tree.
func (t *ElementTree) Root() *Element {
	return t.root
}

// Parse Load external XML document into element tree.
// *source* is a file name or file object, *xmlpar* is an optional parser
// instance that defaults to XMLParser.
// ParseError is raised if the parser fails to parse the document.
// Returns the root element of the given source document.
func (t *ElementTree) Parse(source string, _xmlpar ...XMLParser) error {
	content, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	var xmlpar XMLParser
	if len(_xmlpar) > 0 {
		xmlpar = _xmlpar[0]
	} else {
		xmlpar = NewExpatParser("UTF-8", true, NewBuilder())
	}
	root, err := xmlpar.ParseWhole(string(content))
	if err != nil {
		return err
	}
	t.root = root
	return nil
}

// TreeBuilder is a generic element structure builder.
// This builder converts a sequence of start, data, and end method
// calls to a well-formed element structure.
// You can use this class to build an element structure using a custom XML
// parser, or a parser for some other XML-like format.
// *element_factory* is an optional element factory which is called
// to create new Element instances, as necessary.
type TreeBuilder struct {
	data    []string
	elem    []*Element
	last    *Element
	tail    bool
	factory TreeFactory
}

// TreeFactory is
type TreeFactory func(tag string, attrs map[string]string) *Element

// NewBuilder returns new TreeBuilder instance
func NewBuilder(factory ...TreeFactory) *TreeBuilder {
	b := &TreeBuilder{}
	if len(factory) > 0 {
		b.factory = factory[0]
	} else {
		// set default factory:
		b.factory = func(tag string, attrs map[string]string) *Element {
			return &Element{
				Tag:    tag,
				Attrib: attrs,
			}
		}
	}
	return b
}

// Close flush builder buffers and return toplevel document Element
func (t *TreeBuilder) Close() (*Element, error) {
	if len(t.elem) != 0 {
		return nil, fmt.Errorf("missing end tags")
	}
	if t.last == nil {
		return nil, fmt.Errorf("missing top level element")
	}
	return t.last, nil
}

func (t *TreeBuilder) flush() error {
	if len(t.data) > 0 {
		if t.last != nil {
			text := strings.Join(t.data, "")
			if t.tail {
				if t.last.Tail != nil {
					return fmt.Errorf("internal error (tail)")
				}
				t.last.Tail = &text
			} else {
				if t.last.Text != nil {
					return fmt.Errorf("internal error (text)")
				}
				t.last.Text = &text
			}
		}
		t.data = []string{}
	}
	return nil
}

// Data add text to current element
func (t *TreeBuilder) Data(data string) error {
	t.data = append(t.data, data)
	return nil
}

// Start open new element and return it.
// *tag* is the element name, *attrs* is a dict containing element attributes
func (t *TreeBuilder) Start(tag string, attrs map[string]string) (*Element, error) {
	if err := t.flush(); err != nil {
		return nil, err
	}
	elem := t.factory(tag, attrs)
	t.last = elem
	if len(t.elem) > 0 {
		t.elem[len(t.elem)-1].Append(elem)
	}
	t.elem = append(t.elem, elem)
	t.tail = false
	return elem, nil
}

// End close and return current Element.
// *tag* is the element name.
func (t *TreeBuilder) End(tag string) (*Element, error) {
	t.flush()
	t.last = t.elem[len(t.elem)-1]
	t.elem = t.elem[:len(t.elem)-1]
	if t.last.Tag != tag {
		return nil, fmt.Errorf("end tag mismatch (expected %s, got %s)", t.last.Tag, tag)
	}
	t.tail = true
	return t.last, nil
}
