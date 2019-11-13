package fxlex

import (
	"bufio"
	"testing"
	"strings"
)

func TestLexOp(t *testing.T) {

	//This code test if operators are detected correctly. Then, because of the unexpected ,
	//lexOp raises a panic, which is the expected behaviour (should never happen)
	defer func() {
		if r := recover(); r != nil {
	  	t.Log("Bad op")
	  }
	}()

	const filename = "testfile"
  reader := bufio.NewReader(strings.NewReader("+>=<=:===-4"))
  var myLexer *Lexer = NewLexer(reader, filename, true)//true indicates if debug is activated
	for{
		token, _ := myLexer.lexOp()
		token.PrintToken()
    if token.Type == TokEof {
      break
    }
  }

}


func TestLexNum(t *testing.T) {

	//test if it raises a panic with the unexpected , in the int
	defer func() {
		if r := recover(); r != nil {
	  	t.Log("Bad int")
	  }
	}()

	const filename = "testfile"
  reader := bufio.NewReader(strings.NewReader("0xff,"))
  var myLexer *Lexer = NewLexer(reader, filename, true)//true indicates if debug is activated
	for{
		token, _ := myLexer.lexNum()
		token.PrintToken()
    if token.Type == TokEof {
      break
    }
  }

}

func TestLexSep(t *testing.T) {

	//should raise a panic when it reaches the 1
	defer func() {
		if r := recover(); r != nil {
	  	t.Log("Bad separator")
	  }
	}()

	const filename = "testfile"
  reader := bufio.NewReader(strings.NewReader("(),1"))
  var myLexer *Lexer = NewLexer(reader, filename, true)//true indicates if debug is activated
	for{
		token, _ := myLexer.lexSep()
		token.PrintToken()
    if token.Type == TokEof {
      break
    }
  }
}

func TestLexComment(t *testing.T) {

	const filename = "testfile"
  reader := bufio.NewReader(strings.NewReader("//comment\nhello"))
  var myLexer *Lexer = NewLexer(reader, filename, true)//true indicates if debug is activated
	myLexer.lexComment()
}
