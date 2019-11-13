package fxlex_test

import(
  "testing"
  "os"
  . "fxlex"
  "bufio"
  "flag"
  "strings"
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
	var myLexer *Lexer = NewLexer(reader, filename, true)//true indicates if debug is activated

	for{
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
  var myLexer *Lexer = NewLexer(reader, filename, true)//true indicates if debug is activated
  for{
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
  var myLexer *Lexer = NewLexer(reader, filename, true)//true indicates if debug is activated
  for{
    token, _ := myLexer.Lex()
    token.PrintToken()
    if token.Type == TokEof {
      break
    }
  }
}
