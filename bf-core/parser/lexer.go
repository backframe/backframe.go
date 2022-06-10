package parser

import (
	"errors"
	"regexp"
)

var Spec = []([2]string){
	// Whitespace and comments
	{`^\s+`, "nil"},
	{`^#.*`, "nil"},
	{`^\/\/.*`, "nil"},

	// Delimiters and symbols
	{`^\{`, "{"},
	{`^\}`, "}"},
	{`^;`, ";"},
	{`^\[`, "["},
	{`^\]`, "]"},

	// Numbers
	{`^\d+`, "NUMBER"},

	// Strings
	{`^"[^"]*"`, "STRING"},
	{`^'[^']*'`, "STRING"},

	// Operators
	{`^=`, "ASSIGNMENT_OP"},
	{`^,`, "COMMA"},

	{`^\btrue\b`, "BOOLEAN"},
	{`^\bfalse\b`, "BOOLEAN"},

	// Keywords
	{`^\bsection\b`, "SECTION"},
	{`^\bprovider\b`, "PROVIDER"},
	{`^\binterface\b`, "INTERFACE"},
	{`^\bdatabase\b`, "DATABASE"},
	{`^\bintegration\b`, "INTEGRATION"},
	{`^\bresource\b`, "RESOURCE"},
	{`^\bmethod\b`, "METHOD"},
	{`^\bmigration\b`, "MIGRATION"},
	{`^\bschemaDef\b`, "SCHEMA"},
	{`^\bversion\b`, "VERSION"},

	// Others
	{`^\w+`, "IDENTIFIER"},
}

type Lexer struct {
	_cursor int
	_string string
}

type Token struct {
	_type string
	value string
}

func (l *Lexer) init(val string) {
	l._cursor = 0
	l._string = val
}

func (l *Lexer) hasMoreTokens() bool {
	return l._cursor < len(l._string)
}

func (l *Lexer) getNextToken() (Token, error) {
	if !l.hasMoreTokens() {
		return Token{}, errors.New("no more tokens")
	}

	str := l._string[l._cursor:]
	for _, v := range Spec {
		re, _type := v[0], v[1]
		regex, _ := regexp.Compile(re)
		value := l._match(regex, str)

		if value == "" {
			continue
		}

		if _type == "nil" {
			return l.getNextToken()
		} else {
			return Token{
				_type: _type,
				value: value,
			}, nil
		}

	}

	return Token{}, nil
}

func (l *Lexer) _match(re *regexp.Regexp, val string) string {
	matched := re.FindAllString(val, 1)
	if len(matched) == 0 {
		return ""
	}
	l._cursor += len(matched[0])

	return matched[0]
}
