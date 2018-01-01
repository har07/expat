package main

import (
	"github.com/har07/expat/allc2"
)

func main() {
	raw := `<root>
	<child><grandchild/></child>
	<child></child>
	<child/>
</root>`
	allc2.Parse(raw)
}
