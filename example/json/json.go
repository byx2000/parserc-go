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
	integer  = digits.Map(toInt).Surround(ws)
	decimal  = Seq(digits, Ch('.'), digits).Surround(ws).Map(join).Map(toFloat)
	str      = Skip(Ch('"')).And(Not('"').Many()).Skip(Ch('"')).Map(join).Surround(ws)
	boolean  = Str("true").Or(Str("false")).Map(toBool).Surround(ws)
	objStart = Ch('{').Surround(ws)
	objEnd   = Ch('}').Surround(ws)
	arrStart = Ch('[').Surround(ws)
	arrEnd   = Ch(']').Surround(ws)
	colon    = Ch(':').Surround(ws)
	comma    = Ch(',').Surround(ws)
	jsonObj  = NewParser()
	arr      = Skip(arrStart).And(Separate(comma, jsonObj).Opt([]any{})).Skip(arrEnd)
	pair     = str.Skip(colon).And(jsonObj)
	obj      = Skip(objStart).And(Separate(comma, pair).Opt([]any{})).Skip(objEnd).Map(buildObj)
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
