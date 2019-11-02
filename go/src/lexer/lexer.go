package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"errors"
	"io"
	"runtime/debug"
	"strings"
	"unicode"
)

type RuneScanner interface {
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
}

type Lexer struct {
	file     string
	line     int
	rs       RuneScanner
	lastrune rune
	accepted []rune
}

type Token struct{
	lexema string
	tokType int
	value int
	line int
}

const RuneEOF rune = 'ᛯ'

func parseArguments() (string, bool) {

	filenamePtr := flag.String("file", "", "filename to read")
	dflagPtr := flag.Bool("debug", false, "enable debug info")
	flag.Parse()
	if *filenamePtr == "" {
		fmt.Println("Error: at least argument -file is necessary.")
		os.Exit(1)
	}

	return *filenamePtr, *dflagPtr
}

func NewLexer(rs RuneScanner, filename string) (l *Lexer) {
	l = &Lexer{line: 1}
	l.file = filename
	l.rs = rs
	return l
}

func (l *Lexer) get() (r rune) {

	var err error

	rune, _, err := l.rs.ReadRune()

	if err == nil {
		l.lastrune = rune
		if rune == '\n' {
			l.line++
		}
	} else if err == io.EOF {
		l.lastrune = RuneEOF
		return RuneEOF
	}

	if err != nil {
		panic(err)
	}
	l.accepted = append(l.accepted, rune)

	return rune
}

func (l *Lexer) unget() {

	var err error

	if l.lastrune == RuneEOF {
		return
	}

	err = l.rs.UnreadRune()

	if err == nil && l.lastrune == '\n' {
		l.line--
	}
	l.lastrune = unicode.ReplacementChar
	if len(l.accepted) != 0{
		l.accepted = l.accepted[0:len(l.accepted)-1]
	}

	if err != nil {
		panic(err)
	}

}

func (l *Lexer) accept() (tok string){

	tok  = string(l.accepted)
	if tok == "" && l.lastrune != RuneEOF{
		panic(errors.New("empty token"))
	}
	l.accepted = nil
	return tok

}


func (l *Lexer) lexComment(){

	for r := l.get(); ;r = l.get(){
		fmt.Println(r)
		if r == '\n'{
			fmt.Println("end of comment")
			l.accept()
			break
		}
	}

	return
}

func (l *Lexer) Lex() (t Token, err error){

	for r := l.get(); ; r = l.get(){
		if unicode.IsSpace(r) && r != '\n'{
			l.accept()
			continue
		}
		switch r{

		case '+', '-', '*', '/', '>', '<': //operator or comment
			fmt.Println("This is an operator or a comment")
			look_token := l.get()
			if look_token == '/'{ //it's a comment
				l.lexComment()
			}else{ //not a comment so unget and continue
				l.unget()
				t.lexema = l.accept()
				return t, nil
			}

		case RuneEOF:
			l.accept()
			return t, nil

		case '\n':
			t.lexema = l.accept()
			return t, nil
		}
	}
}


func main() {

	const (
		BugMsg = "compiler error:"
		RunMsg = "runtime error:"
	)

	filename, Dflag := parseArguments()

	defer func() {
		if r := recover(); r != nil {
			errs := fmt.Sprint(r)
			if strings.HasPrefix(errs, "runtime error:") {
				errs = strings.Replace(errs, RunMsg, BugMsg, 1)
			}
			err := errors.New(errs)
			if Dflag {
				fmt.Fprintf(os.Stderr, "%s\n%s", err, debug.Stack())
			}
		}
	}()

	file, err := os.Open(filename)
	fmt.Println(file)
	if err != nil {
		fmt.Println("Error: No such file or directory.")
		return
	}

	reader := bufio.NewReader(file)
	var myLexer *Lexer = NewLexer(reader, filename)
	token, error := myLexer.Lex()
	fmt.Println(token)
	fmt.Println(error)

}
