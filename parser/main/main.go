package main

import (
	"fmt"

	"github.com/har07/expat/parser"
)

func main() {
	raw := `<root>
	<child><grandchild foo="bar" name="fulan"/></child>
	<chil%d></child>
	<child/>
</root>`
	p := parser.Create("UTF-8", true)
	start := func(tag string, attrib map[string]string) {
		fmt.Printf("start %s\n", tag)
		if len(attrib) > 0 {
			fmt.Printf("  attrib: %+v\n", attrib)
		}
	}
	end := func(tag string) {
		fmt.Printf("end %s\n", tag)
	}
	p.SetElementHandler(start, end)
	if err := p.Parse(raw); err != nil {
		fmt.Println(err.Error())
	}
	p.Free()
}
