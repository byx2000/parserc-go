package parserc

import (
	"errors"
	"fmt"
)

// ParseResult 解析结果
type ParseResult struct {
	Result any   // 结果
	Remain Input // 剩余输入
}

var emptyParseResult = ParseResult{}

// ParseFunc 解析函数
type ParseFunc func(Input) (ParseResult, error)

// Parser 解析器
type Parser struct {
	parse ParseFunc
}

func parseError(input Input, msg string) error {
	return errors.New(fmt.Sprintf("parse error at row %d, col %d: %s", input.Row(), input.Col(), msg))
}

// Fail 直接失败
func Fail(msg string) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		return emptyParseResult, parseError(input, msg)
	}}
}

// Any 匹配任意字符
func Any() *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		if input.End() {
			return emptyParseResult, parseError(input, "unexpected end of input")
		}
		c := input.Current()
		return ParseResult{c, input.Next()}, nil
	}}
}

// Ch 匹配指定字符
func Ch(c rune) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		if input.End() {
			return emptyParseResult, parseError(input, "unexpected end of input")
		}
		ch := input.Current()
		if c != ch {
			return emptyParseResult, parseError(input, fmt.Sprintf("expected %c", c))
		}
		return ParseResult{c, input.Next()}, nil
	}}
}

// Chs 匹配字符集
func Chs(chs ...rune) *Parser {
	set := make(map[rune]bool)
	for _, c := range chs {
		set[c] = true
	}
	return &Parser{func(input Input) (ParseResult, error) {
		if input.End() {
			return emptyParseResult, parseError(input, "unexpected end of input")
		}
		c := input.Current()
		_, exist := set[c]
		if !exist {
			return emptyParseResult, parseError(input, fmt.Sprintf("unexpected %c", c))
		}
		return ParseResult{c, input.Next()}, nil
	}}
}

// Not 匹配不等于指定字符的字符
func Not(c rune) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		if input.End() {
			return emptyParseResult, parseError(input, "unexpected end of input")
		}
		ch := input.Current()
		if c == ch {
			return emptyParseResult, parseError(input, fmt.Sprintf("unexpected %c", ch))
		}
		return ParseResult{ch, input.Next()}, nil
	}}
}

// Range 匹配指定范围内的字符
func Range(c1 rune, c2 rune) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		if input.End() {
			return emptyParseResult, parseError(input, "unexpected end of input")
		}
		c := input.Current()
		if (c-c1)*(c-c2) > 0 {
			return emptyParseResult, parseError(input, fmt.Sprintf("unexpected %c", c))
		}
		return ParseResult{c, input.Next()}, nil
	}}
}

// Str 匹配字符串前缀
func Str(s string) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		i := input
		for _, c := range s {
			if i.End() || i.Current() != c {
				return emptyParseResult, parseError(input, fmt.Sprintf("expected %s", s))
			}
			i = i.Next()
		}
		return ParseResult{s, i}, nil
	}}
}

// Map 转换解析结果
func Map(p *Parser, mapper func(any) any) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		r, err := p.parse(input)
		if err != nil {
			return emptyParseResult, err
		}
		return ParseResult{mapper(r.Result), r.Remain}, nil
	}}
}

// And 连接两个解析器
func And(lhs *Parser, rhs *Parser) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		r1, err := lhs.parse(input)
		if err != nil {
			return emptyParseResult, err
		}
		r2, err := rhs.parse(r1.Remain)
		if err != nil {
			return emptyParseResult, err
		}
		return ParseResult{Pair{r1.Result, r2.Result}, r2.Remain}, nil
	}}
}

// Seq 连接多个解析器
func Seq(parsers ...*Parser) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		rs := make([]any, 0)
		for _, p := range parsers {
			r, err := p.parse(input)
			if err != nil {
				return emptyParseResult, err
			}
			rs = append(rs, r.Result)
			input = r.Remain
		}
		return ParseResult{rs, input}, nil
	}}
}

// Or 有序选择两个解析器
func Or(lhs *Parser, rhs *Parser) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		r, err := lhs.parse(input)
		if err == nil {
			return r, nil
		}
		r, err = rhs.parse(input)
		if err != nil {
			return emptyParseResult, err
		}
		return r, nil
	}}
}

// OneOf 有序选择多个解析器
func OneOf(p1 *Parser, p2 *Parser, parsers ...*Parser) *Parser {
	p := Or(p1, p2)
	for _, pp := range parsers {
		p = Or(p, pp)
	}
	return p
}

// SkipFirst 连接两个解析器，并丢弃第一个解析器的结果
func SkipFirst(p1 *Parser, p2 *Parser) *Parser {
	return p1.And(p2).Map(func(p any) any {
		return p.(Pair).Second
	})
}

// SkipSecond 连接两个解析器，并丢弃第二个解析器的结果
func SkipSecond(p1 *Parser, p2 *Parser) *Parser {
	return p1.And(p2).Map(func(p any) any {
		return p.(Pair).First
	})
}

type SkipWrapper struct {
	lhs *Parser
	And func(*Parser) *Parser
}

func Skip(lhs *Parser) SkipWrapper {
	return SkipWrapper{lhs, func(rhs *Parser) *Parser {
		return SkipFirst(lhs, rhs)
	}}
}

// Many 应用指定解析器零次或多次
func Many(p *Parser) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		rs := make([]any, 0)
		for {
			r, err := p.parse(input)
			if err != nil {
				break
			}
			rs = append(rs, r.Result)
			input = r.Remain
		}
		return ParseResult{rs, input}, nil
	}}
}

// Many1 应用指定解析器一次或多次
func Many1(p *Parser) *Parser {
	return p.And(p.Many()).Map(func(p any) any {
		pair := p.(Pair)
		rs := make([]any, 0)
		rs = append(rs, pair.First)
		rs = append(rs, pair.Second.([]any)...)
		return rs
	})
}

// Optional 尝试应用解析器，并在失败时返回默认值
func Optional(p *Parser, defaultValue any) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		r, err := p.parse(input)
		if err != nil {
			return ParseResult{defaultValue, input}, nil
		}
		return r, nil
	}}
}

// Peek 根据probe的执行成功与否，选择执行success或failed
func Peek(probe *Parser, success *Parser, failed *Parser) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		_, err := probe.parse(input)
		if err != nil {
			return failed.parse(input)
		}
		return success.parse(input)
	}}
}

// SeparatedBy 匹配被给定分隔符分隔的输入
func SeparatedBy(delimiter *Parser, p *Parser) *Parser {
	return p.And(Skip(delimiter).And(p).Many()).Map(func(p any) any {
		var result []any
		result = append(result, p.(Pair).First)
		for _, e := range p.(Pair).Second.([]any) {
			result = append(result, e)
		}
		return result
	})
}

// Fatal 指定解析器解析失败时，抛出关键错误
func Fatal(p *Parser) *Parser {
	return &Parser{func(input Input) (ParseResult, error) {
		r, e := p.parse(input)
		if e != nil {
			panic(e)
		}
		return r, nil
	}}
}

// NewParser 创建空解析器，该解析器随后通过Set方法设置
func NewParser() *Parser {
	return &Parser{nil}
}

// ParseToEnd 解析输入直到末尾
func (p Parser) ParseToEnd(s string) (any, error) {
	r, err := p.parse(CreateInput(s))
	if err != nil {
		return nil, err
	}
	remain := r.Remain
	if !remain.End() {
		return nil, parseError(remain, "end of input not reached")
	}
	return r.Result, nil
}

// Set 设置解析器
func (p *Parser) Set(parser *Parser) {
	p.parse = parser.parse
}

// And 连接另一个解析器
func (p *Parser) And(rhs *Parser) *Parser {
	return And(p, rhs)
}

// Or 有序选择另一个解析器
func (p *Parser) Or(rhs *Parser) *Parser {
	return Or(p, rhs)
}

// Many 应用当前解析器零次或多次
func (p *Parser) Many() *Parser {
	return Many(p)
}

// Many1 应用当前解析器一次或多次
func (p *Parser) Many1() *Parser {
	return Many1(p)
}

// Map 转换当前解析器的解析结果
func (p *Parser) Map(mapper func(any) any) *Parser {
	return Map(p, mapper)
}

// Skip 连接另一个解析器并丢弃解析结果
func (p *Parser) Skip(rhs *Parser) *Parser {
	return SkipSecond(p, rhs)
}

// SurroundedBy 在当前解析器周围应用另一个解析器
func (p *Parser) SurroundedBy(parser *Parser) *Parser {
	return Seq(parser, p, parser).Map(func(rs any) any {
		return rs.([]any)[1]
	})
}

// ManyUntil 应用当前解析器零次或多次，直到指定解析器执行成功
func (p *Parser) ManyUntil(until *Parser) *Parser {
	return Peek(until, Fail("no error message"), p).Many()
}

// Optional 将当前解析器变为可选，并提供默认解析结果
func (p *Parser) Optional(defaultValue any) *Parser {
	return Optional(p, defaultValue)
}

// Fatal 当前解析器失败时，抛出关键错误
func (p *Parser) Fatal() *Parser {
	return Fatal(p)
}
