package stream

import (
	"github.com/chinalhr/go-stream/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeline(t *testing.T) {
	testCases := []func(p *referencePipeline){
		func(p *referencePipeline) {
			assert.NotNil(t, p.it)
			assert.NotNil(t, p.currentOpt)
			assert.Nil(t, p.currentOpt.preOpt)
			assert.Nil(t, p.currentOpt.wrapStage)
		},

		func(p *referencePipeline) {
			p.addOperation(func(next stage) stage {
				flag := false
				return newDefaultIntermediateStage(next,
					beginFunc(func(size int) {
						flag = true
						assert.Equal(t, size, 5)
						next.Begin(size)
					}),
					acceptFunc(func(e types.T) {
						next.Accept(e.(int) * e.(int))
					}),
					endFunc(func() {
						assert.Equal(t, flag, true)
						next.End()
					}),
				)
			})
			assert.NotNil(t, p.currentOpt)
			assert.NotNil(t, p.currentOpt.preOpt)
			assert.NotNil(t, p.currentOpt.wrapStage)
		},

		func(p *referencePipeline) {
			p.addOperation(func(next stage) stage {
				return newDefaultIntermediateStage(next,
					beginFunc(func(size int) {
						assert.Equal(t, size, 5)
						next.Begin(-1)
					}),
					acceptFunc(func(e types.T) {
						if e.(int) != 4 {
							next.Accept(e)
						}
					}))
			})
		},

		func(p *referencePipeline) {
			actualResult := []int{1, 9, 16, 25}
			evalResult := make([]int, 0, 4)
			p.evaluate(newDefaultTerminalStage(
				beginFunc(func(size int) {
					assert.Equal(t, size, -1)
				}),
				acceptFunc(func(e types.T) {
					evalResult = append(evalResult, e.(int))
				}),
				endFunc(func() {
					assert.Equal(t, evalResult, actualResult)
				}),
			))
		},
	}

	p := newPipeline(buildSliceIterator(1, 2, 3, 4, 5))
	for _, testCase := range testCases {
		t.Run("testCase", func(t *testing.T) {
			testCase(p)
		})
	}

}
