package parser

import (
	"fmt"
	"strings"
)

// TODO: probably use QName type for Tag, and Attribs key
type Element struct {
	Tag      string
	Tail     *string
	Text     *string
	Children []*Element
	Attribs  map[string]string
}

type TreeBuilder struct {
	data    []string
	elem    []*Element
	last    *Element
	tail    bool
	factory TreeFactory
}

type TreeFactory func(tag string, attrs map[string]string) *Element

// Append add Element as child of current Element
func (e *Element) Append(child *Element) {
	e.Children = append(e.Children, child)
}

// Create returns new TreeBuilder instance
func Create(factory ...TreeFactory) TreeBuilder {
	b := TreeBuilder{}
	if len(factory) > 0 {
		b.factory = factory[0]
	} else {
		// set default factory:
		b.factory = func(tag string, attrs map[string]string) *Element {
			return &Element{
				Tag:     tag,
				Attribs: attrs,
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
func (t *TreeBuilder) Data(data string) {
	t.data = append(t.data, data)
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
	t.tail = false
	return t.last, nil
}
