package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	_string     string
	_lexer      Lexer
	_lookahead  Token
	knownBlocks []string
}

func NewParser() Parser {
	return Parser{
		_string: "",
		_lexer:  Lexer{},
		knownBlocks: []string{
			"INTERFACE",
			"PROVIDER",
			"INTEGRATION",
			"METHOD",
			"VERSION",
			"MIGRATION",
			"DATABASE",
			"RESOURCE",
		},
	}
}

type Block struct {
	Type  string
	Id    string
	Body  []Block
	Value string
}

type SpecAst struct {
	Type     string
	Sections []Block
}

/*
Specfile
	: SectionList
	;

SectionList
	: StatementList
	;

StatementList
	: Statement
	| StatementList Statement
	;

Statement
	: BlockStatement
	| ExpressionStatement
	;

BlockStatement
	: 'TYPE' IDENTIFIER '{' OptStatementList '}'
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
func (p *Parser) Parse(val string) SpecAst {
	p._string = val
	p._lexer.init(val)
	p._lookahead, _ = p._lexer.getNextToken()

	return SpecAst{
		Type:     "Specfile",
		Sections: p.sectionList(),
	}
}

func (p *Parser) sectionList() []Block {
	blocks := []Block{p.section()}

	for {
		if !p._lexer.hasMoreTokens() {
			break
		}

		blocks = append(blocks, p.section())
	}

	return blocks
}

func (p *Parser) section() Block {
	p._eat("SECTION")
	id := p._eat("IDENTIFIER")
	p._eat("{")
	body := p.statementList()
	p._eat("}")

	return Block{
		Type: "SECTION",
		Id:   id.value,
		Body: body,
	}

}

func (p *Parser) statementList() []Block {
	stmnts := []Block{}

	for {
		if p._lookahead._type == "}" {
			break
		}
		stmnts = append(stmnts, p.statement())
	}

	return stmnts
}

func (p *Parser) statement() Block {
	value := p._lexer._string
	idx := p._lexer._cursor + 1

	if string(value[idx]) == "=" {
		return p.expressionStatement()
	} else {
		return p.blockStatement()
	}
}

func (p *Parser) expressionStatement() Block {
	id := p._eat("IDENTIFIER")
	p._eat("ASSIGNMENT_OP")
	token := p.literal()

	p._eat(";")
	return Block{
		Type:  "ASSIGNMENT",
		Id:    id.value,
		Value: token.value,
	}
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

func (p *Parser) blockStatement() Block {
	name := p.blockType()
	id := p._eat("IDENTIFIER")
	p._eat("{")

	body := p.statementList()
	p._eat("}")
	return Block{
		Type: name._type,
		Id:   id.value,
		Body: body,
	}
}

func (p *Parser) blockType() Token {
	nxt := p._lookahead._type

	for _, b := range p.knownBlocks {
		if b == nxt {
			return p._eat(b)
		}
	}

	panic(fmt.Sprintf("SyntaxError on line: %d, column: %d. Unknown block: `%s` found", p.getLine(), p.getCol(), nxt))
}

func (p *Parser) _eat(_type string) Token {
	token := p._lookahead
	if _type != token._type {
		panic(fmt.Sprintf("SyntaxError on line: %d, column: %d. Unexpected Token. Expected: %v but instead found: %v\n", p.getLine(), p.getCol(), _type, token._type))
	}

	p._lookahead, _ = p._lexer.getNextToken()
	return token
}

func (p *Parser) getLine() int {
	return p._lexer.line
}

func (p *Parser) getCol() int {
	return p._lexer.column
}
