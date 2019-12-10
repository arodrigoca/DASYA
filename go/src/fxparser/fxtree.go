package fxparser

import (
	"fxlex"
)

type AST struct {
	FirstNode *ProgNode
}

type ProgNode struct {
	FuncNodes []*DecFuncNode
}

type DecFuncNode struct {
	Id         string
	Arguments []string
	Statements []*StatementNode
}

type StatementNode struct {
	Iter    *IterNode
	Funcall *FuncallNode
}

type FuncallNode struct {
	Id   string
	Args []fxlex.Token
}

type IterNode struct {
	Params     []fxlex.Token
	Statements []*StatementNode
}

func NewAST() (tree *AST) {
	return &AST{&ProgNode{}}
}

func NewDecFuncNode(id string) *DecFuncNode {

	var stringarray []string
	var stat []*StatementNode
	return &DecFuncNode{Id: id, Arguments: stringarray, Statements: stat}
}

func NewStatementNode() *StatementNode {

	return &StatementNode{}
}

func NewFuncallNode(id string) *FuncallNode {

	return &FuncallNode{Id: id}
}

func NewIterNode() *IterNode {

	var stat []*StatementNode
	return &IterNode{Statements: stat}
}

func (t *AST) PrintAST() {
	//
}
