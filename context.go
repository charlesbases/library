package library

import (
	"context"

	"github.com/charlesbases/library/sonyflake"
)

// HeaderTraceID traceid in context and header
const HeaderTraceID = "X-Trace-ID"

// ContextValueTraceID get tarceid with context
func ContextValueTraceID(ctx context.Context) sonyflake.ID {
	if ctx == context.Background() || ctx == nil {
		return 0
	}
	id, _ := ctx.Value(HeaderTraceID).(sonyflake.ID)
	return id
}
