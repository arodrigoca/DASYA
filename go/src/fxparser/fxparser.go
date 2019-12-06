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

func (p *Parser) Rfuncall() error {
//TODO
//<RFUNCALL> := <FARGS> ')' ';' | ')' ';'

  p.pushTrace("RFUNCALL")
	defer p.popTrace()
	_, err, isRpar := p.match(fxlex.TokType(')'))
	if err != nil{
		return err
	}

	if isRpar{
		//es la segunda regla
		_, err, isRpar := p.match(fxlex.TokType(';'))
		if err != nil || !isRpar {
			err = errors.New("Missing ';' token on function call")
			return err
		}
		return nil
	}
	//es la primera regla
  return nil
}

func (p *Parser) Funcall() error {
//TODO
//<FUNCALL> ::= '(' <RFUNCALL>

  p.pushTrace("FUNCALL")
	defer p.popTrace()
	_, err, isLpar := p.match(fxlex.TokType('('))
	if err != nil || !isLpar {
		err = errors.New("Missing '(' on function call")
		return err
	}

  return p.Rfuncall()
}

func (p *Parser) Expr() error {

	_, _ = p.l.Lex()
	return nil
}

func (p *Parser) Iter() error {
//TODO
//<ITER> ::= 'iter' '(' id ':=' <EXPR> ';' <EXPR> ',' <EXPR> ')' '{' <BODY> '}'
  p.pushTrace("ITER")
	defer p.popTrace()

	_, err, isIter := p.match(fxlex.TokIter)
	if err != nil || !isIter {
		err = errors.New("Missing 'iter' on iter definition")
		return err
	}

	_, err, isLpar := p.match(fxlex.TokType('('))
	if err != nil || !isLpar {
		err = errors.New("Missing '(' token on iter definition")
		return err
	}

	_, err, isId := p.match(fxlex.TokId)
	if err != nil || !isId {
		err = errors.New("Missing id on iter definition")
		return err
	}

	_, err, isDDEq := p.match(fxlex.TokDDEq)
	if err != nil || !isDDEq {
		err = errors.New("Missing ':=' token on iter definition")
		return err
	}

	err = p.Expr()
	if err != nil{
		return err
	}

	_, err, isSemic := p.match(fxlex.TokType(';'))
	if err != nil || !isSemic {
		err = errors.New("Missing ';' token on iter definition")
		return err
	}

	err = p.Expr()
	if err != nil{
		return err
	}

	_, err, isComma := p.match(fxlex.TokType(','))
	if err != nil || !isComma {
		err = errors.New("Missing '(' token on iter definition")
		return err
	}

	err = p.Expr()
	if err != nil{
		return err
	}

	_, err, isRpar := p.match(fxlex.TokType(')'))
	if err != nil || !isRpar {
		err = errors.New("Missing ')' token on iter definition")
		return err
	}

	_, err, isLbra := p.match(fxlex.TokType('{'))
	if err != nil || !isLbra {
		err = errors.New("Missing '{' token on iter definition")
		return err
	}

	err = p.Body()
	if err != nil{
		return err
	}

	_, err, isRbra := p.match(fxlex.TokType('}'))
	if err != nil || !isRbra {
		err = errors.New("Missing '}' token on iter definition")
		return err
	}

  return nil
}

func (p *Parser) Stmnt() error {
  //<STMNT> ::= id <FUNCALL> |
  //            <ITER>
  p.pushTrace("STMNT")
	defer p.popTrace()

  _, err, isId := p.match(fxlex.TokId)
  if err != nil{
    return err
  }

  if isId {
    //es la primera regla
    return p.Funcall()
  }
  //es la segunda regla
  return p.Iter()
}

func (p *Parser) Stmntend() error {
  //<STMNTEND> ::= <BODY> |
  //               <EMPTY>
  p.pushTrace("STMNTEND")
	defer p.popTrace()

  t, err := p.l.Peek()
  if err != nil{
    return err
  }
  if t.Type == fxlex.TokType('}'){
    //ha acabado el body, por lo tanto empty
    return nil
  }else{
		//hay más sentencias
  	return p.Body()
	}

	err = errors.New("Unkown or malformed statement")
	return err

}

func (p *Parser) Body() error {

  //TODO

  //<BODY> ::= <STMNT> <STMNTEND>

	p.pushTrace("BODY")
	defer p.popTrace()

  if err := p.Stmnt(); err != nil {
		return err
	}

  if err := p.Stmntend(); err != nil {
		return err
	}

  return nil

}

func (p *Parser) Fdecargs() error {
  //<FDECARGS> ::= ',' id id <FDECARGS> |
  //               id id <FDECARGS> |
  //               <EMPTY>

	p.pushTrace("FDECARGS")
	defer p.popTrace()
  //fmt.Println(fxlex.TokDefInt)
  //fmt.Println(p.l.Peek())
  _, err, isComma := p.match(fxlex.TokType(','))

  if err != nil {
		return err
	}

	if isComma {
    //Es la primera regla
    //comprobar todos los componentes de la primera regla
    //comprobar el primer ID
    //fmt.Println("ES LA PRIMERA REGLA")
    _, err, isInt := p.match(fxlex.TokDefInt)
    _, err1, isBool := p.match(fxlex.TokDefBool)
    if err != nil || err1 != nil || (isInt||isBool) == false{
  		err = errors.New("Missing Id on function arguments")
  		return err
  	}
    //comprobar el segundo id
    _, err, isId := p.match(fxlex.TokId)
    if err != nil || !isId {
  		err = errors.New("Missing Id on function arguments")
  		return err
  	}
    return p.Fdecargs()
	}

  //comprobar si es la segunda regla
  _, err, isInt := p.match(fxlex.TokDefInt)
  _, err1, isBool := p.match(fxlex.TokDefBool)

  if err != nil || err1 != nil || (isInt||isBool) == false{
		return err
	}

	if isInt || isBool {
    //Es la segunda regla
    //Comprobar todos los componentes de la segunda regla
    //comprobar el segundo id
    //fmt.Println("ES LA SEGUNDA REGLA")
    _, err, isId := p.match(fxlex.TokId)
    if err != nil || !isId {
  		err = errors.New("Missing Id on function arguments")
  		return err
  	}
    return p.Fdecargs()
	}

  //es la tercera regla, con lo cual empty
  //fmt.Println("ES LA TERCERA REGLA")
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
