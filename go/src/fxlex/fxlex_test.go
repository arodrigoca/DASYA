package fxlex_test

import(
  "testing"
  "os"
  . "fxlex"
  "bufio"
)

func TestLexer(t *testing.T) {

  filename := "../../bin/lang.fx"

  file, err := os.Open(filename)
	if err != nil {
		t.Log("Error: No such file or directory.")
		return
	}

	reader := bufio.NewReader(file)
	var myLexer *Lexer = NewLexer(reader, filename)

	for i := 0; i <= 200; i++ {
		token, _ := myLexer.Lex()
		//token.PrintToken()
		if token.Type == TokEof {
			break
		}
	}

}
