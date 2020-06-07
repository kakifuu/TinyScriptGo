package lexer

import (
	"TinyScriptGo/util"
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestAnalyse(t *testing.T) {
	tokens := Analyse(util.NewIterator(bufio.NewReader(strings.NewReader("(a+b)^100.12==+100-20"))))
	tests := []struct {
		index int
		want  *Token
	}{
		{0, &Token{BRACKET, "("}},
		{1, &Token{VARIABLE, "a"}},
		{2, &Token{OPERATOR, "+"}},
		{3, &Token{VARIABLE, "b"}},
		{4, &Token{BRACKET, ")"}},
		{5, &Token{OPERATOR, "^"}},
		{6, &Token{FLOAT, "100.12"}},
		{7, &Token{OPERATOR, "=="}},
		{8, &Token{INTEGER, "+100"}},
		{9, &Token{OPERATOR, "-"}},
		{10, &Token{INTEGER, "20"}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("test%v", test.index), func(t *testing.T) {
			got := tokens[test.index]
			if !got.Equals(test.want) {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestAnalyseFromFile(t *testing.T) {
	filename := "../tsg/test.tsg"
	tokens := AnalyseFromFile(filename)
	tests := []struct {
		index int
		want  *Token
	}{
		{0, &Token{KEYWORD, "func"}},
		{1, &Token{VARIABLE, "foo"}},
		{2, &Token{BRACKET, "("}},
		{3, &Token{VARIABLE, "a"}},
		{4, &Token{OPERATOR, ","}},
		{5, &Token{VARIABLE, "b"}},
		{6, &Token{BRACKET, ")"}},
		{7, &Token{BRACKET, "{"}},
		{8, &Token{VARIABLE, "print"}},
		{9, &Token{BRACKET, "("}},
		{10, &Token{VARIABLE, "a"}},
		{11, &Token{OPERATOR, "+"}},
		{12, &Token{VARIABLE, "b"}},
		{13, &Token{BRACKET, ")"}},
		{14, &Token{BRACKET, "}"}},
		{15, &Token{VARIABLE, "foo"}},
		{16, &Token{BRACKET, "("}},
		{17, &Token{FLOAT, "-100.00"}},
		{18, &Token{OPERATOR, ","}},
		{19, &Token{INTEGER, "100"}},
		{20, &Token{BRACKET, ")"}},
	}

	for _, test := range tests {
		got := tokens[test.index]
		if !got.Equals(test.want) {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}
}
