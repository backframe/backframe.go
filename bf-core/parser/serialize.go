package parser

import (
	"fmt"
	"strings"
)

func Serialize(blocks []Block) []byte {
	var contents strings.Builder
	// Start with each section
	for _, b := range blocks {
		// b is a section
		contents.WriteString(WriteBlock(b, 0))
		contents.WriteString("\n")
	}
	return []byte(contents.String())
}

func WriteBlock(b Block, depth int) string {
	var tmpl strings.Builder
	var blocksTmpl strings.Builder
	var ident strings.Builder

	hasContent := len(b.variables) > 0 || len(b.blocks) > 0

	for i := 0; i < depth; i++ {
		ident.WriteString("\t")
	}

	tmpl.WriteString(fmt.Sprintf("%v%v %v {", ident.String(), strings.ToLower(b._type), b.id))

	if hasContent {
		tmpl.WriteString("\n")
	}

	// handle block variables
	for k, v := range b.variables {
		if strings.ContainsAny(v, "|") {
			// its an array
			val := strings.ReplaceAll(v, "|", ",")
			tmp := fmt.Sprintf("%v = [%v,];", k, val)
			tmpl.WriteString(fmt.Sprintf("\t%v%v\n", ident.String(), tmp))
		} else {
			tmp := fmt.Sprintf("%v = %v;", k, v)
			tmpl.WriteString(fmt.Sprintf("\t%v%v\n", ident.String(), tmp))
		}
	}

	// handle nested blocks
	for i := 0; i < len(b.blocks); i++ {
		currentBlock := b.blocks[i]
		tmpl := WriteBlock(currentBlock, depth+1)
		blocksTmpl.WriteString(fmt.Sprintf("\n%v\n", tmpl))
	}

	if len(blocksTmpl.String()) > 0 {
		tmpl.WriteString(fmt.Sprintf("%v", blocksTmpl.String()))
	}

	if hasContent {
		tmpl.WriteString(ident.String())
	}
	tmpl.WriteString("}")

	return tmpl.String()
}
