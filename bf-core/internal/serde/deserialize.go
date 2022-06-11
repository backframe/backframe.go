package serde

import "backframe.io/backframe/bf-core/internal/parser"

func Deserialize(contents []byte) []parser.Block {
	p := parser.NewParser()
	return p.Parse(string(contents))
}
