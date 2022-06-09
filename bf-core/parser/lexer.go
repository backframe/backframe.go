package parser

import (
	"errors"
	"fmt"
	"regexp"
)

var tokens = []([2]string){
	{`^\s+`, "nil"},
	{`^section`, "SECTION"},
}

type Tokenizer struct {
	_string string
	_cursor int
}

func (t *Tokenizer) init(val string) {
	t._cursor = 0
	t._string = val
}

func (t *Tokenizer) hasMoreTokens() bool {
	return t._cursor < len(t._string)
}

func (t *Tokenizer) getNextToken() (Token, error) {
	if !t.hasMoreTokens() {
		return Token{}, errors.New("no more tokens")
	}

	str := t._string[:t._cursor]

	for _, v := range tokens {
		re, _type := v[0], v[1]
		fmt.Println(re)
		fmt.Println(_type)
		regex, _ := regexp.Compile(re)
		value := t._match(regex, str)

		if _type == "nil" {
			return t.getNextToken()
		} else {
			return Token{
				_type: _type,
				value: value,
			}, nil
		}

	}

	panic("Syntax Error, unexpected token")
}

func (t *Tokenizer) _match(re *regexp.Regexp, val string) string {
	fmt.Println(val)
	fmt.Println(re.Match([]byte(val)))
	matched := re.FindAllString(val, 1)
	if matched[0] == "" {
		panic("no match")
	}
	t._cursor += len(matched[0])

	return matched[0]
}

type Token struct {
	_type string
	value string
}

type Parser struct {
	_string    string
	_tokenizer Tokenizer
	_lookahead Token
}

func (p *Parser) Parse(val string) Spec {
	p._string = val
	p._tokenizer.init(val)
	p._lookahead, _ = p._tokenizer.getNextToken()

	sections := make(map[SectionT]Section)
	s := p.section()
	sections[s.sectionType] = s

	return Spec{
		sections: sections,
	}
}

func (p *Parser) section() Section {
	p._eat("SECTION")
	return Section{
		sectionType: Providers,
	}
}

func (p *Parser) _eat(t string) Token {
	token := p._lookahead
	fmt.Printf("%v", token)
	if t != token._type {
		panic("Errorrrrrrr")
	}

	return token
}
