package logx

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

var AttrsCtxKey ctxKey

func ContextWithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(attrs) == 0 {
		return ctx
	}
	existsAttrs, ok := ctx.Value(AttrsCtxKey).([]slog.Attr)
	if !ok {
		existsAttrs = make([]slog.Attr, 0)
	}
	newAttrs := append(existsAttrs, attrs...)
	return context.WithValue(ctx, AttrsCtxKey, newAttrs)

}

func ContextHandlerFunc() HandlerFunc {
	return func(ctx context.Context, record slog.Record) (context.Context, slog.Record) {
		attrs, ok := ctx.Value(AttrsCtxKey).([]slog.Attr)
		if !ok {
			return ctx, record
		}
		record.AddAttrs(attrs...)
		return ctx, record
	}
}
