package lexer

import (
	"TinyScriptGo/util"
	"bufio"
	"errors"
	"os"
)

var (
	ErrCommentsNotMatch = errors.New("Comments not match ")
)

func Analyse(it *util.Iterator) []*Token {
	var tokens []*Token

	for it.HasNext() {
		s := it.Next()
		if s == util.END_SYM {
			break
		}
		if s == " " || s == "\n" {
			continue
		}

		p := it.Peek()

		if s == "/" {
			if p == "/" {
				for it.HasNext() {
					if s = it.Next(); s == "\n" {
						break
					}
				}
				continue
			} else if p == "*" {
				it.Next()
				valid := false
				for it.HasNext() {
					tmp := it.Next()
					if tmp == "*" && it.Peek() == "/" {
						it.Next()
						valid = true
						break
					}
				}
				if !valid {
					panic(ErrCommentsNotMatch)
				}
				continue
			}
		}

		if s == "{" || s == "}" || s == "(" || s == ")" {
			tokens = append(tokens, &Token{BRACKET, s})
			continue
		}

		if s == "\"" || s == "'" {
			it.PutBack(1)
			tokens = append(tokens, MakeString(it))
			continue
		}

		if util.IsLetter(s) {
			it.PutBack(1)
			tokens = append(tokens, MakeVarOrKeyword(it))
			continue
		}

		if util.IsNumber(s) {
			it.PutBack(1)
			tokens = append(tokens, MakeNumber(it))
			continue
		}

		if (s == "+" || s == "-" || s == ".") && util.IsNumber(p) {
			var lastToken *Token
			if len(tokens) > 0 {
				lastToken = tokens[len(tokens)-1]
			}
			if lastToken == nil || !lastToken.IsValue() || lastToken.IsOperator() {
				it.PutBack(1)
				tokens = append(tokens, MakeNumber(it))
				continue
			}
		}

		if util.IsOperator(s) {
			it.PutBack(1)
			tokens = append(tokens, MakeOperator(it))
			continue
		}

		panic(util.ErrWithArgs(ErrUnexpectedCharFormat, s))
	}

	return tokens
}

func AnalyseFromFile(filename string) []*Token {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return Analyse(util.NewIterator(bufio.NewReader(file)))
}
