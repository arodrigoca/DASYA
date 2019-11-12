package fxlex_test

import(
  "testing"
  "os"
  . "fxlex"
  "bufio"
  "flag"
)

var filename = flag.String("file", "", "file to read")

func TestLexerBasic(t *testing.T) {

  filename := *filename
  file, err := os.Open(filename)
	if err != nil {
		t.Log("Error: No such file or directory.")
		return
	}

	reader := bufio.NewReader(file)
	var myLexer *Lexer = NewLexer(reader, filename, true)

	for{
		token, _ := myLexer.Lex()
		token.PrintToken()
		if token.Type == TokEof {
			break
		}
	}

}
