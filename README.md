# parserc-go

parserc-go是用golang实现的解析器组合子（Parser Combinator）库，可以方便地以自底向上的方式构建复杂的解析器。

## 计算器示例

```go
package main

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

func eval(s string) float64 {
    r, err := expr.ParseToEnd(s)
    if err != nil {
        panic(err)
    }
    return r.(float64)
}

func main() {
    fmt.Println(eval("77.58* ( 6 / 3.14+55.2234 ) -2 * 6.1/ ( 1.0+2/ (4.0-3.8*5))"))
}
```

## JSON示例

```go
package main

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

func parse(s string) any {
    r, err := jsonObj.ParseToEnd(s)
    if err != nil {
        panic(err)
    }
    return r
}

func main() {
    fmt.Println(parse(`
    {
        "a": 123,
        "b": 3.14,
        "c": "hello",
        "d": {
            "x": 100,
            "y": "world!"
        },
        "e": [
            12,
            34.56,
            {
                "name": "Xiao Ming",
                "age": 18,
                "score": [99.8, 87.5, 60.0]
            },
            "abc"
        ],
        "f": [],
        "g": {},
        "h": [true, {"m": false}]
    }`))
}

```