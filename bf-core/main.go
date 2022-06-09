package main

import (
	"fmt"

	"backframe.io/backframe/bf-core/parser"
)

func main() {
	p := parser.Parser{}
	ast := p.Parse(`section`)

	fmt.Printf("%v", ast)
}
