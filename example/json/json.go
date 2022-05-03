package json

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

func toInt(s any) any {
	v, _ := strconv.Atoi(s.(string))
	return v
}

func toFloat(a any) any {
	v, _ := strconv.ParseFloat(a.(string), 64)
	return v
}

func toBool(a any) any {
	r, _ := strconv.ParseBool(a.(string))
	return r
}

func buildObj(ps any) any {
	m := map[string]any{}
	for _, p := range ps.([]any) {
		m[p.(Pair).First.(string)] = p.(Pair).Second
	}
	return m
}

var (
	w        = Chs(' ', '\t', '\n', '\r')
	ws       = w.Many()
	digit    = Range('0', '9').Map(toString)
	digits   = digit.Many1().Map(join)
	integer  = digits.Map(toInt).SurroundedBy(ws)
	decimal  = Seq(digits, Ch('.'), digits).SurroundedBy(ws).Map(join).Map(toFloat)
	str      = Skip(Ch('"')).And(Not('"').Many()).Skip(Ch('"')).Map(join).SurroundedBy(ws)
	boolean  = Str("true").Or(Str("false")).Map(toBool).SurroundedBy(ws)
	objStart = Ch('{').SurroundedBy(ws)
	objEnd   = Ch('}').SurroundedBy(ws)
	arrStart = Ch('[').SurroundedBy(ws)
	arrEnd   = Ch(']').SurroundedBy(ws)
	colon    = Ch(':').SurroundedBy(ws)
	comma    = Ch(',').SurroundedBy(ws)
	jsonObj  = NewParser()
	arr      = Skip(arrStart).And(SeparatedBy(comma, jsonObj).Optional([]any{})).Skip(arrEnd)
	pair     = str.Skip(colon).And(jsonObj)
	obj      = Skip(objStart).And(SeparatedBy(comma, pair).Optional([]any{})).Skip(objEnd).Map(buildObj)
)

func init() {
	jsonObj.Set(OneOf(decimal, integer, str, boolean, arr, obj))
}

func Parse(s string) any {
	r, err := jsonObj.ParseToEnd(s)
	if err != nil {
		panic(err)
	}
	return r
}
