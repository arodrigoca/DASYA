package fxlex_test

import (
	"bufio"
	"flag"
	. "fxlex"
	"os"
	"strings"
	"testing"
)

var filename = flag.String("file", "", "file to read")

func TestLexer(t *testing.T) {

	//to run this specific test:
	//go test -v -args -file "<filepath>"
	//otherwise, if this test file is ran with go test, the rest of the
	//test file will be executed

	filename := *filename
	file, err := os.Open(filename)
	if err != nil {
		t.Log("Error: No such file or directory.")
		return
	}

	reader := bufio.NewReader(file)
	var myLexer *Lexer = NewLexer(reader, filename, true) //true indicates if debug is activated

	for {
		token, _ := myLexer.Lex()
		token.PrintToken()
		if token.Type == TokEof {
			break
		}
	}

}

func TestVoidFile(t *testing.T) {

	const filename = "testfile"
	reader := bufio.NewReader(strings.NewReader(""))
	var myLexer *Lexer = NewLexer(reader, filename, true) //true indicates if debug is activated
	for {
		token, _ := myLexer.Lex()
		token.PrintToken()
		if token.Type == TokEof {
			break
		}
	}

}

func TestStringFile(t *testing.T) {

	const filename = "testfile"
	reader := bufio.NewReader(strings.NewReader("//comment\ntype record vector(int x, int y, int z)\ntpe recrd vector intx int,y, intz"))
	var myLexer *Lexer = NewLexer(reader, filename, true) //true indicates if debug is activated
	for {
		token, _ := myLexer.Lex()
		token.PrintToken()
		if token.Type == TokEof {
			break
		}
	}
}

func TestLexOp(t *testing.T) {

	const filename = "testfile"
	reader := bufio.NewReader(strings.NewReader("+-*/><=%|&!^=:<=>=:="))
	var myLexer *Lexer = NewLexer(reader, filename, true) //true indicates if debug is activated
	for {
		token, _ := myLexer.Lex()
		token.PrintToken()
		if token.Type == TokEof {
			break
		}
	}
}

func TestLexNum(t *testing.T) {

	const filename = "testfile"
	reader := bufio.NewReader(strings.NewReader("0xff 48"))
	var myLexer *Lexer = NewLexer(reader, filename, true) //true indicates if debug is activated
	for {
		token, _ := myLexer.Lex()
		token.PrintToken()
		if token.Type == TokEof {
			break
		}
	}

}

func TestLexSep(t *testing.T) {

	const filename = "testfile"
	reader := bufio.NewReader(strings.NewReader("( ) , ; [ ] { } ."))
	var myLexer *Lexer = NewLexer(reader, filename, true) //true indicates if debug is activated
	for {
		token, _ := myLexer.Lex()
		token.PrintToken()
		if token.Type == TokEof {
			break
		}
	}
}

func TestLexComment(t *testing.T) {

	const filename = "testfile"
	reader := bufio.NewReader(strings.NewReader("//comment\n"))
	var myLexer *Lexer = NewLexer(reader, filename, true) //true indicates if debug is activated
	myLexer.Lex()
}

func TestLexId(t *testing.T) {

	const filename = "testfile"
	reader := bufio.NewReader(strings.NewReader("abcd1234_0xff"))
	var myLexer *Lexer = NewLexer(reader, filename, true) //true indicates if debug is activated
	for {
		token, _ := myLexer.Lex()
		token.PrintToken()
		if token.Type == TokEof {
			break
		}
	}
}
