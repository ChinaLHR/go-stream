package stream

import (
	"github.com/chinalhr/go-stream/types"
	"reflect"
)

//iterator Is a source data iterator.
//GetSize Returns the size of the iterator source data.
//HasNext Returns true if the iterator has more elements.
//Next Returns the next element in the iterator.
type iterator interface {
	GetSize() int
	HasNext() bool
	Next() types.T
}

//iteratorBaseInfo The basic information of an iterator.
type iteratorBaseInfo struct {
	currentIndex int
	size         int
}

func (i *iteratorBaseInfo) GetSize() int {
	return i.size
}

func (i *iteratorBaseInfo) HasNext() bool {
	return i.currentIndex < i.size
}

type iteratorInfiniteBaseInfo struct{}

func (i *iteratorInfiniteBaseInfo) GetSize() int {
	return -1
}

func (i *iteratorInfiniteBaseInfo) HasNext() bool {
	return true
}

//sliceIterator A general type of slice iterator.
type sliceIterator struct {
	*iteratorBaseInfo
	elements []types.T
}

func (iterator *sliceIterator) Next() types.T {
	element := iterator.elements[iterator.currentIndex]
	iterator.currentIndex++
	return element
}

//sliceReflectIterator A general type slice iterator based on reflect.
type sliceReflectIterator struct {
	*iteratorBaseInfo
	sliceValue reflect.Value
}

func (iterator *sliceReflectIterator) Next() types.T {
	element := iterator.sliceValue.
		Index(iterator.currentIndex).
		Interface()
	iterator.currentIndex++
	return element
}

//mapReflectIterator A general type map iterator based on reflect.
type mapReflectIterator struct {
	*iteratorBaseInfo
	mapValue *reflect.MapIter
}

func (iterator *mapReflectIterator) Next() types.T {
	iterator.currentIndex++
	iterator.mapValue.Next()
	return types.KV{
		KEY:   iterator.mapValue.Key().Interface(),
		VALUE: iterator.mapValue.Value().Interface(),
	}
}

//supplierIterator A general type supplier iterator based on supplier function.
type supplierIterator struct {
	*iteratorInfiniteBaseInfo
	get func() types.T
}

func (iterator *supplierIterator) Next() types.T {
	return iterator.get()
}

//build iterator

func buildSliceIterator(elements ...types.T) iterator {
	return &sliceIterator{
		iteratorBaseInfo: &iteratorBaseInfo{
			currentIndex: 0,
			size:         len(elements),
		},
		elements: elements,
	}
}

func buildSliceReflectIterator(value reflect.Value) iterator {
	return &sliceReflectIterator{
		iteratorBaseInfo: &iteratorBaseInfo{
			currentIndex: 0,
			size:         value.Len(),
		},
		sliceValue: value,
	}
}

func buildMapReflectIterator(value reflect.Value) iterator {
	return &mapReflectIterator{
		iteratorBaseInfo: &iteratorBaseInfo{
			currentIndex: 0,
			size:         value.Len(),
		},
		mapValue: value.MapRange(),
	}
}

func buildSupplierIterator(fn func() types.T) iterator {
	return &supplierIterator{
		iteratorInfiniteBaseInfo: &iteratorInfiniteBaseInfo{},
		get:                      fn,
	}
}
