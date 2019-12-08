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
	Statements []*StatementNode
}

type StatementNode struct {
	Iter    *IterNode
	Funcall *FuncallNode
}

type FuncallNode struct {
	Id   fxlex.Token
	Args []*fxlex.Token
}

type IterNode struct {
	Params     []*fxlex.Token
	Statements []*StatementNode
}

func NewAST() (tree *AST) {
	return &AST{&ProgNode{}}
}

func NewDecFuncNode(id string) *DecFuncNode {

	return &DecFuncNode{Id: id}
}

func NewStatementNode(iter *IterNode, funcall *FuncallNode) *StatementNode {
	return &StatementNode{Iter: iter, Funcall: funcall}
}
