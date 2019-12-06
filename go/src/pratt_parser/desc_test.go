package pratt_test

import (
	"comp/lex"
	"comp/pratt"
	"errors"
	"fmt"
	"math"
	"os"
	"testing"
)

type progExamp struct {
	input string
	isBad bool
	val   float64
}

var genProgs = []progExamp{
	{"1.0 *2.0+3.0", false, 5.0},
	{"1.0 +2.0*3.0", false, 7.0},
	{"3.0 *4.7+5.2", false, 19.3},
	{"3.0 *(4.7+5.2)", false, 29.7},
	//(2^2)^2 = 8 and 2^(2^2) = 16, correct, right associative
	{"2.0 ^ 2.0 ^ 2.0 ", false, 16.0},
	{"2.0 ^ 2.0 ^ 2.0 ^ 2.0", false, 65536.0},
	{"-(2.0)", false, -2.0},
	{"--(2.0)", false, 2.0},
	{"3.0 *+5.2", false, 15.60},
	//bad expr
	{"", true, -1.0},
	{"3.0 **5.2", true, -1.0},
	{"3.0 *", true, -1.0},
	{"* 3.0", true, -1.0},
	{"3.0 *4.7 5.2", true, -1.0},
	{"3.0 * 4.7+5.2)", true, -1.0},
	{"3.0 * (4.7+5.2", true, -1.0},
	{"3.0 * (4.7+5.2", true, -1.0},
	{"3.0 * 4.7+)5.2", true, -1.0},
	{"2.0 ^ *2.0 ^ (2.0 ^ 2.0", true, -1.0},
	{"*", true, -1.0},
	{"()", true, -1.0},
	{"(", true, -1.0},
	{"-", true, -1.0},
}

const Eps = 1e-9

func almostEqual(f, g float64) bool {
	return math.Abs(f-g) <= Eps
}

func TestGen(t *testing.T) {
	var expr *pratt.Expr
	for _, v := range genProgs {
		if testing.Verbose() {
			fmt.Fprintf(os.Stderr, "--> %s\n", v.input)
		}
		l, err := lex.NewFakeLexer(v.input)
		if err != nil {
			t.Fatal(err)
		}
		p := pratt.NewParser(l)
		if err, expr = p.Parse(); err != nil && !v.isBad {
			errs := fmt.Sprintf("%s: %s", err, v.input)
			t.Fatal(errs)
		}

		val := expr.Eval()
		if v.isBad && err == nil {
			errs := fmt.Sprintf("%s should fail evals to %f", v.input, val)
			t.Fatal(errors.New(errs))
		} else if !v.isBad && !almostEqual(val, v.val) {
			errs := fmt.Sprintf("%s  is %f should be %f", v.input, val, v.val)
			t.Fatal(errs)
		}
	}
}
