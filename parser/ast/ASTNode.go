package ast

import (
	"TinyScriptGo/lexer"
	"errors"
)

const (
	BLOCK = iota
	BINARY_EXPR
	UNARY_EXPR
	CALL_EXPR
	VARIABLE
	SCALAR
	IF_STMT
	FOR_STMT
	WHILE_STMT
	RETURN_STMT
	ASSIGN_STMT
	DECLARE_STMT
	FUNC_DECLARE_STMT
)

var (
	ErrIndexOutOfBound = errors.New("Index out of bound ")
)

type Node struct {
	parent   *Node
	children []*Node
	lexeme   *lexer.Token
	label    string
	typ      int
	props    map[string]interface{}
}

func (node *Node) AddChild(child *Node) {
	child.parent = node
	node.children = append(node.children, child)
}

func (node *Node) GetChild(index int) *Node {
	if index >= len(node.children) {
		return nil
	}
	return node.children[index]
}

func (node *Node) GetChildren() []*Node {
	return node.children
}

func (node *Node) ReplaceChild(index int, child *Node) error {
	if index >= len(node.children) {
		return ErrIndexOutOfBound
	}
	node.children[index] = child
	return nil
}

func (node *Node) SetLexeme(lexeme *lexer.Token) {
	node.lexeme = lexeme
}

func (node *Node) GetLexeme() *lexer.Token {
	return node.lexeme
}

func (node *Node) SetLabel(label string) {
	node.label = label
}

func (node *Node) GetLabel() string {
	return node.label
}

func (node *Node) SetType(typ int) {
	node.typ = typ
}

func (node *Node) GetType() int {
	return node.typ
}

func (node *Node) SetProp(key string, val interface{}) {
	node.props[key] = val
}

func (node *Node) GetProps() map[string]interface{} {
	return node.props
}

func (node *Node) GetProp(key string) interface{} {
	if val, ok := node.props[key]; ok {
		return val
	}
	return nil
}
