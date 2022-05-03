package calc

import (
	"fmt"
	. "parserc-go/parserc"
	"strconv"
)

func toString(v any) any {
	switch v.(type) {
	case rune:
		return fmt.Sprintf("%c", v)
	default:
		return fmt.Sprintf("%s", v)
	}
}

func join(list any) any {
	str := ""
	for _, e := range list.([]any) {
		str += toString(e).(string)
	}
	return str
}

func toFloat(a any) any {
	v, _ := strconv.ParseFloat(a.(string), 64)
	return v
}

func calc(p any) any {
	v := p.(Pair).First.(float64)
	for _, e := range p.(Pair).Second.([]any) {
		switch e.(Pair).First.(string) {
		case "+":
			v += e.(Pair).Second.(float64)
		case "-":
			v -= e.(Pair).Second.(float64)
		case "*":
			v *= e.(Pair).Second.(float64)
		case "/":
			v /= e.(Pair).Second.(float64)
		}
	}
	return v
}

var (
	w           = Chs(' ', '\t', '\n', '\r')
	ws          = w.Many()
	digit       = Range('0', '9').Map(toString)
	digits      = digit.Many1().Map(join)
	integer     = digits.Map(toFloat).SurroundedBy(ws)
	decimal     = Seq(digits, Ch('.'), digits).Map(join).Map(toFloat).SurroundedBy(ws)
	add         = Str("+").SurroundedBy(ws)
	sub         = Str("-").SurroundedBy(ws)
	mul         = Str("*").SurroundedBy(ws)
	div         = Str("/").SurroundedBy(ws)
	lp          = Str("(").SurroundedBy(ws)
	rp          = Str(")").SurroundedBy(ws)
	expr        = NewParser()
	bracketExpr = Skip(lp).And(expr).Skip(rp)
	fact        = OneOf(decimal, integer, bracketExpr)
	term        = fact.And(mul.Or(div).And(fact).Many()).Map(calc)
)

func init() {
	expr.Set(term.And(add.Or(sub).And(term).Many()).Map(calc))
}

func Eval(s string) float64 {
	r, err := expr.ParseToEnd(s)
	if err != nil {
		panic(err)
	}
	return r.(float64)
}
