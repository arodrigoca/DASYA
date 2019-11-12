package fxlex

import (
	"bufio"
	"os"
	"testing"
)

func TestLexOp(t *testing.T) {

	file, _ := os.Open("../../bin/lexer_test.txt")
	reader := bufio.NewReader(file)
	var myLexer *Lexer = NewLexer(reader, "file")
	token, _ := myLexer.Lex()
	t.Log(token.lexema)

}
