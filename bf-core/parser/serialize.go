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
		contents.WriteString(WriteBlock(b))
		contents.WriteString("\n")
	}
	return []byte(contents.String())
}

func WriteBlock(b Block) string {
	var tmpl strings.Builder
	var blocksTmpl strings.Builder

	tmpl.WriteString(fmt.Sprintf("%v %v { \n", strings.ToLower(b._type), b.id))

	// handle block variables
	for k, v := range b.variables {
		if strings.ContainsAny(v, "|") {
			// its an array
			val := strings.ReplaceAll(v, "|", ",")
			tmp := fmt.Sprintf("%v = [%v,];", k, val)
			tmpl.WriteString(fmt.Sprintf("\t%v\n", tmp))
		} else {
			tmp := fmt.Sprintf("%v = %v;", k, v)
			tmpl.WriteString(fmt.Sprintf("\t%v\n", tmp))
		}
	}

	// handle nested blocks
	for i := 0; i < len(b.blocks); i++ {
		currentBlock := b.blocks[i]
		tmpl := WriteBlock(currentBlock)
		blocksTmpl.WriteString(tmpl)
	}
	if len(blocksTmpl.String()) > 0 {
		tmpl.WriteString(fmt.Sprintf("%v\n", blocksTmpl.String()))
	}
	tmpl.WriteString("}")

	return tmpl.String()
}
