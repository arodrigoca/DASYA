package fxlex

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"unicode"
)

//token types

type TokType rune

const (
	RuneEOF = unicode.MaxRune + 1 + iota

	//reserved words
	TokFunc
	TokMain
	TokTypeDef
	TokRecord
	TokIf
	TokElse
	TokIter
	TokDefInt
	TokDefBool
	//

	TokId
	TokValInt
	TokValBool
	TokDMul
	TokGreater
	TokSmaller
	TokEqual
	TokDDEq
	TokEof = TokType(RuneEOF)
	TokBad = TokType(0)
)

//key is token lexema and value is the corrected token
var reserved_words_map = map[string]Token{

	"func":   Token{lexema: "func", Type: TokFunc},
	"main":   Token{lexema: "main", Type: TokMain},
	"type":   Token{lexema: "type", Type: TokTypeDef},
	"if":     Token{lexema: "Rect", Type: TokIf},
	"else":   Token{lexema: "else", Type: TokElse},
	"iter":   Token{lexema: "iter", Type: TokIter},
	"record": Token{lexema: "record", Type: TokRecord},
	"True":   Token{lexema: "True", Type: TokValBool, TokValBool: true},
	"False":  Token{lexema: "False", Type: TokValBool, TokValBool: false},
	"int":    Token{lexema: "int", Type: TokDefInt},
	"bool":   Token{lexema: "bool", Type: TokDefBool},
}

type Token struct {
	lexema       string
	Type         TokType
	TokValInt    int64
	TokValBool   bool
	TokValString string
	line         int
}

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
	dflag    bool
}

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

func NewLexer(rs RuneScanner, filename string, debug ...bool) (l *Lexer) {
	l = &Lexer{line: 1}
	l.file = filename
	l.rs = rs
	if v := len(debug); v > 0 {
		if debug[0] {
			l.dflag = true
		}
	}
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
		switch r {
		case '\n':
			l.unget()
			l.accept()
			return
		}
	}
}

func (l *Lexer) lexOp() (t Token, err error) {

	const (
		ops = "+-*/><=%|&!^=:"
	)

	//special case. RuneScanner doesn't allow to use unget() twice in a row (line 403), so
	//the lookahead rune and the rune before that cannot be ungetted (comment case).

	if string(l.accepted) == "/" {
		t.lexema = l.accept()
		t.Type = TokType('/')
		return t, nil
	}

	r := l.get()

	switch r {

	case '*':
		look_token := l.get()
		if look_token == '*' {
			//power operator
			t.lexema = l.accept()
			t.Type = TokDMul

		} else {
			l.unget()
			t.lexema = l.accept()
			t.Type = TokType(r)
		}

	case '>', '<':
		look_token := l.get()
		if look_token == '=' {
			//comparison operator
			t.lexema = l.accept()
			if r == '>' {
				t.Type = TokGreater
			} else {
				t.Type = TokSmaller
			}

		} else {
			l.unget()
			t.lexema = l.accept()
			t.Type = TokType(r)
		}

	case '=':
		//equality operator
		look_token := l.get()
		if look_token == '=' {
			t.lexema = l.accept()
			t.Type = TokEqual
		} else {
			l.unget()
			t.lexema = l.accept()
			t.Type = TokType(r)
		}

	case ':':
		look_token := l.get()
		if look_token == '=' {
			t.lexema = l.accept()
			t.Type = TokDDEq
		} else {
			l.unget()
			t.lexema = l.accept()
			t.Type = TokType(r)
		}

	default:
		if strings.ContainsRune(ops, r) {
			//correct operator
			t.lexema = l.accept()
			t.Type = TokType(r)
		} else {
			panic(errors.New("Bad operator"))
		}

	}
	return t, err

}

func (l *Lexer) lexId() (t Token, err error) {

	for r := l.get(); ; r = l.get() {

		switch {

		case unicode.IsLetter(r), strings.ContainsRune("_", r), unicode.IsNumber(r):

			continue

		default:

			l.unget()
			t.lexema = l.accept()
			t.Type = TokId
			t.TokValString = t.lexema
			break

		}

		if value, found := reserved_words_map[t.lexema]; found {
			t = value
		}
		return t, nil
	}
}

func (l *Lexer) lexSep() (t Token, err error) {

	const sep = "(),;[]{}."

	r := l.get()

	if strings.ContainsRune(sep, r) {
		t.lexema = l.accept()
		t.Type = TokType(r)
		return t, nil

	} else {
		panic(errors.New("Bad separator"))
	}

}

func (l *Lexer) lexNum() (t Token, err error) {

	const validHex = "ABCDEFabcdef"

	var isHex = false

	for r := l.get(); ; r = l.get() {

		switch {

		case r == '0':
			if !isHex {
				//fmt.Println("Number is zero. Might be hexadecimal")
				look_token := l.get()
				if look_token == 'x' {
					//fmt.Println("Hexadecimal")
					isHex = true
				} else {
					l.unget()
				}
			}

		case unicode.IsNumber(r):
			//fmt.Println("rune is an integer")

		case strings.ContainsRune(validHex, r):
			//fmt.Println("rune is a hexadecimal letrter")

		default:
			l.unget()
			t.lexema = l.accept()
			t.Type = TokValInt
			t.TokValInt, err = strconv.ParseInt(t.lexema, 10, 64)
			return t, nil
		}
	}
}

func (l *Lexer) Lex() (t Token, err error) {

	const (
		BugMsg = "compiler error:"
		RunMsg = "runtime error:"
	)

	defer func() {
		if e := recover(); e != nil {
			errs := fmt.Sprint(e)
			if strings.HasPrefix(errs, "runtime error:") {
				errs = strings.Replace(errs, RunMsg, BugMsg, 1)
			}
			err := errors.New(errs)
			if l.dflag {
				fmt.Fprintf(os.Stderr, "%s\n%s", err, debug.Stack())
			}
		}
	}()

	for r := l.get(); ; r = l.get() {
		if unicode.IsSpace(r) {
			l.accept()
			continue
		}
		switch r {

		case '+', '-', '*', '/', '>', '<', '=', ':', '%', '|', '&', '!', '^': //operator or comment

			if r == '/' {
				look_token := l.get()
				if look_token == '/' { //it's a comment
					l.lexComment()
				} else { //not a comment so unget and continue
					l.unget()
					t, err = l.lexOp()
					t.line = l.line
					return t, err
				}
			} else {
				l.unget()
				t, err = l.lexOp()
				t.line = l.line
				return t, err
			}

		case RuneEOF:
			//fmt.Println("End of file")
			t.lexema = l.accept()
			t.Type = TokEof
			t.line = l.line
			//fmt.Println(t.lexema)
			return t, nil

		case '\n':
			l.accept()
			continue

		case '(', ')', ',', ';', '[', ']', '{', '}', '.':

			l.unget()
			t, err = l.lexSep()
			t.line = l.line
			return t, err

		default:
			break
		}

		switch {

		case unicode.IsLetter(r):
			l.unget()
			t, err = l.lexId()
			t.line = l.line
			return t, err

		case unicode.IsNumber(r):
			l.unget()
			t, err = l.lexNum()
			t.line = l.line
			return t, err
		}
	}
}

func (t *Token) PrintToken() {

	if t.Type > unicode.MaxRune {
		switch t.Type {

		case TokEof:
			//fmt.Println("End of file token")
			fmt.Printf("Token type: TokEof\n")

		case TokId:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokId\n")
			fmt.Printf("Value: %s\n", t.TokValString)

		case TokValInt:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokValInt\n")
			fmt.Printf("Value: %v\n", t.TokValInt)

		case TokValBool:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokValBool\n")
			fmt.Printf("Value: %v\n", t.TokValBool)

		case TokFunc:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokFunc\n")

		case TokMain:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokMain\n")

		case TokTypeDef:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokTypeDef\n")

		case TokRecord:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokRecord\n")

		case TokIf:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokIf\n")

		case TokElse:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokElse\n")

		case TokIter:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokIter\n")

		case TokDMul:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokDMul\n")

		case TokGreater:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokGreater\n")

		case TokSmaller:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokSmaller\n")

		case TokEqual:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokEqual\n")

		case TokDDEq:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokDDEq\n")

		case TokDefInt:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokDefInt\n")

		case TokDefBool:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: TokDefBool\n")

		default:
			fmt.Printf("Lexema: %s\n", t.lexema)
			fmt.Printf("Token type: %v\n", t.Type)
		}
	} else {
		fmt.Printf("Lexema: %s\n", t.lexema)
		fmt.Printf("Token type: %c\n", t.Type)
	}
	fmt.Printf("Line: %v\n", t.line)
	fmt.Printf("\n")

}
