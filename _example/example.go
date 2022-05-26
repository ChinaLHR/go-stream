package main

import (
	"fmt"
	"github.com/chinalhr/go-stream"
	"github.com/chinalhr/go-stream/types"
)

func main() {

	type widget struct {
		color  string
		weight int
	}

	widgets := []widget{
		{
			color:  "yellow",
			weight: 4,
		},
		{
			color:  "red",
			weight: 3,
		},
		{
			color:  "yellow",
			weight: 2,
		},
		{
			color:  "blue",
			weight: 1,
		},
	}

	sum := stream.OfSlice(widgets).
		Filter(func(e types.T) bool {
			return e.(widget).color == "yellow"
		}).
		Map(func(e types.T) (r types.R) {
			return e.(widget).weight
		}).
		Reduce(func(e1 types.T, e2 types.T) types.T {
			return e1.(int) + e2.(int)
		})

	fmt.Printf("sum:%v", sum.(int))
}
