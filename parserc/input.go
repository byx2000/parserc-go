package parserc

import "unicode/utf8"

// Input 输入流
type Input struct {
	str   string
	index int
	row   int
	col   int
}

// CreateInput 创建输入流
func CreateInput(s string) Input {
	return Input{s, 0, 1, 1}
}

// End 判断是否到达输入流末尾
func (p Input) End() bool {
	return p.index == utf8.RuneCountInString(p.str)
}

// Next 输入流向后移一位
func (p Input) Next() Input {
	row := p.row
	col := p.col + 1
	if p.Current() == '\n' {
		row++
		col = 1
	}
	return Input{p.str, p.index + 1, row, col}
}

// Current 获取当前字符
func (p Input) Current() rune {
	return []rune(p.str)[p.index]
}

// Row 获取当前行号
func (p Input) Row() int {
	return p.row
}

// Col 获取当前列号
func (p Input) Col() int {
	return p.col
}
