package main

import (
	"fmt"

	"github.com/har07/expat/debug"
	"github.com/har07/expat/parser"
)

func main() {
	raw := `<root>
	<child><grandchild foo="bar" name="fulan"/></child>
	<child></child>
	<child/>
</root>`
	root, err := parser.FromString(raw)
	if err != nil {
		fmt.Printf("FromString call error: %s\n", err.Error())
		return
	}
	fmt.Println(debug.PrintJSON(root))
}
