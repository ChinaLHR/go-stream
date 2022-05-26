package stream

import (
	"github.com/chinalhr/go-stream/types"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

//get data source test

func TestStream_OfElements(t *testing.T) {
	tests := []struct {
		name   string
		input  []types.T
		actual []types.T
	}{
		{
			name:   "normalCase",
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{1, 2, 3, 4, 5},
		},
		{
			name:   "emptyCase",
			input:  []types.T{},
			actual: []types.T{},
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream := OfElements(test.input...)
			iterator := stream.p.it
			result := make([]types.T, 0, 5)
			for iterator.HasNext() {
				result = append(result, iterator.Next())
			}
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_OfSlice(t *testing.T) {
	tests := []struct {
		name   string
		input  []types.T
		actual []types.T
	}{
		{
			name:   "normalCase",
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{1, 2, 3, 4, 5},
		},
		{
			name:   "emptyCase",
			input:  []types.T{},
			actual: []types.T{},
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream := OfSlice(test.input)
			iterator := stream.p.it
			result := make([]types.T, 0, 5)
			for iterator.HasNext() {
				result = append(result, iterator.Next())
			}
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_OfMap(t *testing.T) {
	tests := []struct {
		name   string
		input  map[int]string
		actual map[int]string
	}{
		{
			name: "normalCase",
			input: map[int]string{
				1: "A", 2: "B", 3: "C",
			},
			actual: map[int]string{
				1: "A", 2: "B", 3: "C",
			},
		},
		{
			name:   "emptyCase",
			input:  map[int]string{},
			actual: map[int]string{},
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: map[int]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream := OfMap(test.input)
			iterator := stream.p.it
			result := make(map[int]string)
			for iterator.HasNext() {
				kv := iterator.Next().(types.KV)
				result[kv.KEY.(int)] = kv.VALUE.(string)
			}
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Generate(t *testing.T) {
	var initialVal = 0

	tests := []struct {
		name   string
		input  func() types.T
		actual []int
	}{
		{
			name: "normalCase",
			input: func() types.T {
				initialVal = initialVal + 1
				return initialVal
			},
			actual: []int{1, 2, 3, 4, 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream := Generate(test.input)
			iterator := stream.p.it
			result := make([]int, 0, 5)
			for i := 0; i < 5; i++ {
				assert.Equal(t, iterator.HasNext(), true)
				assert.Equal(t, iterator.GetSize(), -1)
				result = append(result, iterator.Next().(int))
			}
			assert.Equal(t, result, test.actual)
		})
	}
}

//intermediateStage stateless operation test

func TestStream_Filter(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(e types.T) bool
		input  []types.T
		actual []types.T
	}{
		{
			name: "normalCase",
			fn: func(e types.T) bool {
				return e.(int)%2 == 0
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{2, 4},
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Filter(test.fn).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Map(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(e types.T) types.R
		input  []types.T
		actual []types.T
	}{
		{
			name: "normalCase",
			fn: func(e types.T) types.R {
				return strconv.Itoa(e.(int) * e.(int))
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{"1", "4", "9", "16", "25"},
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Map(test.fn).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Peek(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(e types.T)
		input  []types.T
		actual []types.T
	}{
		{
			name: "normalCase",
			fn: func(e types.T) {
				strconv.Itoa(e.(int) * e.(int))
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{1, 2, 3, 4, 5},
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Peek(test.fn).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_FlatMap(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(t types.T) Stream
		input  [][]types.T
		actual []types.T
	}{
		{
			name: "normalCase",
			fn: func(t types.T) Stream {
				return OfSlice(t)
			},
			input:  [][]types.T{{1, 2, 3}, {4, 5}},
			actual: []types.T{1, 2, 3, 4, 5},
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				FlatMap(test.fn).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

//intermediateStage stateful operation test

func TestStream_Distinct(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(item types.T) types.R
		input  []types.T
		actual []types.T
	}{
		{
			name: "normalCase",
			fn: func(t types.T) types.R {
				return t
			},
			input:  []types.T{1, 1, 2, 2, 3, 4, 5, 5, 5},
			actual: []types.T{1, 2, 3, 4, 5},
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Distinct(test.fn).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Sorted(t *testing.T) {
	tests := []struct {
		name    string
		compare func(first types.T, second types.T) int
		input   []types.T
		actual  []types.T
	}{
		{
			name: "normalCase",
			compare: func(first types.T, second types.T) int {
				return first.(int) - second.(int)
			},
			input:  []types.T{3, 2, 5, 1, 4},
			actual: []types.T{1, 2, 3, 4, 5},
		},
		{
			name:    "nilCase",
			compare: nil,
			input:   nil,
			actual:  []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Sorted(test.compare).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Skip(t *testing.T) {
	tests := []struct {
		name   string
		skip   int
		input  []types.T
		actual []types.T
	}{
		{
			name:   "normalCase",
			skip:   3,
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{4, 5},
		},
		{
			name:   "nilCase",
			skip:   0,
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Skip(test.skip).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Limit(t *testing.T) {
	tests := []struct {
		name   string
		limit  int
		input  []types.T
		actual []types.T
	}{
		{
			name:   "normalCase",
			limit:  3,
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{1, 2, 3},
		},
		{
			name:   "nilCase",
			limit:  0,
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Limit(test.limit).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_TakeWhile(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(t types.T) bool
		input  []types.T
		actual []types.T
	}{
		{
			name: "normalCase",
			fn: func(t types.T) bool {
				return t.(int) <= 3
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{1, 2, 3},
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				TakeWhile(test.fn).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_DropWhile(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(t types.T) bool
		input  []types.T
		actual []types.T
	}{
		{
			name: "normalCase",
			fn: func(t types.T) bool {
				return t.(int) <= 3
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{4, 5},
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				DropWhile(test.fn).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

//terminal non-short-circuiting operation

func TestStream_ForEach(t *testing.T) {

	tests := []struct {
		name   string
		mapFn  func(e types.T) types.R
		input  []types.T
		actual []types.T
	}{
		{
			name: "normalCase",
			mapFn: func(e types.T) types.R {
				return strconv.Itoa(e.(int))
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: []types.T{"1", "2", "3", "4", "5"},
		},
	}

	for _, test := range tests {
		forEachIndex := 0
		t.Run(test.name, func(t *testing.T) {
			OfSlice(test.input).
				Map(test.mapFn).
				ForEach(func(e types.T) {
					assert.Equal(t, e, test.actual[forEachIndex])
					forEachIndex++
				})

		})
	}
}

func TestStream_FindLast(t *testing.T) {
	tests := []struct {
		name   string
		mapFn  func(e types.T) types.R
		input  []types.T
		actual types.T
	}{
		{
			name:   "normalCase",
			input:  []types.T{1, 2, 3, 4, 5},
			actual: 5,
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				FindLast()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Reduce(t *testing.T) {
	tests := []struct {
		name        string
		accumulator func(e1 types.T, e2 types.T) types.T
		input       []types.T
		actual      types.T
	}{
		{
			name: "normalCase",
			accumulator: func(e1 types.T, e2 types.T) types.T {
				return e1.(int) + e2.(int)
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: 15,
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Reduce(test.accumulator)
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_ReduceFromIdentity(t *testing.T) {
	tests := []struct {
		name        string
		identity    types.T
		accumulator func(e1 types.T, e2 types.T) types.T
		input       []types.T
		actual      types.T
	}{
		{
			name:     "normalCase",
			identity: 10,
			accumulator: func(e1 types.T, e2 types.T) types.T {
				return e1.(int) + e2.(int)
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: 25,
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				ReduceFromIdentity(test.identity, test.accumulator)
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Count(t *testing.T) {
	tests := []struct {
		name   string
		input  []types.T
		actual int
	}{
		{
			name:   "normalCase",
			input:  []types.T{1, 2, 3, 4, 5},
			actual: 5,
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Count()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Max(t *testing.T) {
	tests := []struct {
		name    string
		compare func(first types.T, second types.T) int
		input   []types.T
		actual  types.T
	}{
		{
			name: "normalCase",
			compare: func(first types.T, second types.T) int {
				return first.(int) - second.(int)
			},
			input:  []types.T{6, 5, 9, 4, 2},
			actual: 9,
		},
		{
			name:    "nilCase",
			compare: nil,
			input:   nil,
			actual:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Max(test.compare)
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_Min(t *testing.T) {
	tests := []struct {
		name    string
		compare func(first types.T, second types.T) int
		input   []types.T
		actual  types.T
	}{
		{
			name: "normalCase",
			compare: func(first types.T, second types.T) int {
				return first.(int) - second.(int)
			},
			input:  []types.T{6, 5, 9, 4, 2},
			actual: 2,
		},
		{
			name:    "nilCase",
			compare: nil,
			input:   nil,
			actual:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				Min(test.compare)
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_ToSlice(t *testing.T) {
	tests := []struct {
		name   string
		input  []types.T
		actual []types.T
	}{
		{
			name:   "normalCase",
			input:  []types.T{5, 4, 3, 2, 1},
			actual: []types.T{5, 4, 3, 2, 1},
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: []types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				ToSlice()
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_ToMap(t *testing.T) {
	tests := []struct {
		name        string
		keyMapper   func(t types.T) types.K
		valueMapper func(t types.T) types.R
		input       []types.T
		actual      map[types.K]types.R
	}{
		{
			name: "normalCase",
			keyMapper: func(t types.T) types.K {
				return t
			},
			valueMapper: func(t types.T) types.R {
				return "[" + strconv.Itoa(t.(int)) + "]"
			},
			input: []types.T{5, 4, 3, 2, 1},
			actual: map[types.K]types.R{
				5: "[5]",
				4: "[4]",
				3: "[3]",
				2: "[2]",
				1: "[1]",
			},
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: map[types.K]types.R{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				ToMap(test.keyMapper, test.valueMapper)
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_GroupingBy(t *testing.T) {
	tests := []struct {
		name       string
		classifier func(t types.T) types.K
		input      []types.T
		actual     map[types.K][]types.T
	}{
		{
			name: "normalCase",
			classifier: func(t types.T) types.K {
				if t.(int)%2 == 0 {
					return true
				} else {
					return false
				}
			},
			input: []types.T{1, 2, 3, 4, 5, 6},
			actual: map[types.K][]types.T{
				true: {
					2, 4, 6,
				},
				false: {
					1, 3, 5,
				},
			},
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: map[types.K][]types.T{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				GroupingBy(test.classifier)
			assert.Equal(t, result, test.actual)
		})
	}
}

//terminal short-circuiting operation test
func TestStream_AllMatch(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(e types.T) bool
		input  []types.T
		actual bool
	}{
		{
			name: "nonMatchCase",
			fn: func(e types.T) bool {
				if e.(int)%2 == 0 {
					return true
				}
				return false
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: false,
		},
		{
			name: "MatchCase",
			fn: func(e types.T) bool {
				if e.(int) < 10 {
					return true
				}
				return false
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: true,
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				AllMatch(test.fn)
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_AnyMatch(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(e types.T) bool
		input  []types.T
		actual bool
	}{
		{
			name: "nonMatchCase",
			fn: func(e types.T) bool {
				if e.(int) > 10 {
					return true
				}
				return false
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: false,
		},
		{
			name: "MatchCase",
			fn: func(e types.T) bool {
				if e.(int) == 1 {
					return true
				}
				return false
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: true,
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				AnyMatch(test.fn)
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_NoneMatch(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(e types.T) bool
		input  []types.T
		actual bool
	}{
		{
			name: "nonMatchCase",
			fn: func(e types.T) bool {
				if e.(int)%2 == 0 {
					return true
				}
				return false
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: false,
		},
		{
			name: "MatchCase",
			fn: func(e types.T) bool {
				if e.(int) > 10 {
					return true
				}
				return false
			},
			input:  []types.T{1, 2, 3, 4, 5},
			actual: true,
		},
		{
			name:   "nilCase",
			fn:     nil,
			input:  nil,
			actual: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				NoneMatch(test.fn)
			assert.Equal(t, result, test.actual)
		})
	}
}

func TestStream_FindFirst(t *testing.T) {
	tests := []struct {
		name   string
		input  []types.T
		actual types.T
	}{
		{
			name:   "normalCase",
			input:  []types.T{1, 2, 3, 4, 5},
			actual: 1,
		},
		{
			name:   "nilCase",
			input:  nil,
			actual: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := OfSlice(test.input).
				FindFirst()
			assert.Equal(t, result, test.actual)
		})
	}
}
