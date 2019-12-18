package fxparser_test

import (
	"bufio"
	///"errors"
	. "fxlex"
	. "fxparser"
	"strings"
	"testing"
)

func TestConsumer(t *testing.T) {

	var test_text string = "circle);\ncircle();"

	fake_reader := strings.NewReader(test_text)

	reader := bufio.NewReader(fake_reader)
	var myLexer *Lexer = NewLexer(reader, "consumer_test.txt", true) //true indicates if debug is activated
	var myParser *Parser = NewParser(myLexer)
	parseerror := myParser.ConsumeUntilMarker("{}();")
	if parseerror != nil {
		t.Error(parseerror)
	}

}

func TestErrExpected(t *testing.T) {

	//var test_text string = "func line ( int x , int y ){iter (i := 0; x , 1){circle (2 , 3, y , 5);}}"
	var test_text string = "func line ( int x , int y ){circle);\ncircle();}"
	fake_reader := strings.NewReader(test_text)
	reader := bufio.NewReader(fake_reader)
	var myLexer *Lexer = NewLexer(reader, "consumer_test.txt", true) //true indicates if debug is activated
	var myParser *Parser = NewParser(myLexer)

	parseerror := myParser.Parse()
	if parseerror != nil {
		t.Error(parseerror)
	}

}
