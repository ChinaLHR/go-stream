package stream

import (
	"github.com/chinalhr/go-stream/types"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func BenchmarkTestStream_Parallel_CPU(b *testing.B) {
	benchmarkInput := sequenceSlice(10)
	randCount := 50000

	tests := []struct {
		name   string
		worker int
	}{
		{
			name:   "non worker",
			worker: 0,
		},
		{
			name:   "two worker",
			worker: 2,
		},
		{
			name:   "four worker",
			worker: 4,
		},
		{
			name:   "eight worker",
			worker: 8,
		},
		{
			name:   "sixteen worker",
			worker: 16,
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				OfSlice(benchmarkInput).
					Parallel(test.worker).
					Map(func(e types.T) (r types.R) {
						slice := randSlice(randCount)
						return slice
					}).
					ForEach(func(e types.T) {
						sort.Ints(e.([]int))
					})
			}
		})
	}
}

func BenchmarkTestStream_Parallel_IO(b *testing.B) {
	benchmarkInput := sequenceSlice(20)
	mockIOWaitingTime := time.Millisecond * 500

	tests := []struct {
		name   string
		worker int
	}{
		{
			name:   "non worker",
			worker: 0,
		},
		{
			name:   "two worker",
			worker: 2,
		},
		{
			name:   "four worker",
			worker: 4,
		},
		{
			name:   "eight worker",
			worker: 8,
		},
		{
			name:   "sixteen worker",
			worker: 16,
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				OfSlice(benchmarkInput).
					Parallel(test.worker).
					Map(func(e types.T) (r types.R) {
						time.Sleep(mockIOWaitingTime)
						return e
					}).
					ForEach(func(e types.T) {

					})
			}
		})
	}
}

func randSlice(count int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s := make([]int, count)
	for i := 0; i < count; i++ {
		s[i] = r.Intn(count)
	}
	return s
}

func sequenceSlice(count int) []int {
	s := make([]int, count)
	for i := 0; i < count; i++ {
		s[i] = i
	}
	return s
}
