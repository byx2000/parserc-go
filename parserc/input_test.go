package parserc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateInput(t *testing.T) {
	input := CreateInput("abc")
	assert.Equal(t, "abc", input.str)
	assert.Equal(t, 0, input.index)
	assert.Equal(t, 1, input.row)
	assert.Equal(t, 1, input.col)
}

func TestInput_Next(t *testing.T) {
	str := `abc
def
ghi`
	input := CreateInput(str)
	assert.Equal(t, 'a', input.Current())
	assert.Equal(t, 1, input.row)
	assert.Equal(t, 1, input.col)

	input = input.Next().Next()
	assert.Equal(t, 'c', input.Current())
	assert.Equal(t, 1, input.row)
	assert.Equal(t, 3, input.col)

	input = input.Next().Next()
	assert.Equal(t, 'd', input.Current())
	assert.Equal(t, 2, input.row)
	assert.Equal(t, 1, input.col)

	input = input.Next().Next().Next().Next().Next()
	assert.Equal(t, 'h', input.Current())
	assert.Equal(t, 3, input.row)
	assert.Equal(t, 2, input.col)

	assert.False(t, input.End())
	input = input.Next().Next()
	assert.True(t, input.End())
}
