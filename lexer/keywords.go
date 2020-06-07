package lexer

var keywords = map[string]int{
	"var":    0,
	"int":    0,
	"float":  0,
	"bool":   0,
	"void":   0,
	"string": 0,
	"if":     0,
	"else":   0,
	"for":    0,
	"while":  0,
	"break":  0,
	"func":   0,
	"return": 0,
}

func IsKeyword(str string) bool {
	_, exist := keywords[str]
	return exist
}
