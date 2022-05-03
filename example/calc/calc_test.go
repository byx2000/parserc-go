package calc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testEvalSuccess(t *testing.T, s string, v float64) {
	assert.InDelta(t, v, Eval(s), 1e-6)
}

func testEvalFailed(t *testing.T, s string) {
	assert.Panics(t, func() {
		Eval(s)
	})
}

func TestEval(t *testing.T) {
	testEvalSuccess(t, "123.456*67.89", 123.456*67.89)
	testEvalSuccess(t, " 0.78 / 10.4 ", 0.78/10.4)
	testEvalSuccess(t, "(2+3)*(7-4)", (2+3)*(7-4))
	testEvalSuccess(t, "2.4 / 5.774 * (6 / 3.57 + 6.37) - 2 * 7 / 5.2 + 5", 2.4/5.774*(6/3.57+6.37)-2*7/5.2+5)
	testEvalSuccess(t, "77.58* ( 6 / 3.14+55.2234 ) -2 * 6.1/ ( 1.0+2/ (4.0-3.8*5))  ", 77.58*(6/3.14+55.2234)-2*6.1/(1.0+2/(4.0-3.8*5)))

	testEvalFailed(t, "")
	testEvalFailed(t, "1+")
	testEvalFailed(t, "1 * (2 + 3")
	testEvalFailed(t, "1 * 2 + 3)")
	testEvalFailed(t, "a + 12")
	testEvalFailed(t, " 1 2  4")
}
