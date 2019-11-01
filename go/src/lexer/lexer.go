package main

import (
	"flag"
	"fmt"
	"os"
	//"bufio"
	"errors"
	"io"
	"runtime/debug"
	"strings"
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

	if err != nil {
		panic(err)
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

	//reader := bufio.NewReader(file)

}
