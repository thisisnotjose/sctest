package processors

import (
	"fmt"

	"github.com/thisisnot/sctest/internal/types"
)

type deadlettersProcessor struct {
	ctx     *types.Context
	printDL func(string)
}

// NewDeadLettersProcessor returns a new events processor
func NewDeadLettersProcessor(ctx *types.Context) types.Processor {
	return deadlettersProcessor{
		ctx: ctx,
		printDL: func(v string) {
			fmt.Printf("Dead letter event %v", v)
		},
	}
}

// Start kicks off the listener on the events channel and for every message it
// adds it to the queue and then processes as many messages as it can
func (p deadlettersProcessor) Start() {
	for newDeadLetter := range p.ctx.DeadLetterChannel {
		p.printDL(newDeadLetter)
	}
}
