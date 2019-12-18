package fxparser

import (
	"errors"
	"fmt"
	"fxlex"
	"os"
	"strings"
)

type Parser struct {
	l           *fxlex.Lexer
	depth       int
	DebugDesc   bool
	ErrorNumber int
	Errors      []error
}

func NewParser(l *fxlex.Lexer) *Parser {

	var erarray []error
	return &Parser{l, 0, true, 0, erarray}
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
	if err != nil {
		return fxlex.Token{}, err, false
	}
	if t.Type != tT {
		return t, nil, false
	}
	t, err = p.l.Lex()
	return t, nil, true

}

func (p *Parser) ErrExpected(place string, found fxlex.Token, wanted string) {

	err := fmt.Errorf("%s:%d: Expected %s in %s, found %s", found.File, found.Line, wanted, place, found.Lexema)
	fmt.Println(err)
	p.ErrorNumber += 1
	p.Errors = append(p.Errors, err)
	if p.ErrorNumber > 5 {
		panic("Too many syntax errors")
	}

}

func (p *Parser) ConsumeUntilMarker(markers string) error {

	for t, _ := p.l.Peek(); ; t, _ = p.l.Peek() {
		//t.PrintToken()
		if t.Type != fxlex.TokEof {
			if strings.Contains(markers, t.Lexema) {
				_, _ = p.l.Lex()
				return nil
			} else {
				_, err := p.l.Lex()
				if err != nil {
					return err
				}
			}
		} else {
			panic("Found EOF")
		}
	}

	return nil
}

func (p *Parser) Exprend() error {
	//<EXPREND> ::= ',' <FARGS> | <EMPTY>

	p.pushTrace("EXPREND")
	defer p.popTrace()
	t, err := p.l.Peek()
	if err != nil {
		return err
	}
	if t.Type == fxlex.TokType(',') {
		//Es la primera regla
		t, err = p.l.Lex()
		if err != nil {
			return err
		}

		return p.Fargs()
	}

	//es la segunda regla
	return nil
}

func (p *Parser) Fargs() error {
	//<FARGS> ::= <EXPR> <EXPREND>
	p.pushTrace("FARGS")
	defer p.popTrace()
	if err := p.Expr(); err != nil {
		return err
	}
	if err := p.Exprend(); err != nil {
		return err
	}

	return nil
}

func (p *Parser) Rfuncall() error {
	//<RFUNCALL> := <FARGS> ')' ';' | ')' ';'

	p.pushTrace("RFUNCALL")
	defer p.popTrace()
	_, err, isRpar := p.match(fxlex.TokType(')'))
	if err != nil {
		return err
	}

	if isRpar {
		//es la segunda regla
		_, err, isSemic := p.match(fxlex.TokType(';'))
		if err != nil || !isSemic {
			err = errors.New("Missing ';' token on function call")
			return err
		}
		return nil
	}

	err = p.Fargs()
	if err != nil {
		return err
	}

	_, err, isRpar = p.match(fxlex.TokType(')'))
	if err != nil || !isRpar {
		err = errors.New("Missing ')' token on function call")
		return err
	}

	_, err, isSemic := p.match(fxlex.TokType(';'))
	if err != nil || !isSemic {
		err = errors.New("Missing ';' token on function call")
		return err
	}

	return nil
}

func (p *Parser) Funcall() error {
	//<FUNCALL> ::= '(' <RFUNCALL>

	p.pushTrace("FUNCALL")
	defer p.popTrace()
	tok_1, err, isLpar := p.match(fxlex.TokType('('))
	if err != nil || !isLpar {
		//err = errors.New("Missing '(' on function call")
		//return err
		p.ErrExpected("function call", tok_1, "(")
		p.ConsumeUntilMarker(";")
		return nil
	}

	return p.Rfuncall()
}

func (p *Parser) Atom() error {
	//<ATOM> ::= id | intval | boolVal
	p.pushTrace("ATOM")
	defer p.popTrace()
	t, err := p.l.Peek()
	if err != nil {
		return err
	}
	if ((t.Type == fxlex.TokId) || (t.Type == fxlex.TokValInt) || (t.Type == fxlex.TokValBool)) != false {
		_, err = p.l.Lex()
		if err != nil {
			return err
		}
		return nil
	}
	err = errors.New("Bad atom")
	return err
}

func (p *Parser) Expr() error {
	//TODO
	//<EXPR> :: = <ATOM>
	p.pushTrace("EXPR")
	defer p.popTrace()
	err := p.Atom()
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) Iter() error {
	//<ITER> ::= 'iter' '(' id ':=' <EXPR> ';' <EXPR> ',' <EXPR> ')' '{' <BODY> '}'
	p.pushTrace("ITER")
	defer p.popTrace()

	tok_1, err, isIter := p.match(fxlex.TokIter)
	if err != nil || !isIter {
		//err = errors.New("Missing 'iter' on iter definition")
		//return err
		p.ErrExpected("iter declaration", tok_1, "Iter")
	}

	tok_2, err, isLpar := p.match(fxlex.TokType('('))
	if err != nil || !isLpar {
		//err = errors.New("Missing '(' token on iter definition")
		//return err
		p.ErrExpected("iter declaration", tok_2, "(")
		return nil
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
	if err != nil {
		return err
	}

	_, err, isSemic := p.match(fxlex.TokType(';'))
	if err != nil || !isSemic {
		err = errors.New("Missing ';' token on iter definition")
		return err
	}

	err = p.Expr()
	if err != nil {
		return err
	}

	_, err, isComma := p.match(fxlex.TokType(','))
	if err != nil || !isComma {
		err = errors.New("Missing '(' token on iter definition")
		return err
	}

	err = p.Expr()
	if err != nil {
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
	if err != nil {
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
	if err != nil {
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
	if err != nil {
		return err
	}
	if t.Type == fxlex.TokType('}') {
		//ha acabado el body, por lo tanto empty
		return nil
	} else {
		//hay más sentencias
		return p.Body()
	}

	err = errors.New("Unkown or malformed statement")
	return err

}

func (p *Parser) Body() error {

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
	_, err, isComma := p.match(fxlex.TokType(','))

	if err != nil {
		return err
	}

	if isComma {
		//Es la primera regla
		//comprobar todos los componentes de la primera regla
		//comprobar el primer ID
		_, err, isInt := p.match(fxlex.TokDefInt)
		tok_2, err1, isBool := p.match(fxlex.TokDefBool)
		if err != nil || err1 != nil || (isInt || isBool) == false {
			//err = errors.New("Missing Id on function arguments")
			//return err
			p.ErrExpected("function declaration", tok_2, "Id")
			return nil
		}
		//comprobar el segundo id
		tok_1, err, isId := p.match(fxlex.TokId)
		if err != nil || !isId {
			//err = errors.New("Missing Id on function arguments")
			//return err
			p.ErrExpected("function declaration", tok_1, "Id")
			return nil
		}

		return p.Fdecargs()
	}

	//comprobar si es la segunda regla
	_, err, isInt := p.match(fxlex.TokDefInt)
	_, err1, isBool := p.match(fxlex.TokDefBool)

	if err != nil || err1 != nil || (isInt || isBool) == false {
		return err
	}

	if isInt || isBool {
		//Es la segunda regla
		//Comprobar todos los componentes de la segunda regla
		//comprobar el segundo id
		tok_1, err, isId := p.match(fxlex.TokId)
		if err != nil || !isId {
			//err = errors.New("Missing Id on function arguments")
			//return err
			p.ErrExpected("function declaration", tok_1, "Id")
			return nil
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

	tok_1, err, isRpar := p.match(fxlex.TokType(')'))

	if err != nil || !isRpar {
		p.ErrExpected("function declaration", tok_1, ")")
		//err = errors.New("Missing ')' token on function definition")
		//return err
		return nil
	}

	return nil

}

func (p *Parser) Fsig() error {
	//<FSIG> :: = 'func' ID '(' <FINSIDE>
	p.pushTrace("FSIG")
	defer p.popTrace()
	tok_1, err, isFunc := p.match(fxlex.TokFunc)

	if err != nil || !isFunc {
		p.ErrExpected("function declaration", tok_1, "func")
		//return err
		return nil
	}

	tok_2, err, isId := p.match(fxlex.TokId)

	if err != nil || !isId {
		p.ErrExpected("function declaration", tok_2, "id")
		//return err
		return nil
	}

	tok_3, err, isLpar := p.match(fxlex.TokType('('))

	if err != nil || !isLpar {
		p.ErrExpected("function declaration", tok_3, "(")
		//return err
		return nil
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

	tok_1, err, isLbra := p.match(fxlex.TokType('{'))

	if err != nil || !isLbra {
		p.ErrExpected("function", tok_1, "{")
		p.ConsumeUntilMarker("}")
		//return err
		return nil
	}

	if err := p.Body(); err != nil {
		return err
	}

	tok_2, err, isRbra := p.match(fxlex.TokType('}'))
	if err != nil || !isRbra {
		p.ErrExpected("function", tok_2, "}")
		//return err
		return nil
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
		fmt.Println("SYNTAX ERROR")
		return err
	}

	return nil
}
