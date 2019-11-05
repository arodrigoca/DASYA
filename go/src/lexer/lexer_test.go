package main

import (
	"bufio"
	"os"
	"testing"
)

func TestCMDline(t *testing.T) {

	os.Args = []string{"cmd", "-file ../../bin/lang.fx", "-debug"}
	filename, Dflag := parseArguments()
	t.Log(filename)
	t.Log(Dflag)

}

func TestLexOp(t *testing.T) {

	file, _ := os.Open("../../bin/lexer_test.txt")
	reader := bufio.NewReader(file)
	var myLexer *Lexer = NewLexer(reader, "file")
	token, _ := myLexer.Lex()
	t.Log(token.lexema)

}
