package main

import (
	"fmt"

	"github.com/har07/expat/parser"
	"github.com/har07/expat/debug"
)

func main() {
	raw := `<root>
	<child><grandchild foo="bar" name="fulan"/></child>
	<chil%d></child>
	<child/>
</root>`
	root := parser.FromString(raw)
	fmt.Println(debug.PrintJSON(root))
}
