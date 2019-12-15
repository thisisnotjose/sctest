package processors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/thisisnotjose/sctest/internal/types"
)

func TestNewDeadLettersProcessor(t *testing.T) {
	ctx := resetContext()

	tests := []struct {
		name         string
		events       []types.Event
		wantErr      bool
		emptyResults bool
	}{
		{
			"TestFollow",
			[]types.Event{
				types.Event{
					Sequence:       1,
					ReceiverUserID: 1,
					EmitterUserID:  2,
					EventType:      "F",
					Payload:        "1|F|1|2\n",
				},
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := []string{}
			p := deadlettersProcessor{
				ctx: &ctx,
				printDL: func(v string) {
					results = append(results, v)
				},
			}

			for _, evt := range tt.events {
				ctx.DeadLetterChannel <- evt.Payload

				go p.Start()
				time.Sleep(300 * time.Millisecond)
				assert.Equal(t, evt.Payload, results[0])
				close(ctx.DeadLetterChannel)
				close(ctx.EventChannel)
			}

			ctx = resetContext()
		})
	}
}
