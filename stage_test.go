package stream

import (
	"github.com/chinalhr/go-stream/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChainedStage(t *testing.T) {
	var terminalStage *chainedStage
	var intermediateStageWrap func(stage) stage
	actualResult := []int{1, 4, 9, 16, 25}
	terminalStageResult := make([]int, 0, 5)

	testCases := []func(){
		func() {
			intermediateStageWrap = func(next stage) stage {
				return newDefaultIntermediateStage(terminalStage,
					beginFunc(func(size int) {
						assert.Equal(t, size, 5)
						next.Begin(size)
					}),
					acceptFunc(func(e types.T) {
						assert.NotNil(t, e)
						next.Accept(e.(int) * e.(int))
					}),
				)
			}
		},

		func() {
			terminalStage = newDefaultTerminalStage(
				beginFunc(func(size int) {
					assert.Equal(t, size, 5)
				}),
				acceptFunc(func(e types.T) {
					assert.NotNil(t, e)
					terminalStageResult = append(terminalStageResult, e.(int))
				}),
				endFunc(func() {
					assert.Equal(t, terminalStageResult, actualResult)
				}),
			)
		},

		func() {
			source := buildSliceIterator(1, 2, 3, 4, 5)
			wrapStage := intermediateStageWrap(terminalStage)
			wrapStage.Begin(source.GetSize())
			for source.HasNext() && !wrapStage.CancellationRequested() {
				wrapStage.Accept(source.Next())
			}
			wrapStage.End()
		},
	}

	for _, testCase := range testCases {
		t.Run("testCase", func(t *testing.T) {
			testCase()
		})
	}
}
