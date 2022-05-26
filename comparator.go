package stream

import "github.com/chinalhr/go-stream/types"

type Comparator struct {
	source  []types.T
	compare func(e1 types.T, e2 types.T) int
}

func (c *Comparator) Len() int {
	return len(c.source)
}

func (c *Comparator) Less(i, j int) bool {
	return c.compare(c.source[i], c.source[j]) < 0
}

func (c *Comparator) Swap(i, j int) {
	c.source[i], c.source[j] = c.source[j], c.source[i]
}
