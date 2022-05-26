package stream

import "github.com/chinalhr/go-stream/types"

//stage Is a stage on the stream pipeline.used to conduct values through the stages of a stream pipeline,
//with additional methods to manage size information,control flow, etc.
//Begin Informs the stage how much data will flow in. The Begin method needs to be called before the Accept method.
//Accept Processing the data flowing into the stage.
//End After all the data has been sent, you need to call the End method.
//CancellationRequested Decide whether you still need to pass data to stage by calling CancellationRequested.
type stage interface {
	Begin(size int)
	Accept(e types.T)
	End()
	CancellationRequested() bool
}

type chainedStage struct {
	begin                 func(int)
	accept                func(e types.T)
	end                   func()
	cancellationRequested func() bool
}

func (c *chainedStage) Begin(size int) {
	c.begin(size)
}

func (c *chainedStage) Accept(e types.T) {
	c.accept(e)
}

func (c *chainedStage) CancellationRequested() bool {
	return c.cancellationRequested()
}

func (c *chainedStage) End() {
	c.end()
}

//stageFunc Returns the function that operates on stage.
type stageFunc func(c *chainedStage)

func beginFunc(fn func(int)) stageFunc {
	return func(c *chainedStage) {
		c.begin = fn
	}
}

func acceptFunc(fn func(e types.T)) stageFunc {
	return func(c *chainedStage) {
		c.accept = fn
	}
}

func cancellationRequestedFunc(fn func() bool) stageFunc {
	return func(c *chainedStage) {
		c.cancellationRequested = fn
	}
}

func endFunc(fn func()) stageFunc {
	return func(c *chainedStage) {
		c.end = fn
	}
}

//newDefaultIntermediateStage Create a default intermediate stage. The method accepts the next stage,
//and the stage with the default intermediate operation will call the method of the next stage to connect
//the pipelines.
func newDefaultIntermediateStage(nextStage stage, fn ...stageFunc) *chainedStage {
	c := &chainedStage{
		begin:                 nextStage.Begin,
		accept:                nextStage.Accept,
		cancellationRequested: nextStage.CancellationRequested,
		end:                   nextStage.End,
	}
	for _, fun := range fn {
		fun(c)
	}
	return c
}

// newDefaultIntermediateStage Create a default terminal stage.
func newDefaultTerminalStage(fn ...stageFunc) *chainedStage {
	c := &chainedStage{
		begin:                 func(i int) {},
		accept:                func(e types.T) {},
		cancellationRequested: func() bool { return false },
		end:                   func() {},
	}
	for _, fun := range fn {
		fun(c)
	}
	return c
}
