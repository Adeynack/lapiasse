package ctxval

import "context"

// FallbackValues returns a new context that looks up values first in the main
// context (ctx), and if not found, looks them up in another one (otherContext).
func FallbackValues(ctx context.Context, otherContext context.Context) context.Context {
	return &joinedValuesContext{
		Context:      ctx,
		otherContext: otherContext,
	}
}

type joinedValuesContext struct {
	context.Context
	otherContext context.Context
}

func (j *joinedValuesContext) Value(key any) any {
	// Try to find the value in the base context first.
	if val := j.Context.Value(key); val != nil {
		return val
	}

	// When not found, look in the other context.
	if val := j.otherContext.Value(key); val != nil {
		return val
	}

	return nil
}
