package session

import (
	"context"
)

type ctxKey int

const (
	sessionCtxKey ctxKey = 0
	dbTxCtxKey    ctxKey = iota
)

// SetSessionContextKey set the session in the given context object
// and returns new context with session
func SetSessionContextKey(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, sessionCtxKey, s)
}

// GetSessionFromContext extracts session from the given context
// with flag indicating whether session found or not
func GetSessionFromContext(ctx context.Context) (*Session, bool) {
	sess, ok := ctx.Value(sessionCtxKey).(*Session)
	return sess, ok
}
