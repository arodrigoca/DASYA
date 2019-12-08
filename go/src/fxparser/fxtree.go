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
	return &DecFuncNode{Id: id, Arguments: stringarray}
}

func NewStatementNode(iter *IterNode, funcall *FuncallNode) *StatementNode {

	return &StatementNode{Iter: iter, Funcall: funcall}
}

func NewFuncallNode(id string, args []fxlex.Token) *FuncallNode {

	return &FuncallNode{Id: id, Args: args}
}

func NewIterNode(params []fxlex.Token, stat []*StatementNode) *IterNode {

	return &IterNode{Params: params, Statements: stat}
}

func (t *AST) PrintAST() {
	//
}
