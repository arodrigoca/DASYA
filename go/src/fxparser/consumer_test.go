package fxparser_test

import (
	"bufio"
	"errors"
	. "fxlex"
	. "fxparser"
	"strings"
	"testing"
)

func TestConsumer(t *testing.T) {

	var test_text string = "circle();"

	fake_reader := strings.NewReader(test_text)

	reader := bufio.NewReader(fake_reader)
	var myLexer *Lexer = NewLexer(reader, "consumer_test.txt", true) //true indicates if debug is activated
	var myParser *Parser = NewParser(myLexer)
	err := errors.New("test error")
	parseerror := myParser.ConsumeUntilMarker(err)
	if parseerror != nil {
		t.Error(parseerror)
	}

}

func TestConsumer(t *testing.T) {

	var test_text string = "circle();"

	fake_reader := strings.NewReader(test_text)

	reader := bufio.NewReader(fake_reader)
	var myLexer *Lexer = NewLexer(reader, "consumer_test.txt", true) //true indicates if debug is activated
	var myParser *Parser = NewParser(myLexer)
	err := errors.New("test error")
	parseerror := myParser.ConsumeUntilMarker(err)
	if parseerror != nil {
		t.Error(parseerror)
	}

}
