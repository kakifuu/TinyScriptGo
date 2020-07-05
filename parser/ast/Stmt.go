package ast

import "TinyScriptGo/util"

type Stmt Node

func (stmt *Stmt) ParseStmt(it util.Iterator) (*Node, error) {
	if !it.HasNext() {
		return nil, nil
	}
	return nil, nil
}
