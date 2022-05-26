package stream

import (
	"github.com/chinalhr/go-stream/types"
	"sync"
)

//operation Represents an operation on the pipeline.
//wrapStage stage wrapper function, pass parameters between successive stages by passing next stage.
//preOpt Reference to the previous operation.
type operation struct {
	wrapStage func(stage) stage
	preOpt    *operation
}

//referencePipeline
//currentOpt is the latest intermediate operation in the pipeline.
//workers is the number of parallel executions.
type referencePipeline struct {
	it         iterator
	currentOpt *operation
	workers    int
}

func newPipeline(source iterator) *referencePipeline {
	return &referencePipeline{
		it: source,
		//head has non operation.
		currentOpt: &operation{},
	}
}

func (p *referencePipeline) addOperation(wrap func(stage) stage) {
	op := &operation{
		wrapStage: wrap,
	}
	if p.currentOpt == nil {
		p.currentOpt = op
		return
	}
	op.preOpt = p.currentOpt
	p.currentOpt = op
}

//evaluate the pipeline with a terminal operation to produce a result.
//By passing terminalStage, the stage chain is constructed based on the pipeline based wrapStage method。
//If the number of workers of Pipeline is greater than 1 and the size of pipeline Iterator is greater than 1,
//will be parallel evaluate, otherwise it will be sequential evaluate.
func (p *referencePipeline) evaluate(terminalStage stage) {
	if p.workers > 1 && p.it.GetSize() > 1 {
		p.evaluateParallel(terminalStage)
	} else {
		p.evaluateSequential(terminalStage)
	}
}

func (p *referencePipeline) evaluateSequential(c stage) {
	stage := c
	for i := p.currentOpt; i.preOpt != nil; i = i.preOpt {
		stage = i.wrapStage(stage)
	}
	source := p.it
	stage.Begin(source.GetSize())
	for source.HasNext() && !stage.CancellationRequested() {
		stage.Accept(source.Next())
	}
	stage.End()
}

func (p *referencePipeline) evaluateParallel(c stage) {
	stage := c
	for i := p.currentOpt; i.preOpt != nil; i = i.preOpt {
		stage = i.wrapStage(stage)
	}
	source := p.it

	sharding := sourceSharding(p.it.GetSize(), p.workers)
	var wg sync.WaitGroup
	stage.Begin(source.GetSize())
	wg.Add(len(sharding))
	for _, s := range sharding {
		shardingBucket := s
		sources := make([]types.T, 0, s)
		for shardingBucket != 0 && source.HasNext() {
			sources = append(sources, source.Next())
			shardingBucket--
		}
		go parallelRun(sources, stage, wg.Done)
	}
	wg.Wait()
	stage.End()
}

//sourceSharding Distribute the source data evenly to the worker.
//forExample:
//	input: sourceSize = 22 workerSize = 4
//	return: []int{6,6,5,5}
func sourceSharding(sourceSize int, workerSize int) []int {
	if workerSize > sourceSize {
		workerSize = sourceSize
	}

	sharding := make([]int, workerSize)
	shardingIdx := 0
	for i := 0; i < sourceSize; i++ {
		if shardingIdx < workerSize {
			sharding[shardingIdx]++
			shardingIdx++
		} else {
			sharding[0]++
			shardingIdx = 1
		}
	}
	return sharding
}

func parallelRun(sources []types.T, stage stage, done func()) {
	defer done()

	for i := 0; i < len(sources) && !stage.CancellationRequested(); i++ {
		stage.Accept(sources[i])
	}
}
