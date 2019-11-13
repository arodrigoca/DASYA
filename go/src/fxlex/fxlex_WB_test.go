package fxlex

import (
	"bufio"
	"testing"
	"strings"
)

func TestLexOp(t *testing.T) {

	const filename = "testfile"
  reader := bufio.NewReader(strings.NewReader(""))
  var myLexer *Lexer = NewLexer(reader, filename, true)//true indicates if debug is activated
	myLexer.accepted = append(myLexer.accepted, '+')

}

/*
func TestLexId(t *testing.T) {

}

func TestLexNum(t *testing.T) {
}

func TestLexSep(t *testing.T) {
}

func TestLexComment(t *testing.T) {
}
*/
