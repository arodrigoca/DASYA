package expr_parser

import (
	"errors"
	"fmt"
	"fxlex"
	"os"
	"strings"
)

var precTab = map[rune]int{
	')': 1,
	'+': 20,
	'-': 20,
	'*': 30,
	'^': 40,
	'(': 50,
}

var leftTab = map[rune]bool{
	'^': true,
}
var unaryTab = map[rune]bool{
	'+': true,
	'-': true,
	'(': true,
}

type Parser struct {
	l           *fxlex.Lexer
	depth       int
	DebugDesc   bool
	ErrorNumber int
	Errors      []error
}

type Expr struct {
	tok    lex.Token
	ERight *Expr
	ELeft  *Expr
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

///////////////////////////////////////////////////////////////////pratt Parser

//no left context, null-denotation: nud
func (p *Parser) Nud(tok lex.Token) (expr *Expr, err error) {
	var rExpr *Expr
	var rbp int
	p.dPrintf("Nud:  %d, %s \n", rbp, tok)
	if tok.Type == lex.TokLPar { //special unary, parenthesis
		expr, err = p.Expr(rbp)
		if err != nil {
			return nil, err
		}
		if _, err, isClosed := p.match(lex.TokRPar); err != nil {
			return nil, err
		} else if !isClosed {
			return nil, errors.New("unmatched parenthesis")
		}
		return expr, nil
	}
	expr = NewExpr(tok)
	rbp = bindPow(tok)
	rTok := rune(tok.Type)
	if rbp != defRbp { //regular unary operators
		if !unaryTab[rTok] {
			errs := fmt.Sprintf("%s  is not unary", tok.Type)
			return nil, errors.New(errs)
		}
		rExpr, err = p.Expr(rbp)
		if rExpr == nil {
			return nil, errors.New("unary operator without operand")
		}
		expr.ERight = rExpr
	}
	return expr, nil
}

//left context, left-denotation: led
func (p *Parser) Led(left *Expr, tok lex.Token) (expr *Expr, err error) {
	var rbp int
	expr = NewExpr(tok)
	expr.ELeft = left
	rbp = bindPow(tok)
	if isleft := leftTab[rune(tok.Type)]; isleft {
		rbp -= 1
	}
	p.dPrintf("Led: %d, {{%s}} %s \n", rbp, left, tok)
	rExpr, err := p.Expr(rbp)
	if err != nil {
		return nil, err
	}
	if rExpr == nil {
		errs := fmt.Sprintf("missing operand for %s\n", tok.Type)
		return nil, errors.New(errs)
	}
	expr.ERight = rExpr
	return expr, nil
}

const defRbp = 0

func bindPow(tok lex.Token) int {
	if rbp, ok := precTab[rune(tok.Type)]; ok {
		return rbp
	}
	return defRbp
}


func NewExpr(tok lex.Token) (expr *Expr) {
	return &Expr{tok: tok}
}

////////////////////////////////////////////////////////////////////////////////

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
	p.pushTrace("STMNT")
	defer p.popTrace()

	_, err, isId := p.match(fxlex.TokId)
	if err != nil {
		return err
	}

	if isId {
		//es la primera regla
		//return p.Funcall()
		err = p.Funcall()
		if err != nil{
			err = p.ConsumeUntilMarker(";", true)
			if err != nil{
				return err
			}

		}

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
	defer p.popTrace()

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
