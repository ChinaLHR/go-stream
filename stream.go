package stream

import (
	"errors"
	"github.com/chinalhr/go-stream/types"
	"reflect"
	"sort"
)

//Stream is a sequence of elements supporting sequential and parallel aggregate operations.
//Stream operations are combined into a stream pipeline, and computational processing is performed during terminal operation.
//Stream consists of a source(OfElements OfSlice OfMap Generate), zero or more intermediate operations(Filter Map Peek FlatMap
//Distinct Sorted Skip Limit TakeWhile DropWhile), and terminal operations(ForEach FindLast FindFirst Reduce ReduceFromIdentity
//Count Max Min ToSlice ToMap GroupingBy AllMatch AnyMatch NoneMatch FindFirst).
//Stream has lazy evaluation and short-circuit evaluation.
//Example See: _example/example.go
type Stream struct {
	p *referencePipeline
}

//Create a Stream and set Stream source

//OfElements Return a sequential Stream containing elements.
func OfElements(elements ...types.T) Stream {
	iterator := buildSliceIterator(elements...)
	pipeline := newPipeline(iterator)
	return Stream{pipeline}
}

//OfSlice Return a sequential Stream containing slice.
func OfSlice(slice types.T) Stream {
	if slice == nil {
		return OfElements()
	}
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		panic(errors.New("reflect type is not slice"))
	}
	iterator := buildSliceReflectIterator(reflect.ValueOf(slice))
	pipeline := newPipeline(iterator)
	return Stream{pipeline}
}

//OfMap Return a sequential Stream containing mapValue.
func OfMap(mapValue types.T) Stream {
	if mapValue == nil {
		return OfElements()
	}
	if reflect.TypeOf(mapValue).Kind() != reflect.Map {
		panic(errors.New("reflect type is not map"))
	}

	iterator := buildMapReflectIterator(reflect.ValueOf(mapValue))
	pipeline := newPipeline(iterator)
	return Stream{pipeline}
}

//Generate Return an infinite sequential Stream,elements are generated by the supplier.
func Generate(supplier func() types.T) Stream {
	iterator := buildSupplierIterator(supplier)
	pipeline := newPipeline(iterator)
	return Stream{pipeline}
}

//IntermediateStage stateless operation

//Filter Returns a Stream consisting of the elements of this stream that match the given predicate function.
func (s Stream) Filter(predicate func(e types.T) bool) Stream {
	pipeline := s.p
	pipeline.addOperation(
		func(next stage) stage {
			return newDefaultIntermediateStage(next, beginFunc(func(size int) {
				next.Begin(-1)
			}), acceptFunc(func(e types.T) {
				if predicate(e) {
					next.Accept(e)
				}
			}))
		})
	return s
}

//Map Returns a Stream of elements transformed by the mapper function
func (s Stream) Map(mapper func(e types.T) (r types.R)) Stream {
	pipeline := s.p
	pipeline.addOperation(
		func(next stage) stage {
			return newDefaultIntermediateStage(next, acceptFunc(func(e types.T) {
				r := mapper(e)
				next.Accept(r)
			}))
		})
	return s
}

//Peek Does not transform the Stream, executes the consumer function on the elements in the Stream.
func (s Stream) Peek(consumer func(e types.T)) Stream {
	pipeline := s.p
	pipeline.addOperation(func(next stage) stage {
		return newDefaultIntermediateStage(next, acceptFunc(func(e types.T) {
			consumer(e)
			next.Accept(e)
		}))
	})
	return s
}

//FlatMap Returns a Stream consisting of the results of replacing each element of this Stream with the contents of
//a mapped Stream produced by applying the provided mapper function to each element.
func (s Stream) FlatMap(mapper func(t types.T) Stream) Stream {
	pipeline := s.p
	pipeline.addOperation(func(next stage) stage {
		return newDefaultIntermediateStage(next, beginFunc(func(size int) {
			next.Begin(-1)
		}), acceptFunc(func(e types.T) {
			stream := mapper(e)
			stream.ForEach(next.Accept)
		}))
	})
	return s
}

//IntermediateStage stateful operation

//Distinct Returns a Stream consisting of the distinct elements,confirm the uniqueness of the element through
//the distinctFn
func (s Stream) Distinct(distinctFn func(item types.T) types.R) Stream {
	pipeline := s.p
	pipeline.addOperation(func(next stage) stage {
		var keyMap map[types.R]bool
		return newDefaultIntermediateStage(next, beginFunc(func(size int) {
			keyMap = make(map[types.R]bool)
			next.Begin(-1)
		}), acceptFunc(func(e types.T) {
			key := distinctFn(e)
			_, hasKey := keyMap[key]
			if !hasKey {
				next.Accept(e)
				keyMap[key] = true
			}
		}), endFunc(func() {
			keyMap = nil
			next.End()
		}))
	})
	return s
}

//Sorted Return to orderly Stream. sort elements in Stream based on compare function.
func (s Stream) Sorted(compare func(first types.T, second types.T) int) Stream {
	pipeline := s.p
	pipeline.addOperation(func(next stage) stage {
		var sortedList []types.T
		return newDefaultIntermediateStage(next, beginFunc(func(size int) {
			if size == -1 {
				sortedList = make([]types.T, 0, size)
			} else {
				sortedList = make([]types.T, 0)
			}
			next.Begin(size)
		}), acceptFunc(func(e types.T) {
			sortedList = append(sortedList, e)
		}), endFunc(func() {
			c := &Comparator{
				source:  sortedList,
				compare: compare,
			}
			sort.Sort(c)
			next.Begin(len(sortedList))
			for i := 0; i < len(sortedList) && !next.CancellationRequested(); i++ {
				next.Accept(sortedList[i])
			}
			next.End()
			sortedList = nil
		}))
	})
	return s
}

//Skip Discard the previous n elements, return the Stream of the remaining elements.
func (s Stream) Skip(n int) Stream {
	pipeline := s.p
	pipeline.addOperation(func(next stage) stage {
		var totalSkip = n
		if n < 0 {
			n = 0
		}
		return newDefaultIntermediateStage(next, beginFunc(func(size int) {
			if size == -1 {
				next.Begin(-1)
				return
			}
			if n > size {
				next.Begin(0)
			} else {
				next.Begin(size - n)
			}
		}), acceptFunc(func(e types.T) {
			if totalSkip == 0 {
				next.Accept(e)
			} else {
				totalSkip--
			}
		}))
	})
	return s
}

//Limit Returns a stream consisting of the elements of this Stream, truncated to be no longer than maxSize in length.
func (s Stream) Limit(maxSize int) Stream {
	pipeline := s.p
	pipeline.addOperation(func(next stage) stage {
		var totalLimit = 0
		if maxSize < 0 {
			maxSize = 0
		}
		return newDefaultIntermediateStage(next, beginFunc(func(size int) {
			if size == -1 {
				next.Begin(-1)
				return
			}
			next.Begin(maxSize)

		}), acceptFunc(func(e types.T) {
			if totalLimit < maxSize {
				next.Accept(e)
			}
			totalLimit++
		}), cancellationRequestedFunc(func() bool {
			return totalLimit == maxSize
		}))
	})
	return s
}

//TakeWhile Truncate Stream when function does not match.
func (s Stream) TakeWhile(predicate func(t types.T) bool) Stream {
	pipeline := s.p
	pipeline.addOperation(func(next stage) stage {
		take := true
		return newDefaultIntermediateStage(next, beginFunc(func(size int) {
			next.Begin(-1)
		}), acceptFunc(func(e types.T) {
			if take {
				if predicate(e) {
					next.Accept(e)
				} else {
					take = false
				}
			}
		}), cancellationRequestedFunc(func() bool {
			return !take || next.CancellationRequested()
		}))
	})
	return s
}

//DropWhile When the element matching function, start passing the element to Stream.
func (s Stream) DropWhile(predicate func(t types.T) bool) Stream {
	pipeline := s.p
	pipeline.addOperation(func(next stage) stage {
		take := false
		return newDefaultIntermediateStage(next, beginFunc(func(size int) {
			next.Begin(-1)
		}), acceptFunc(func(e types.T) {
			if take {
				next.Accept(e)
				return
			}
			if !predicate(e) {
				take = true
				next.Accept(e)
			}
		}))
	})
	return s
}

//Terminal non-short-circuiting operation

//ForEach Performs an action for each element of this stream.
func (s Stream) ForEach(action func(e types.T)) {
	pipeline := s.p
	pipeline.evaluate(newDefaultTerminalStage(
		acceptFunc(func(e types.T) {
			action(e)
		}),
	))
}

//FindLast Return The last element of the Stream.
func (s Stream) FindLast() types.T {
	pipeline := s.p
	var result types.T
	pipeline.evaluate(newDefaultTerminalStage(
		acceptFunc(func(e types.T) {
			result = e
		}),
	))
	return result
}

//Reduce Performs a reduction on the elements of this Stream, using the associative
//accumulator function, and returns the reduced value.
func (s Stream) Reduce(accumulator func(e1 types.T, e2 types.T) types.T) types.T {
	pipeline := s.p
	var empty bool
	var state types.T
	pipeline.evaluate(newDefaultTerminalStage(
		beginFunc(func(i int) {
			empty = true
			state = nil
		}),
		acceptFunc(func(e types.T) {
			if empty {
				empty = false
				state = e
			} else {
				state = accumulator(state, e)
			}
		}),
	))
	return state
}

//ReduceFromIdentity Performs a reduction on the elements of this Stream, using the provided identity value
//and an associative accumulator function, and returns the reduced value.
func (s Stream) ReduceFromIdentity(identity types.T, accumulator func(e1 types.T, e2 types.T) types.T) types.T {
	pipeline := s.p
	var state = identity
	pipeline.evaluate(newDefaultTerminalStage(
		acceptFunc(func(e types.T) {
			state = accumulator(state, e)
		}),
	))
	return state
}

//Count Returns the count of elements in this Stream.
func (s Stream) Count() int {
	return s.ReduceFromIdentity(0, func(e1 types.T, e2 types.T) types.T {
		return e1.(int) + 1
	}).(int)
}

//Max Compare through the compare function, return the max value in Stream.
func (s Stream) Max(compare func(first types.T, second types.T) int) types.T {
	return s.Reduce(func(e1 types.T, e2 types.T) types.T {
		if compare(e1, e2) >= 0 {
			return e1
		}
		return e2
	})
}

//Min Compare through the compare function, return the min value in Stream.
func (s Stream) Min(compare func(first types.T, second types.T) int) types.T {
	return s.Reduce(func(e1 types.T, e2 types.T) types.T {
		if compare(e1, e2) <= 0 {
			return e1
		}
		return e2
	})
}

//ToSlice Returns a Slice Containing all elements of the Stream.
func (s Stream) ToSlice() []types.T {
	pipeline := s.p
	var resultSlice []types.T
	pipeline.evaluate(newDefaultTerminalStage(
		beginFunc(func(size int) {
			if size != -1 {
				resultSlice = make([]types.T, 0, size)
				return
			}
			resultSlice = make([]types.T, 0)
		}),
		acceptFunc(func(e types.T) {
			resultSlice = append(resultSlice, e)
		}),
	))
	return resultSlice
}

//ToMap Returns a Map containing all elements of the Stream transformed by the keyMapper function.
func (s Stream) ToMap(keyMapper func(t types.T) types.K, valueMapper func(t types.T) types.R) map[types.K]types.R {
	pipeline := s.p
	var resultMap map[types.K]types.R
	pipeline.evaluate(newDefaultTerminalStage(
		beginFunc(func(size int) {
			if size != -1 {
				resultMap = make(map[types.K]types.R, size)
				return
			}
			resultMap = make(map[types.K]types.R)
		}),
		acceptFunc(func(e types.T) {
			key := keyMapper(e)
			value := valueMapper(e)
			resultMap[key] = value
		}),
	))
	return resultMap
}

//GroupingBy Returns a Map containing all elements of the Stream transformed by the classifier function.
func (s Stream) GroupingBy(classifier func(t types.T) types.K) map[types.K][]types.T {
	pipeline := s.p
	var resultGroupingMap map[types.K][]types.T
	pipeline.evaluate(newDefaultTerminalStage(
		beginFunc(func(size int) {
			if size != -1 {
				resultGroupingMap = make(map[types.K][]types.T, size)
				return
			}
			resultGroupingMap = make(map[types.K][]types.T)
		}),
		acceptFunc(func(e types.T) {
			groupingKey := classifier(e)
			if resultGroupingMap[groupingKey] == nil {
				resultGroupingMap[groupingKey] = make([]types.T, 0)
			}
			resultGroupingMap[groupingKey] = append(resultGroupingMap[groupingKey], e)
		}),
	))
	return resultGroupingMap
}

//Terminal short-circuiting operation

//AllMatch Returns whether all elements of this Stream match the predicate function.
func (s Stream) AllMatch(predicate func(e types.T) bool) bool {
	pipeline := s.p
	var matchRes = true
	pipeline.evaluate(newDefaultTerminalStage(
		acceptFunc(func(e types.T) {
			if !predicate(e) {
				matchRes = false
			}
		}), cancellationRequestedFunc(func() bool {
			return !matchRes
		})))

	return matchRes
}

//AnyMatch Returns whether any elements of this Stream match the predicate function.
func (s Stream) AnyMatch(predicate func(e types.T) bool) bool {
	pipeline := s.p
	var matchRes = false
	pipeline.evaluate(newDefaultTerminalStage(
		acceptFunc(func(e types.T) {
			if predicate(e) {
				matchRes = true
			}
		}), cancellationRequestedFunc(func() bool {
			return matchRes
		})))

	return matchRes
}

//NoneMatch Returns whether any elements of this Stream non match the predicate function.
func (s Stream) NoneMatch(predicate func(e types.T) bool) bool {
	pipeline := s.p
	var matchRes = true
	pipeline.evaluate(newDefaultTerminalStage(
		acceptFunc(func(e types.T) {
			if predicate(e) {
				matchRes = false
			}
		}),
		cancellationRequestedFunc(func() bool {
			return !matchRes
		}),
	))

	return matchRes
}

//FindFirst Return the first element of the Stream.
func (s Stream) FindFirst() types.T {
	pipeline := s.p
	var result types.T
	var find = false
	pipeline.evaluate(newDefaultTerminalStage(
		acceptFunc(func(e types.T) {
			if !find {
				result = e
				find = true
			}
		}), cancellationRequestedFunc(func() bool {
			return find
		}),
	))
	return result
}

//Parallel operation

//Parallel Set the number of workers to perform Stream operations in parallel.
func (s Stream) Parallel(workers int) Stream {
	pipeline := s.p
	pipeline.workers = workers
	return s
}
