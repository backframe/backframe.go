package parser

import "fmt"

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

/**
Specfile
	: SectionList
	;
*/
func (p *Parser) Parse(val string) BfSpec {
	p._string = val
	p._lexer.init(val)
	p._lookahead, _ = p._lexer.getNextToken()

	return BfSpec{
		Sections: p.sectionList(),
	}
}

/*
SectionList
	: Section
	| SectionList Section -> Section Section
	;
*/
func (p *Parser) sectionList() []Section {
	list := []Section{p.section()}

	for {
		// check if no more lookahed
		if !p._lexer.hasMoreTokens() {
			break
		}

		list = append(list, p.section())
	}

	return list
}

/*
Section
	: 'section' 'IDENTIFIER' '{' StatementList '}'
	;

*/
func (p *Parser) section() Section {
	p._eat("SECTION")
	ident := p._eat("IDENTIFIER")
	p._eat("{")

	var t SectionType

	switch ident.value {
	case "Providers":
		t = Providers
	case "Interfaces":
		t = Interfaces
	case "Integrations":
		t = Integrations
	case "Databases":
		t = Databases
	}

	body := p.statementList(t)
	p._eat("}")
	return Section{
		_type: t,
		Body:  body,
	}
}

/**

StatementList
	: Statement
	| StatementList Statement -> Statement Statement Statement
	;
*/
func (p *Parser) statementList(t SectionType) SectionBody {
	body := SectionBody{}

	if t == Providers {
		providers := ConsumeProviderBlock(p)
		body.Providers = append(body.Providers, providers...)
	} else if t == Interfaces {
		interfaces := ConsumeInterfaceBlock(p)
		body.Interfaces = append(body.Interfaces, interfaces...)
	} else if t == Databases {
		p._eat("DATABASE")
		p._eat("IDENTIFIER")
	}

	return body
}

/**

StatementList
	: Statement
	| StatementList Statement -> Statement Statement Statement
	;

Statement
	: DeclarationStatement
	| BlockStatement
	;

DeclarationStatement
	: 'Identifier' ASSIGNMENT_OP 'Value'
	;

BlockStatement
	: 'KEYWORD' 'IDENTIFIER' '{' Statement '}'
*/

func (p *Parser) _eat(_type string) Token {
	token := p._lookahead
	if _type != token._type {
		panic(fmt.Sprintf("Syntax Error, Unexpected Token. Expected: %v but instead found: %v\n", _type, token._type))
	}

	p._lookahead, _ = p._lexer.getNextToken()
	return token
}

// Common factors
// Blocks have similar syntax
// Blocks can expect variables, or not
// Blocks have custom-defined keywords
// Blocks have nested blocks
