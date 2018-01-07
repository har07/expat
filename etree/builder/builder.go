package builder

import (
	"fmt"
	"strings"

	"github.com/har07/expat/etree/element"
)

// Tree is a generic element structure builder.
// This builder converts a sequence of start, data, and end method
// calls to a well-formed element structure.
// You can use this class to build an element structure using a custom XML
// parser, or a parser for some other XML-like format.
// *element_factory* is an optional element factory which is called
// to create new element.E instances, as necessary.
type Tree struct {
	data    []string
	elem    []*element.E
	last    *element.E
	tail    bool
	factory TreeFactory
}

// TreeFactory is
type TreeFactory func(tag string, attrs map[string]string) *element.E

// New returns new Tree instance
func New(factory ...TreeFactory) *Tree {
	b := &Tree{}
	if len(factory) > 0 {
		b.factory = factory[0]
	} else {
		// set default factory:
		b.factory = func(tag string, attrs map[string]string) *element.E {
			return &element.E{
				Tag:    tag,
				Attrib: attrs,
			}
		}
	}
	return b
}

// Close flush builder buffers and return toplevel document element.E
func (t *Tree) Close() (*element.E, error) {
	if len(t.elem) != 0 {
		return nil, fmt.Errorf("missing end tags")
	}
	if t.last == nil {
		return nil, fmt.Errorf("missing top level element")
	}
	return t.last, nil
}

func (t *Tree) flush() error {
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
func (t *Tree) Data(data string) error {
	t.data = append(t.data, data)
	return nil
}

// Start open new element and return it.
// *tag* is the element name, *attrs* is a dict containing element attributes
func (t *Tree) Start(tag string, attrs map[string]string) (*element.E, error) {
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

// End close and return current element.E.
// *tag* is the element name.
func (t *Tree) End(tag string) (*element.E, error) {
	t.flush()
	t.last = t.elem[len(t.elem)-1]
	t.elem = t.elem[:len(t.elem)-1]
	if t.last.Tag != tag {
		return nil, fmt.Errorf("end tag mismatch (expected %s, got %s)", t.last.Tag, tag)
	}
	t.tail = true
	return t.last, nil
}
