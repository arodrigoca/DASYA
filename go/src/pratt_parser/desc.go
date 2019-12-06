package pratt

import (
	"comp/lex"
	"errors"
	"fmt"
	"os"
	"strings"
)

/*
 * 	An interpretation of Pratt parsing
 *	https://web.archive.org/web/20151223215421/http://hall.org.ua/halls/wizzard/pdf/Vaughan.Pratt.TDOP.pdf
 *	Proceedings of the 1st Annual ACM SIGACT-SIGPLAN
 *	Symposium on Principles of Programming Languages (1973)
 */

type Parser struct {
	l     *lex.Lexer
	depth int
	tag   string
}

func NewParser(l *lex.Lexer) *Parser {
	return &Parser{l, 0, ""}
}

//https://web.archive.org/web/20151223215421/http://hall.org.ua/halls/wizzard/pdf/Vaughan.Pratt.TDOP.pdf

func (p *Parser) match(tT lex.TokType) (tok lex.Token, err error, isMatch bool) {
	p.dPrintf("match: %s", tT)
	tok, err = p.l.Peek()
	if err != nil {
		return lex.Token{}, err, false
	}
	if tok.Type != tT {
		return tok, nil, false
	}
	p.l.Lex() //already peeked
	p.pushTrace(tok.String())
	defer p.popTrace(&err)
	return tok, nil, true
}

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

func (p *Parser) Expr(rbp int) (expr *Expr, err error) {
	var left *Expr

	s := fmt.Sprintf("Expr: %d", rbp)
	p.pushTrace(s)
	defer p.popTrace(&err)

	tok, err := p.l.Peek()
	if err != nil {
		return expr, err
	}
	p.dPrintf("expr: Nud Lex: %s", tok)

	if tok.Type == lex.TokEof {
		return expr, nil
	}
	p.l.Lex() //already peeked
	if left, err = p.Nud(tok); err != nil {
		return nil, err
	}
	expr = left
	for {
		tok, err := p.l.Peek()
		if err != nil {
			return expr, err
		}
		if tok.Type == lex.TokEof || tok.Type == lex.TokRPar {
			return expr, nil
		}
		if bindPow(tok) <= rbp {
			p.dPrintf("Not enough binding: %d <= %d, %s\n", bindPow(tok), rbp, tok)
			return left, nil
		}
		p.l.Lex() //already peeked
		p.dPrintf("expr: led Lex: %s", tok)
		if left, err = p.Led(left, tok); err != nil {
			return expr, err
		}
		expr = left
	}
	return expr, err
}

//PROG :== EXPR EOF
func (p *Parser) Prog() (e error, expr *Expr) {
	p.pushTrace("Prog")
	defer p.popTrace(&e)
	expr, err := p.Expr(defRbp - 1)
	if err != nil {
		return err, nil
	}
	if expr == nil {
		return errors.New("empty expression"), nil
	}
	t, err, isEof := p.match(lex.TokEof)
	if err != nil {
		return err, nil
	}
	if !isEof {
		es := fmt.Sprintf("need %s, got %s", lex.TokEof, t.Type)
		return errors.New(es), nil
	}
	return nil, expr
}

func (p *Parser) Parse() (err error, expr *Expr) {
	p.pushTrace("Parse")
	defer p.popTrace(&err)
	if err, expr = p.Prog(); err != nil {
		return err, nil
	}

	return nil, expr
}

const DebugDesc = true

func (p *Parser) pushTrace(tag string) {
	if DebugDesc {
		tabs := strings.Repeat("\t", p.depth)
		fmt.Fprintf(os.Stderr, "->%s%s\n", tabs, tag)
	}
	p.tag = tag
	p.depth++
}

func (p *Parser) dPrintf(format string, a ...interface{}) {
	if DebugDesc {
		tabs := strings.Repeat("\t", p.depth)
		format = fmt.Sprintf("%s%s", tabs, format)
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

func (p *Parser) popTrace(e *error) {
	if e != nil && *e != nil {
		if DebugDesc {
			tabs := strings.Repeat("\t", p.depth)
			fmt.Fprintf(os.Stderr, "<-%s%s:%s\n", tabs, p.tag, *e)
		}
	}
	p.tag = ""
	p.depth--
}
