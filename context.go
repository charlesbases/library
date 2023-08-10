package library

import (
	"context"
)

// HeaderTraceID traceid in context and header
const HeaderTraceID = "X-Trace-ID"

// ContextValueTraceID get tarceid with context
func ContextValueTraceID(ctx context.Context) string {
	if ctx == context.Background() || ctx == nil {
		return ""
	}
	id, _ := ctx.Value(HeaderTraceID).(string)
	return id
}
