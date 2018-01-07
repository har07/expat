package element

import "fmt"

// E is an XML element.
// This class is the reference implementation of the E interface.
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
type E struct {
	Tag      string
	Tail     *string
	Text     *string
	Children []*E
	Attrib   map[string]string
}

// String ...
func (e *E) String() string {
	return fmt.Sprintf("<Element %s attrib=%+v>", e.Tag, e.Attrib)
}

// Append adds *subelement* to the end of this element.
// The new element will appear in document order after the last existing
// subelement (or directly after the text, if it's the first subelement),
// but before the end tag for this element.
func (e *E) Append(subelement *E) {
	e.Children = append(e.Children, subelement)
}

// Extend appends subelements from a sequence.
// *elements* is a sequence with zero or more elements.
func (e *E) Extend(elements []*E) {
	e.Children = append(e.Children, elements...)
}

// Insert inserts *subelement* at position *index*.
func (e *E) Insert(subelement *E, index int) {
	i := index
	e.Children = append(e.Children, nil)
	copy(e.Children[i+1:], e.Children[i:])
	e.Children[i] = subelement
}

// Remove Remove matching subelement.
// Unlike the find methods, this method compares elements based on
// identity, NOT ON tag value or contents.
func (e *E) Remove(subelement *E) {
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
func (e *E) Clear() {
	e.Attrib = make(map[string]string)
	e.Children = []*E{}
	e.Text, e.Tail = nil, nil
}

// Get returns element attribute.
// *key* is what attribute to look for, and *defaulted* is what to return
// if the attribute was not found.
// Returns a string containing the attribute value, or the default if
// attribute was not found.
func (e *E) Get(key string, defaulted ...string) string {
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
func (e *E) Set(key string, value string) {
	e.Attrib[key] = value
}

// Keys returns list of attribute names.
// Names are returned in an arbitrary order, just like an ordinary
// Golang map.  Equivalent to attrib.keys()
func (e *E) Keys() []string {
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
