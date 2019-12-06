package pratt

import (
	"comp/lex"
	"fmt"
	"math"
	"os"
)

type Expr struct {
	tok    lex.Token
	ERight *Expr
	ELeft  *Expr
}

func NewExpr(tok lex.Token) (expr *Expr) {
	return &Expr{tok: tok}
}

func (e *Expr) String() string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("\t%p EXPR[%s](%f) L->%p R->%p", e, e.tok.Type, e.tok.TokFloatVal, e.ELeft, e.ERight)
}

const DebugExpr = false

func (e *Expr) Eval() float64 {
	if DebugExpr {
		fmt.Fprintf(os.Stderr, "%s\n", e)
	}
	rV := 0.0
	lV := 0.0
	if e == nil {
		return 0
	}
	if e.ERight != nil {
		rV = e.ERight.Eval()
	}
	if e.ELeft != nil {
		lV = e.ELeft.Eval()
	}
	tok := e.tok
	switch tok.Type {
	case lex.TokExp:
		return math.Pow(lV, rV)
	case lex.TokMin:
		return lV - rV
	case lex.TokAdd:
		return lV + rV
	case lex.TokMul:
		return lV * rV
	case lex.TokFloatVal:
		return tok.TokFloatVal
	default:
		panic("Bad subtree")
	}
}
