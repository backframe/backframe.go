package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	_string    string
	_lexer     Lexer
	_lookahead Token
}

func NewParser() Parser {
	return Parser{
		_string: "",
		_lexer:  Lexer{},
	}
}

type Block struct {
	_type     string
	id        string
	variables map[string]string
	blocks    []Block
}

/*
Specfile
	: BlockStatementList
	;

BlockStatementList
	: BlockStatement
	| BlockStatementList BlockStatement -> BlockStatement BlockStatement BlockStatement
	;

BlockStatement
	: Block
	| ExpressionStatement
	| BlockStatementList
	;

Block
	: 'TYPE' IDENTIFIER '{' OptBlockStatementList '}'
	;

ExpressionStatement
	: IDENTIFIER ASSIGNMENT_OP LITERAL
	;

Literal
	: STRING
	| ARRAY
	| OBJECT
	;

*/
func (p *Parser) Parse(val string) []Block {
	p._string = val
	p._lexer.init(val)
	p._lookahead, _ = p._lexer.getNextToken()

	return p.blockStatement()
}

func (p *Parser) statementList() (map[string]string, []Block) {
	variables, blocks := p.statement()

	for {
		if p._lookahead._type == "}" {
			break
		}
		newVars, newBlocks := p.statement()
		blocks = append(blocks, newBlocks...)
		for k, v := range newVars {
			variables[k] = v
		}
	}

	return variables, blocks
}

func (p *Parser) statement() (map[string]string, []Block) {
	switch p._lookahead._type {
	case "IDENTIFIER":
		return p.expressionStatement(), []Block{}
	default:
		return make(map[string]string), p.blockStatement()
	}
}

func (p *Parser) expressionStatement() map[string]string {
	values := make(map[string]string)

	for {
		if p._lookahead._type == "}" || p._lookahead._type != "IDENTIFIER" {
			break
		}

		for {
			if p._lookahead._type == "}" || p._lookahead._type == ";" || p._lookahead._type != "IDENTIFIER" {
				break
			}
			v := p._eat("IDENTIFIER")
			p._eat("ASSIGNMENT_OP")
			// TODO: Implement literals
			id := p.literal()

			values[v.value] = id.value
			p._eat(";")
		}
	}

	return values
}

func (p *Parser) literal() Token {
	switch p._lookahead._type {
	case "STRING":
		return p._eat("STRING")
	case "BOOLEAN":
		return p._eat("BOOLEAN")
	case "[":
		return p.array()
	default:
		return p._eat("STRING")
	}
}

func (p *Parser) array() Token {
	list := []string{}

	p._eat("[")

	for {
		if p._lookahead.value == "]" {
			break
		}
		val := p._eat("STRING")
		p._eat("COMMA")
		list = append(list, val.value)
	}

	p._eat("]")
	return Token{
		_type: "ARRAY",
		value: strings.Join(list, "|"),
	}
}

func (p *Parser) blockStatement() []Block {
	list := []Block{}

	for {
		if p._lookahead._type == "}" || !p._lexer.hasMoreTokens() {
			break
		}
		v := p._eat(p._lookahead._type)
		id := p._eat("IDENTIFIER")
		p._eat("{")

		b := Block{
			_type: v._type,
			id:    id.value,
		}

		vars, blocks := p.statementList()
		b.variables = vars
		b.blocks = blocks

		list = append(list, b)

		p._eat("}")
	}

	return list
}

func (p *Parser) _eat(_type string) Token {
	token := p._lookahead
	if _type != token._type {
		panic(fmt.Sprintf("Syntax Error, Unexpected Token. Expected: %v but instead found: %v\n", _type, token._type))
	}

	p._lookahead, _ = p._lexer.getNextToken()
	return token
}
