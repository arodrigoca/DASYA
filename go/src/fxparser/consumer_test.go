package fxparser_test

import(
  "testing"
  "bufio"
  . "fxlex"
  . "fxparser"
  "strings"
)


func TestConsumer(t *testing.T){

  var test_text string = "circle();"

  fake_reader := strings.NewReader(test_text)

	reader := bufio.NewReader(fake_reader)
	var myLexer *Lexer = NewLexer(reader, "consumer_test.txt", true) //true indicates if debug is activated
	var myParser *Parser = NewParser(myLexer)
	parseerror := myParser.ConsumeUntilMarker()
	if parseerror != nil {
		t.Error(parseerror)
	}

}
