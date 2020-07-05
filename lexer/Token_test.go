package lexer

import (
	"TinyScriptGo/util"
	"bufio"
	"regexp"
	"strings"
	"testing"
)

func TestMakeVarOrKeyword(t *testing.T) {
	reader1 := bufio.NewReader(strings.NewReader("if abc"))
	reader2 := bufio.NewReader(strings.NewReader("true abc"))
	it1 := util.NewRuneIterator(reader1)
	it2 := util.NewRuneIterator(reader2)

	tests := []struct {
		action func()
		it     util.Iterator
		want   *Token
	}{
		{nil, it1, &Token{KEYWORD, "if"}},
		{func() { it1.Next() }, it1, &Token{VARIABLE, "abc"}},
		{nil, it2, &Token{BOOLEAN, "true"}},
	}

	for _, test := range tests {
		if test.action != nil {
			test.action()
		}
		got := MakeVarOrKeyword(test.it)
		if !got.Equals(test.want) {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}
}

func TestMakeString(t *testing.T) {
	tests := []string{"'123'", "\"123\""}
	for _, test := range tests {
		it := util.NewRuneIterator(bufio.NewReader(strings.NewReader(test)))
		got := MakeString(it)
		want := &Token{STRING, test}
		if !got.Equals(want) {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestMakeNumber(t *testing.T) {
	ptn := regexp.MustCompile("[* ]+")
	tests := []string{
		"200",
		"+0 aa",
		"-0 aa",
		"0.3 ccc",
		".5555 ddd",
		"7789.8888 ooo",
		"-1000.123123*123123",
	}
	for _, test := range tests {
		it := util.NewRuneIterator(bufio.NewReader(strings.NewReader(test)))
		got := MakeNumber(it)
		wantVal := ptn.Split(test, -1)[0]
		tokenType := INTEGER
		if strings.Contains(test, ".") {
			tokenType = FLOAT
		}
		want := &Token{tokenType, wantVal}
		if !got.Equals(want) {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestMakeOperator(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantVal string
	}{
		{"test+", "+ 100", "+"},
		{"test++", "++i", "++"},
		{"test+=", "+=g", "+="},

		{"test-", "- 100", "-"},
		{"test--", "i--", "--"},
		{"test-=", "a-=10", "-="},

		{"test*", "a * b", "*"},
		{"test*=", "a *= 10", "*="},

		{"test/", "a / 10", "/"},
		{"test/=", "a /= 10", "/="},

		{"test>", "a > b", ">"},
		{"test>>", "2 >> 1", ">>"},
		{"test>=", "2 >= 1", ">="},

		{"test<", "a < b", "<"},
		{"test<<", "1 << 2", "<<"},
		{"test<=", "1 <= 2", "<="},

		{"test=", "a = 1", "="},
		{"test==", "a == 1", "=="},

		{"test!", "!true", "!"},
		{"test!=", "!=true", "!="},

		{"test&", "you & me", "&"},
		{"test&&", "true && false", "&&"},
		{"test&=", "a &= 2", "&="},

		{"test|", "a | 2", "|"},
		{"test||", "a || b", "||"},
		{"test|=", "a |= b", "|="},

		{"test^", "a ^ b", "^"},
		{"test^^", "a ^^ b", "^^"},
		{"test^=", "a ^= b", "^="},

		{"test%", "a % b", "%"},
		{"test%=", "a %= b", "%="},

		{"test,", "a, b int", ","},
		{"test;", "abc;", ";"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			it := util.NewRuneIterator(bufio.NewReader(strings.NewReader(test.text)))
			got := MakeOperator(it)
			want := &Token{OPERATOR, test.wantVal}
			if !got.Equals(want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
