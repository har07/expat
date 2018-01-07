package main

import (
	"fmt"

	"github.com/har07/expat/debug"
	"github.com/har07/expat/etree"
)

func main() {
	raw := `<root>
	<child><grandchild foo="bar" name="fulan"/></child>
	<child>baz</child>
	<child/>
</root>`
	root, err := etree.FromString(raw)
	if err != nil {
		fmt.Printf("FromString call error: %s\n", err.Error())
		return
	}
	fmt.Println(debug.PrintJSON(root))
}
