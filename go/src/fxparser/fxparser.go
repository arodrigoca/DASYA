package fxparser

import (
	"errors"
	"fmt"
	"fxlex"
	"os"
	"strings"
)

type Parser struct {
	l         *fxlex.Lexer
	depth     int
	DebugDesc bool
}

func NewParser(l *fxlex.Lexer) *Parser {
	return &Parser{l, 0, true}
}

func (p *Parser) pushTrace(tag string) {

	if p.DebugDesc {
		tabs := strings.Repeat("\t", p.depth)
		fmt.Fprintf(os.Stderr, "%s%s\n", tabs, tag)
	}
	p.depth++
}

func (p *Parser) popTrace() {
	p.depth--
}

func (p *Parser) match(tT fxlex.TokType) (t fxlex.Token, e error, isMatch bool) {

	t, err := p.l.Peek()
	//fmt.Print("Peek says: ")
	//t.PrintToken()
	if err != nil {
		return fxlex.Token{}, err, false
	}
	if t.Type != tT {
		return t, nil, false
	}
	t, err = p.l.Lex()
	//fmt.Print("LEXED: ")
	return t, nil, true

}

func (p *Parser) Body() error {
	p.pushTrace("BODY")
	defer p.popTrace()
	_, _ = p.l.Lex()
	return nil

}

func (p *Parser) Fdecargs() error {
  //<FDECARGS> ::= ',' id id <FDECARGS> |
  //               id id <FDECARGS> |
  //               <EMPTY>

	p.pushTrace("FDECARGS")
	defer p.popTrace()

  _, err, isComma := p.match(fxlex.TokType(','))

  if err != nil {
		return err
	}

	if isComma {
    //Es la primera regla
    //comprobar todos los componentes de la primera regla
    //comprobar el primer ID
    _, err, isId := p.match(fxlex.TokId)
    if err != nil || !isId {
  		err = errors.New("Missing Id on function arguments")
  		return err
  	}
    //comprobar el segundo id
    _, err, isId = p.match(fxlex.TokId)
    if err != nil || !isId {
  		err = errors.New("Missing Id on function arguments")
  		return err
  	}
    return p.Fdecargs()
	}

  //comprobar si es la segunda regla
  _, err, isId := p.match(fxlex.TokId)
  if err != nil {
		return err
	}

	if isId {
    //Es la segunda regla
    //Comprobar todos los componentes de la segunda regla
    //comprobar el segundo id
    _, err, isId = p.match(fxlex.TokId)
    if err != nil || !isId {
  		err = errors.New("Missing Id on function arguments")
  		return err
  	}
    return p.Fdecargs()
	}

  //es la tercera regla, con lo cual empty
	return nil

}

func (p *Parser) Finside() error {
	//<FINSIDE> :: = <FDECARGS> ')' |')'

	p.pushTrace("FINSIDE")
	defer p.popTrace()

  _, err, isRpar := p.match(fxlex.TokType(')'))

  if err != nil {
		return err
	}

	if isRpar {
		return nil
	}

  if err := p.Fdecargs(); err != nil {
		return err
	}

  fmt.Println(p.l.Peek())
  _, err, isRpar = p.match(fxlex.TokType(')'))

  if err != nil || !isRpar {
		err = errors.New("Missing ')' token on function definition")
		return err
	}

	return nil

}

func (p *Parser) Fsig() error {
	//<FSIG> :: = 'func' ID '(' <FINSIDE>
	p.pushTrace("FSIG")
	defer p.popTrace()
	_, err, isFunc := p.match(fxlex.TokFunc)

	if err != nil || !isFunc {
		err = errors.New("Missing 'func' token on function definition")
		return err
	}

	_, err, isId := p.match(fxlex.TokId)

	if err != nil || !isId {
		err = errors.New("Missing function id on function definition")
		return err
	}

	_, err, isLpar := p.match(fxlex.TokType('('))

	if err != nil || !isLpar {
		err = errors.New("Missing '(' token on function definition")
		return err
	}

	if err := p.Finside(); err != nil {
		return err
	}

	return nil
}

func (p *Parser) Func() error {
	//<FUNC> ::= <FSIG> '{' <BODY> '}'
	p.pushTrace("FUNC")
	defer p.popTrace()

	if err := p.Fsig(); err != nil {
		return err
	}

	_, err, isLbra := p.match(fxlex.TokType('{'))

	if err != nil || !isLbra {
		err = errors.New("Missing '{' token on function")
		return err
	}

	if err := p.Body(); err != nil {
		return err
	}

	_, err, isRbra := p.match(fxlex.TokType('}'))

	if err != nil || !isRbra {
		err = errors.New("Missing '}' token on function")
		return err
	}

	return nil

}

func (p *Parser) End() error {
	//<END> ::= <PROG> | <EOF>
	p.pushTrace("END")
	defer p.popTrace()
	_, err, isEOF := p.match(fxlex.TokEof)

	if err != nil {
		return err
	}

	if isEOF {
		p.pushTrace("EOF")
		return nil
	}

	return p.Prog()
}

func (p *Parser) Prog() error {
	//<PROG> ::= <FUNC> <END> | <EOF>
	p.pushTrace("PROG")
	defer p.popTrace()
	_, err, isEOF := p.match(fxlex.TokEof)

	if err != nil {
		return err
	}

	if isEOF {
		p.pushTrace("EOF")
		return nil
	}

	if err := p.Func(); err != nil {
		return err
	} else {
		return p.End()
	}

}

func (p *Parser) Parse() error {
	p.pushTrace("Parse")
	defer p.popTrace()
	if err := p.Prog(); err != nil {
		return err
	}

	return nil
}
