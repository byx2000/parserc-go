package parserc

import "fmt"

// Pair 长度为2的元组
type Pair struct {
	First  any
	Second any
}

func (p Pair) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}
