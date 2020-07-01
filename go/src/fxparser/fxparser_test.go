package fxparser_test

import (
	"bufio"
	"flag"
	. "fxlex"
	. "fxparser"
	"os"
	"testing"
)

var filename = flag.String("file", "", "file to read")

func TestLexer(t *testing.T) {

	t.Error("dummy error")
	filename := *filename
	file, err := os.Open(filename)
	if err != nil {
		t.Log("Error: No such file or directory.")
		return
	}

	reader := bufio.NewReader(file)
	var myLexer *Lexer = NewLexer(reader, filename, true) //true indicates if debug is activated
	var myParser *Parser = NewParser(myLexer)
	parseerror := myParser.Parse()
	if parseerror != nil {
		t.Error(parseerror)
	}

}
