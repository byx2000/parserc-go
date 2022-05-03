package parserc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func verifySuccess(t *testing.T, p *Parser, input string, expectedResult any) {
	r, e := p.ParseToEnd(input)
	assert.Nil(t, e)
	assert.Equal(t, expectedResult, r)
}

func verifyFailed(t *testing.T, p *Parser, input string) {
	_, e := p.ParseToEnd(input)
	assert.NotNil(t, e)
}

func TestFail(t *testing.T) {
	verifyFailed(t, Fail("error message"), "")
	verifyFailed(t, Fail("error message"), "abc")
}

func TestAny(t *testing.T) {
	verifyFailed(t, Any(), "")
	verifySuccess(t, Any(), "a", 'a')
	verifySuccess(t, Any(), "b", 'b')
	verifySuccess(t, Any(), "c", 'c')
}

func TestCh(t *testing.T) {
	verifyFailed(t, Ch('a'), "")
	verifySuccess(t, Ch('a'), "a", 'a')
	verifyFailed(t, Ch('a'), "b")
}

func TestChs(t *testing.T) {
	verifyFailed(t, Chs('a', 'b', 'c'), "")
	verifySuccess(t, Chs('a', 'b', 'c'), "a", 'a')
	verifySuccess(t, Chs('a', 'b', 'c'), "b", 'b')
	verifySuccess(t, Chs('a', 'b', 'c'), "c", 'c')
	verifyFailed(t, Chs('a', 'b', 'c'), "d")
}

func TestNot(t *testing.T) {
	verifyFailed(t, Not('a'), "")
	verifyFailed(t, Not('a'), "a")
	verifySuccess(t, Not('a'), "b", 'b')
}

func TestRange(t *testing.T) {
	verifyFailed(t, Range('d', 'f'), "")
	verifySuccess(t, Range('d', 'f'), "d", 'd')
	verifySuccess(t, Range('d', 'f'), "e", 'e')
	verifySuccess(t, Range('d', 'f'), "f", 'f')
	verifyFailed(t, Range('d', 'f'), "c")
	verifyFailed(t, Range('d', 'f'), "g")
}

func TestStr(t *testing.T) {
	verifyFailed(t, Str("abc"), "")
	verifyFailed(t, Str("abc"), "a")
	verifyFailed(t, Str("abc"), "ab")
	verifyFailed(t, Str("abc"), "bc")
	verifyFailed(t, Str("abc"), "ac")
	verifySuccess(t, Str("abc"), "abc", "abc")
}

func TestMap(t *testing.T) {
	verifySuccess(t, Str("abc").Map(func(r any) any {
		return r.(string) + " hello"
	}), "abc", "abc hello")
	verifyFailed(t, Str("abc").Map(func(r any) any {
		return r.(string) + " hello"
	}), "ab")
}

func TestAnd(t *testing.T) {
	verifySuccess(t, Ch('a').And(Ch('b')), "ab", Pair{'a', 'b'})
	verifyFailed(t, Ch('a').And(Ch('b')), "")
	verifyFailed(t, Ch('a').And(Ch('b')), "a")
	verifyFailed(t, Ch('a').And(Ch('b')), "b")
	verifyFailed(t, Ch('a').And(Ch('b')), "ac")
	verifyFailed(t, Ch('a').And(Ch('b')), "cb")
	verifyFailed(t, Ch('a').And(Ch('b')), "xy")
}

func TestSeq(t *testing.T) {
	verifySuccess(t, Seq(Ch('a'), Str(" hello "), Ch('b')), "a hello b", []any{'a', " hello ", 'b'})
	verifyFailed(t, Seq(Ch('a'), Str(" hello "), Ch('b')), "")
	verifyFailed(t, Seq(Ch('a'), Str(" hello "), Ch('b')), "a")
	verifyFailed(t, Seq(Ch('a'), Str(" hello "), Ch('b')), "a hi b")
	verifyFailed(t, Seq(Ch('a'), Str(" hello "), Ch('b')), "a hello c")
}

func TestOr(t *testing.T) {
	verifySuccess(t, Ch('a').Or(Ch('b')), "a", 'a')
	verifySuccess(t, Ch('a').Or(Ch('b')), "b", 'b')
	verifyFailed(t, Ch('a').Or(Ch('b')), "x")
	verifyFailed(t, Ch('a').Or(Ch('b')), "")
	verifyFailed(t, Str("a").Or(Str("ab")), "ab")
}

func TestOneOf(t *testing.T) {
	verifySuccess(t, OneOf(Str("apple"), Str("banana"), Str("cat")), "apple", "apple")
	verifySuccess(t, OneOf(Str("apple"), Str("banana"), Str("cat")), "banana", "banana")
	verifySuccess(t, OneOf(Str("apple"), Str("banana"), Str("cat")), "cat", "cat")
	verifyFailed(t, OneOf(Str("apple"), Str("banana"), Str("cat")), "doctor")
	verifyFailed(t, OneOf(Str("apple"), Str("banana"), Str("cat")), "")
	verifyFailed(t, OneOf(Str("a"), Str("ab")), "ab")
}

func TestSkip(t *testing.T) {
	verifySuccess(t, Skip(Ch('a')).And(Ch('b')), "ab", 'b')
	verifySuccess(t, Ch('a').Skip(Ch('b')), "ab", 'a')
}

func TestMany(t *testing.T) {
	verifySuccess(t, Ch('a').Many(), "", []any{})
	verifySuccess(t, Ch('a').Many(), "a", []any{'a'})
	verifySuccess(t, Ch('a').Many(), "aaa", []any{'a', 'a', 'a'})
}

func TestMany1(t *testing.T) {
	verifyFailed(t, Ch('a').Many1(), "")
	verifySuccess(t, Ch('a').Many1(), "a", []any{'a'})
	verifySuccess(t, Ch('a').Many1(), "aaa", []any{'a', 'a', 'a'})
}

func TestOptional(t *testing.T) {
	verifySuccess(t, Ch('a').Opt('x'), "", 'x')
	verifySuccess(t, Ch('a').Opt('x'), "a", 'a')
}

func TestPeek(t *testing.T) {
	verifySuccess(t, Peek(Str("ab"), Str("abc"), Str("def")), "abc", "abc")
	verifySuccess(t, Peek(Str("ab"), Str("abc"), Str("def")), "def", "def")
	verifyFailed(t, Peek(Str("ab"), Str("abc"), Str("def")), "abx")
	verifyFailed(t, Peek(Str("ab"), Str("abc"), Str("def")), "deg")
}

func TestSeparatedBy(t *testing.T) {
	verifySuccess(t, Separate(Ch(','), Any()), "a,b,c", []any{'a', 'b', 'c'})
	verifySuccess(t, Separate(Ch(','), Any()), "a", []any{'a'})
	verifyFailed(t, Separate(Ch(','), Any()), "")
}

func TestSurroundedBy(t *testing.T) {
	verifySuccess(t, Ch('a').Surround(Ch('b')), "bab", 'a')
	verifyFailed(t, Ch('a').Surround(Ch('b')), "bax")
	verifyFailed(t, Ch('a').Surround(Ch('b')), "xab")
	verifyFailed(t, Ch('a').Surround(Ch('b')), "xay")
	verifyFailed(t, Ch('a').Surround(Ch('b')), "bmb")
	verifyFailed(t, Ch('a').Surround(Ch('b')), "ba")
	verifyFailed(t, Ch('a').Surround(Ch('b')), "ab")
	verifyFailed(t, Ch('a').Surround(Ch('b')), "a")
	verifyFailed(t, Ch('a').Surround(Ch('b')), "")
}

func TestFatal(t *testing.T) {
	verifySuccess(t, Ch('a').Fatal(), "a", 'a')
	assert.Panics(t, func() {
		_, _ = Ch('a').Fatal().ParseToEnd("b")
	})
}

func TestDelaySet(t *testing.T) {
	p1 := NewParser()
	p2 := p1.And(Ch('b'))
	p1.Set(Ch('a'))
	verifySuccess(t, p2, "ab", Pair{'a', 'b'})
	verifyFailed(t, p2, "")
	verifyFailed(t, p2, "a")
}
