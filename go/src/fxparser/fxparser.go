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

func (p *Parser) ErrExpected(place string, found fxlex.Token, wanted string) error {

	err := fmt.Errorf("%s:%d: Expected %s in %s, found %s", found.File, found.Line, wanted, place, found.Lexema)
	fmt.Println(err)

	if p.ErrorNumber >= 5 {
		panic("Too many syntax errors")
		//return nil
		//os.Exit(1)
	}

	p.ErrorNumber += 1
	p.Errors = append(p.Errors, err)

	return err
}

func (p *Parser) ErrGeneric(message string, file string, line int, place string) error {

	err := fmt.Errorf("%s:%d: %s %s", file, line, message, place)

	fmt.Println(err)

	if p.ErrorNumber >= 5 {
		panic("Too many syntax errors")
		//return errors.New("Too many syntax errors")
		//return nil
		//os.Exit(1)
	}

	p.ErrorNumber += 1
	p.Errors = append(p.Errors, err)

	return err
}



func (p *Parser) ConsumeUntilMarker(markers string, consume bool) error {

	for t, _ := p.l.Peek(); ; t, _ = p.l.Peek() {
		//t.PrintToken()
		if t.Type != fxlex.TokEof {
			if strings.Contains(markers, t.Lexema) {
				if consume{
					_, _ = p.l.Lex()
				}
				return nil
			} else {
				_, err := p.l.Lex()
				if err != nil {
					return err
				}
			}
		} else {
			panic("Found EOF")
			//os.Exit(1)
		}
	}

	return nil
}

func (p *Parser) ConsumeUntilToken(token_type fxlex.TokType) error {

	for t, _ := p.l.Peek(); ; t, _ = p.l.Peek() {

		if t.Type != fxlex.TokEof {
			if t.Type == token_type {
				//_, _ = p.l.Lex()
				return nil
			} else {
				_, err := p.l.Lex()
				if err != nil {
					return err
				}
			}
		} else {
			panic("Found EOF")
			//os.Exit(1)
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
		//fmt.Println("CONSUMED UNTIL MARKER")
		p.ConsumeUntilMarker(",", false)
		//return nil
		//return err
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
	tok, err, isRpar := p.match(fxlex.TokType(')'))
	//fmt.Println(tok)
	if err != nil {
		err = p.ErrExpected("on function call", tok, ")")
		return err
	}

	if isRpar {
		//es la segunda regla
		tok, err, isSemic := p.match(fxlex.TokType(';'))
		//fmt.Println(tok)
		if err != nil || !isSemic {
			//err = errors.New("Missing ';' token on function call")
			err = p.ErrExpected("on function call", tok, ";")
			return err
		}
		return nil
	}

	err = p.Fargs()
	if err != nil {
		return err
	}

	tok, err, isRpar = p.match(fxlex.TokType(')'))
	//fmt.Println(tok)
	if err != nil || !isRpar {
		//err = errors.New("Missing ')' token on function call")
		//err := p.ErrGeneric("Missing ')' token", tok.File, tok.Line, "on function call")
		//return err
		err = p.ErrExpected("on function call", tok, ")")
		//p.ConsumeUntilMarker(";")
		return err
	}

	tok, err, isSemic := p.match(fxlex.TokType(';'))
	if err != nil || !isSemic {
		//err = errors.New("Missing ';' token on function call")
		err = p.ErrExpected("on function call", tok, ";")
		//err = p.ErrGeneric("Missing ';' token", tok.File, tok.Line, "on function call")
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
	 	err = p.ErrExpected("function call", tok_1, "(")
		//p.ConsumeUntilMarker(")")
		return err

	}

	err = p.Rfuncall()
	if err != nil{

		//p.ConsumeUntilMarker(";")
		return err
	}

	return nil
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
	//err = errors.New("Bad atom")
	err = p.ErrGeneric("Bad atom", t.File, t.Line, "")

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

	has_error := false

	tok_1, err, isIter := p.match(fxlex.TokIter)
	if err != nil || !isIter {
		//err = errors.New("Missing 'iter' on iter definition")
		//return err
		//err= p.ErrExpected("iter declaration", tok_1, "Iter")
		err = p.ErrGeneric("Bad function call or empty", tok_1.File, tok_1.Line, "on function body")
		return err
	}

	tok_2, err, isLpar := p.match(fxlex.TokType('('))
	if err != nil || !isLpar {
		//err = errors.New("Missing '(' token on iter definition")
		//return err
		err= p.ErrExpected("iter declaration", tok_2, "(")
		p.ConsumeUntilMarker(":=", false)
		has_error = true

	}

	if !has_error{
		has_error = false
		tok, err, isId := p.match(fxlex.TokId)
		if err != nil || !isId {
			//err = errors.New("Missing id on iter definition")
			err = p.ErrGeneric("Missing id", tok.File, tok.Line, "on iter definiton")
			p.ConsumeUntilMarker(":=", false)
			has_error = true
		}
	}

	has_error = false
	tok, err, isDDEq := p.match(fxlex.TokDDEq)
	if err != nil || !isDDEq {
		//err = errors.New("Missing ':=' token on iter definition")
		err = p.ErrGeneric("Missing ':=' token", tok.File, tok.Line, "on iter definition")
		p.ConsumeUntilMarker(";", false)
		has_error = true
	}

	if !has_error{
		has_error = false
		err = p.Expr()
		if err != nil {
			p.ConsumeUntilMarker(";", false)
			has_error = true
		}
	}

	//_, err, isSemic := p.match(fxlex.TokType(';'))
	has_error = false
	tok, err, isSemic := p.match(fxlex.TokType(';'))
	if err != nil || !isSemic {

		//err = errors.New("Missing ';' token on iter definition")
		//return err
		err = p.ErrGeneric("Missing ';' token", tok.File, tok.Line, "on iter definiton")
		p.ConsumeUntilMarker(")", false)
		has_error = true
	}

	if !has_error{
		has_error = false
		err = p.Expr()
		if err != nil {
			//p.ConsumeUntilMarker("{}();")
			return nil
		}

		tok, err, isComma := p.match(fxlex.TokType(','))
		if err != nil || !isComma {
			//err = errors.New("Missing ',' token on iter definition")
			err = p.ErrGeneric("Missing ',' token", tok.File, tok.Line, "on iter definition")
			p.ConsumeUntilMarker(")", false)
		}

		if !has_error{
			err = p.Expr()
			if err != nil {
				p.ConsumeUntilMarker(")", false)
			}
		}
	}

	has_error = false
	tok, err, isRpar := p.match(fxlex.TokType(')'))
	if err != nil || !isRpar {
		//err = errors.New("Missing ')' token on iter definition")
		err = p.ErrGeneric("Missing ')' token", tok.File, tok.Line, "on iter definition")
		p.ConsumeUntilMarker("{", false)
		has_error = true
	}

	has_error = false
	tok, err, isLbra := p.match(fxlex.TokType('{'))

	if err != nil || !isLbra {
		//err = errors.New("Missing '{' token on iter definition")
		err = p.ErrGeneric("Missing '{' token", tok.File, tok.Line, "on iter definition")
	}

	err = p.Body()
	if err != nil {
		p.ConsumeUntilMarker("}", false)
		has_error = true
	}

	tok, err, isRbra := p.match(fxlex.TokType('}'))

	if err != nil || !isRbra {
		//err = errors.New("Missing '}' token on iter definition")
		err = p.ErrGeneric("Missing '}' token", tok.File, tok.Line, "on iter definition")
		return err
	}

	return nil
}

func (p *Parser) Stmnt() error {
	//<STMNT> ::= id <FUNCALL> |
	//            <ITER>
	//UPDATE P5: AÑADIR DECLARACIONES
	//<STMNT> ::= id <FUNCALL>       | done
	//						id "=" <EXPR> ";"	 |
	//						int id ";"         |
	//						bool id ";"        | done
	//            <ITER>

	p.pushTrace("STMNT")
	defer p.popTrace()

	_, err, isId := p.match(fxlex.TokId)
	if err != nil {
		return err
	}

	next_token, _ := p.l.Peek()
	if isId {
		//es la primera regla o la segunda
		//return p.Funcall()

		if next_token.Type == fxlex.TokType('('){
			err = p.Funcall()
			if err != nil{
				err = p.ConsumeUntilMarker(";", true)
				if err != nil{
					return err
				}

			}
		}else if next_token.Type == fxlex.TokType('='){
			//es la segunda regla
			fmt.Println("Second rule")
			_, err, is_eq := p.match(fxlex.TokType('='))
			if err != nil{
				return err
			}

			if !is_eq{
				panic("Something went wrong while parsing")
			}

			err = p.Expr()

			if err != nil{
				p.ConsumeUntilMarker(";", true)
				return nil
			}

			tok_semic, err, is_semic := p.match(fxlex.TokType(';'))

			if err != nil{
				return err
			}

			if !is_semic{
				err = p.ErrExpected("Asignation", tok_semic, ";")
				return nil
			}

			return nil

		}else{
			err = p.ErrGeneric("Malformed asignation", next_token.File, next_token.Line, "")
			p.ConsumeUntilMarker(";", true)
			return nil
		}

		return nil

	}else if next_token.Type == fxlex.TokDefInt{
		//es la tercera regla
		fmt.Println("Third rule")
		//stub
		p.ConsumeUntilMarker(";", true)
		return nil

	}else if next_token.Type == fxlex.TokDefBool{
		//es la cuarta regla
		fmt.Println("Fourth rule")
		//stub
		p.ConsumeUntilMarker(";", true)
		return nil
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
			err= p.ErrExpected("function arguments", tok_2, "Type")
			p.ConsumeUntilMarker(")", false)
			return nil
		}
		//comprobar el segundo id
		tok_1, err, isId := p.match(fxlex.TokId)
		if err != nil || !isId {
			//err = errors.New("Missing Id on function arguments")
			//return err
			err= p.ErrExpected("function declaration", tok_1, "Id")
			p.ConsumeUntilMarker(")", false)
			return nil
		}

		return p.Fdecargs()
	}

	//comprobar si es la segunda regla
	_, err, isInt := p.match(fxlex.TokDefInt)
	_, err1, isBool := p.match(fxlex.TokDefBool)

	if err != nil || err1 != nil {

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
			err= p.ErrExpected("function declaration", tok_1, "Id")
			p.ConsumeUntilMarker(")", false)
			return nil
		}

		return p.Fdecargs()
	}
	//es la tercera regla, con lo cual empty o bien hay algún error
	t, err := p.l.Peek()

	if t.Type == fxlex.TokType(')'){
		return nil
	}

	tok, err, _ := p.match(fxlex.TokType(')'))
	err= p.ErrExpected("function definition", tok, ")")
	p.ConsumeUntilMarker(")", false)
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
		err= p.ErrExpected("function declaration", tok_1, ")")
		//err = errors.New("Missing ')' token on function definition")
		//return err
		return err
	}

	return nil

}

func (p *Parser) Fsig() error {
	//<FSIG> :: = 'func' ID '(' <FINSIDE> |
	//						'func' main '(' <FINSIDE> |

	p.pushTrace("FSIG")
	defer p.popTrace()
	tok_1, err, isFunc := p.match(fxlex.TokFunc)

	if err != nil || !isFunc {
		err= p.ErrExpected("function declaration", tok_1, "func")
		//return err
		return err
	}

	_, err, isId := p.match(fxlex.TokId)
	_, err_main, isMain := p.match(fxlex.TokMain)


	if err != nil || err_main != nil{

		if err != nil{
			//err= p.ErrExpected("function declaration", tok_2, "id")
			//return err
			return err
		}else if err_main != nil{
			//err= p.ErrExpected("function declaration", tok_main, "id")
			//return err
			return err
		}
	}

	if isId || isMain{

		tok_3, err, isLpar := p.match(fxlex.TokType('('))

		if err != nil || !isLpar {
			err= p.ErrExpected("function declaration", tok_3, "(")
			//return err
			return err
		}

		if err := p.Finside(); err != nil {
			//return err
			p.ConsumeUntilMarker("{", false)
			return nil
		}

	}else{
		//even though the function is not correctly defined, keep the flow and evaluate everything to search for more
		//errors
		tok_err, _, _ := p.match(fxlex.TokType('('))
		err= p.ErrExpected("function declaration", tok_err, "main or function id")
		p.ConsumeUntilMarker("(", false)

		tok_3, err, isLpar := p.match(fxlex.TokType('('))

		if err != nil || !isLpar {
			err= p.ErrExpected("function declaration", tok_3, "(")
			//return err
			return err
		}

		if err := p.Finside(); err != nil {
			p.ConsumeUntilMarker("{", false)
			return nil
		}

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
		err= p.ErrExpected("function", tok_1, "{")
		//p.ConsumeUntilMarker("}")
		//return err
		return err
	}

	if err := p.Body(); err != nil {
		return err
	}

	_, err, isRbra := p.match(fxlex.TokType('}'))
	if err != nil || !isRbra {
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

func (p *Parser) Parse() []error {
	p.pushTrace("Parse")
	//defer p.popTrace()
	defer func(){
		p.popTrace()
		if r := recover(); r != nil{
			fmt.Println(p.Errors)
		}
	}()

	p.Prog()


	if p.Errors != nil{
		fmt.Println("SYNTAX ERROR")
		return p.Errors
	}

	/*
	if err := p.Prog(); err != nil {
		fmt.Println("SYNTAX ERROR")
		return err
	}
	*/

	return nil
}
