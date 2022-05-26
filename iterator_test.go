package stream

import (
	"github.com/chinalhr/go-stream/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSliceIterator(t *testing.T) {
	testCases := []func(it iterator){
		func(it iterator) {
			size := it.GetSize()
			assert.Equal(t, 5, size)
		},
		func(it iterator) {
			for i := 1; i <= 5; i++ {
				e := it.Next()
				assert.Equal(t, i, e)
			}
		},
		func(it iterator) {
			hasNext := it.HasNext()
			assert.Equal(t, false, hasNext)
		},
	}

	sliceIterator := buildSliceIterator(1, 2, 3, 4, 5)
	for _, testCase := range testCases {
		t.Run("testCase", func(t *testing.T) {
			testCase(sliceIterator)
		})
	}
}

func TestSliceReflectIterator(t *testing.T) {
	testData := [5]int{1, 2, 3, 4, 5}
	testCases := []func(it iterator){
		func(it iterator) {
			size := it.GetSize()
			assert.Equal(t, 5, size)
		},
		func(it iterator) {
			for i := 1; i <= 5; i++ {
				e := it.Next()
				assert.Equal(t, i, e)
			}
		},
		func(it iterator) {
			hasNext := it.HasNext()
			assert.Equal(t, false, hasNext)
		},
	}

	sliceReflectIterator := buildSliceReflectIterator(reflect.ValueOf(testData))
	for _, testCase := range testCases {
		t.Run("testCase", func(t *testing.T) {
			testCase(sliceReflectIterator)
		})
	}
}

func TestMapReflectIterator(t *testing.T) {
	testData := map[int]string{
		1: "A",
		2: "B",
		3: "C",
		4: "D",
		5: "E",
	}
	testCases := []func(it iterator){
		func(it iterator) {
			size := it.GetSize()
			assert.Equal(t, 5, size)
		},
		func(it iterator) {
			for i := 1; i <= 5; i++ {
				e := it.Next()
				delete(testData, e.(types.KV).KEY.(int))
			}
		},
		func(it iterator) {
			hasNext := it.HasNext()
			assert.Equal(t, false, hasNext)
			assert.Empty(t, testData)
		},
	}
	mapReflectIterator := buildMapReflectIterator(reflect.ValueOf(testData))
	for _, testCase := range testCases {
		t.Run("testCase", func(t *testing.T) {
			testCase(mapReflectIterator)
		})
	}
}

func TestSupplierIterator(t *testing.T) {
	var testData = 0
	testCases := []func(it iterator){
		func(it iterator) {
			size := it.GetSize()
			assert.Equal(t, -1, size)
		},
		func(it iterator) {
			for i := 1; i <= 5; i++ {
				e := it.Next()
				assert.Equal(t, i, e)
			}
		},
		func(it iterator) {
			hasNext := it.HasNext()
			assert.Equal(t, true, hasNext)
		},
	}
	supplierIterator := buildSupplierIterator(func() types.T {
		testData = testData + 1
		return testData
	})
	for _, testCase := range testCases {
		t.Run("testCase", func(t *testing.T) {
			testCase(supplierIterator)
		})
	}
}
