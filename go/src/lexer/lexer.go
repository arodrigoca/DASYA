package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
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

type Token struct {
	lexema  string
	tokType int
	value   int
	line    int
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
	if len(l.accepted) != 0 {
		l.accepted = l.accepted[0 : len(l.accepted)-1]
	}

	if err != nil {
		panic(err)
	}

}

func (l *Lexer) accept() (tok string) {

	tok = string(l.accepted)
	if tok == "" && l.lastrune != RuneEOF {
		panic(errors.New("empty token"))
	}
	l.accepted = nil
	return tok

}

func (l *Lexer) lexComment() {

	for r := l.get(); ; r = l.get() {
		//fmt.Println(r)
		if r == '\n' {
			//fmt.Println("end of comment")
			l.accept()
			break
		} else {
			l.accept()
		}
	}

	return
}

func (l *Lexer) lexOp() (t Token, err error) {

	const (
		ops = "+-*/><="
	)

	r := l.get()

	switch r {

	case '*':
		look_token := l.get()
		if look_token == '*' {
			//power operator
			t.lexema = l.accept()

		} else {
			l.unget()
			t.lexema = l.accept()
		}

	case '>', '<':
		look_token := l.get()
		if look_token == '=' {
			//comparison operator
			t.lexema = l.accept()

		} else {
			l.unget()
			t.lexema = l.accept()
		}

	case '=':
		look_token := l.get()
		if look_token == '=' {
			t.lexema = l.accept()
		} else {
			l.unget()
			t.lexema = l.accept()
		}

	default:
		if strings.ContainsRune(ops, r) {
			//correct operator
			t.lexema = l.accept()
		} else {
			panic(errors.New("Bad operator"))
		}

	}

	return t, err

}

func (l *Lexer) Lex() (t Token, err error) {

	for r := l.get(); ; r = l.get() {
		if unicode.IsSpace(r) && r != '\n' {
			l.accept()
			continue
		}
		switch r {

		case '+', '-', '*', '/', '>', '<': //operator or comment

			if r == '/' {
				look_token := l.get()
				if look_token == '/' { //it's a comment
					l.lexComment()
				} else { //not a comment so unget and continue
					l.unget()
					t, err = l.lexOp()
					return t, nil
				}
			} else {
				l.unget()
				t, err = l.lexOp()
				return t, nil
			}

		case RuneEOF:
			l.accept()
			return t, nil

		case '\n':
			t.lexema = l.accept()
			fmt.Println("Lexemma is line end")
			return t, nil

		default:
			fmt.Println("Not an operator or eof or eol")
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
	token, _ := myLexer.Lex()
	fmt.Println(token)

	token, _ = myLexer.Lex()
	fmt.Println(token)

	token, _ = myLexer.Lex()
	fmt.Println(token)

}
