package lexer

import (
	"TinyScriptGo/util"
	"errors"
)

type Token struct {
	_type  int
	_value string
}

// Token Type
const (
	KEYWORD = iota
	VARIABLE
	BRACKET
	OPERATOR
	INTEGER
	FLOAT
	STRING
	BOOLEAN
)

var (
	ErrUnexpected           = errors.New("Unexpected error ")
	ErrUnexpectedCharFormat = "Unexpected character %v "
)

func (token *Token) GetType() int {
	return token._type
}

func (token *Token) IsVariable() bool {
	return token._type == VARIABLE
}

func (token *Token) IsScalar() bool {
	return token._type == STRING || token._type == INTEGER ||
		token._type == FLOAT || token._type == BOOLEAN
}

func (token *Token) IsValue() bool {
	return token.IsScalar() || token.IsVariable()
}

func (token *Token) IsOperator() bool {
	return token._type == OPERATOR
}

func (token *Token) IsNumber() bool {
	return token._type == INTEGER || token._type == FLOAT
}

func (token *Token) Equals(t *Token) bool {
	return token._type == t._type && token._value == t._value
}

func MakeVarOrKeyword(it *util.Iterator) *Token {
	value := ""
	for it.HasNext() {
		p := it.Peek()
		if util.IsLiteral(p) {
			value += p
		} else {
			break
		}
		it.Next()
	}

	if IsKeyword(value) {
		return &Token{KEYWORD, value}
	}
	if value == "false" || value == "true" {
		return &Token{BOOLEAN, value}
	}

	return &Token{VARIABLE, value}
}

func MakeString(it *util.Iterator) *Token {
	value := ""
	state := 0
	for it.HasNext() {
		s := it.Next()
		switch state {
		case 0:
			if s == "\"" {
				state = 1
			} else {
				state = 2
			}
			value += s
		case 1:
			if s == "\"" {
				return &Token{STRING, value + s}
			} else {
				value += s
			}
		case 2:
			if s == "'" {
				return &Token{STRING, value + s}
			} else {
				value += s
			}
		}
	}

	panic(ErrUnexpected)
}

func MakeNumber(it *util.Iterator) *Token {
	value := ""
	state := 0
	for it.HasNext() {
		p := it.Peek()

		switch state {
		case 0:
			if p == "0" {
				state = 1
			} else if util.IsNumber(p) {
				state = 2
			} else if p == "+" || p == "-" {
				state = 3
			} else if p == "." {
				state = 5
			}
		case 1:
			if p == "0" {
				state = 1
			} else if util.IsNumber(p) {
				state = 2
			} else if p == "." {
				state = 4
			} else {
				return &Token{INTEGER, "0"}
			}
		case 2:
			if util.IsNumber(p) {
				state = 2
			} else if p == "." {
				state = 4
			} else {
				return &Token{INTEGER, value}
			}
		case 3:
			if util.IsNumber(p) {
				state = 2
			} else if p == "." {
				state = 5
			} else {
				panic(util.ErrWithArgs(ErrUnexpectedCharFormat, p))
			}
		case 4:
			if util.IsNumber(p) {
				state = 20
			} else if p == "." {
				panic(util.ErrWithArgs(ErrUnexpectedCharFormat, p))
			} else {
				return &Token{FLOAT, value}
			}
		case 5:
			if util.IsNumber(p) {
				state = 20
			} else {
				panic(util.ErrWithArgs(ErrUnexpectedCharFormat, p))
			}
		case 20:
			if util.IsNumber(p) {
				state = 20
			} else if p == "." {
				panic(util.ErrWithArgs(ErrUnexpectedCharFormat, p))
			} else {
				return &Token{FLOAT, value}
			}
		}

		it.Next()
		value += p
	}

	panic(ErrUnexpected)
}

func MakeOperator(it *util.Iterator) *Token {
	state := 0

	for it.HasNext() {
		s := it.Next()

		switch state {
		case 0:
			switch s {
			case "+":
				state = 1
			case "-":
				state = 2
			case "*":
				state = 3
			case "/":
				state = 4
			case ">":
				state = 5
			case "<":
				state = 6
			case "=":
				state = 7
			case "!":
				state = 8
			case "&":
				state = 9
			case "|":
				state = 10
			case "^":
				state = 11
			case "%":
				state = 12
			case ",":
				return &Token{OPERATOR, ","}
			case ";":
				return &Token{OPERATOR, ";"}
			}
		case 1:
			if s == "+" {
				return &Token{OPERATOR, "++"}
			} else if s == "=" {
				return &Token{OPERATOR, "+="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "+"}
			}
		case 2:
			if s == "-" {
				return &Token{OPERATOR, "--"}
			} else if s == "=" {
				return &Token{OPERATOR, "-="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "-"}
			}
		case 3:
			if s == "=" {
				return &Token{OPERATOR, "*="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "*"}
			}
		case 4:
			if s == "=" {
				return &Token{OPERATOR, "/="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "/"}
			}
		case 5:
			if s == "=" {
				return &Token{OPERATOR, ">="}
			} else if s == ">" {
				return &Token{OPERATOR, ">>"}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, ">"}

			}
		case 6:
			if s == "=" {
				return &Token{OPERATOR, "<="}
			} else if s == "<" {
				return &Token{OPERATOR, "<<"}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "<"}
			}
		case 7:
			if s == "=" {
				return &Token{OPERATOR, "=="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "="}
			}
		case 8:
			if s == "=" {
				return &Token{OPERATOR, "!="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "!"}
			}
		case 9:
			if s == "&" {
				return &Token{OPERATOR, "&&"}
			} else if s == "=" {
				return &Token{OPERATOR, "&="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "&"}
			}
		case 10:
			if s == "|" {
				return &Token{OPERATOR, "||"}
			} else if s == "=" {
				return &Token{OPERATOR, "|="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "|"}
			}
		case 11:
			if s == "^" {
				return &Token{OPERATOR, "^^"}
			} else if s == "=" {
				return &Token{OPERATOR, "^="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "^"}
			}
		case 12:
			if s == "=" {
				return &Token{OPERATOR, "%="}
			} else {
				it.PutBack(1)
				return &Token{OPERATOR, "%"}
			}
		}
	}

	panic(ErrUnexpected)
}
