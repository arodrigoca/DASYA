package fxparser_test

import (
	"bufio"
	. "fxlex"
  . "fxparser"
	"os"
	"testing"
  "flag"
  "fmt"
)

var filename = flag.String("file", "", "file to read")

func TestLexer(t *testing.T) {

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
  fmt.Println(parseerror)

}
