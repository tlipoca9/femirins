package logx

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cockroachdb/errors"
)

func ConsoleErrorStackHandlerFunc() HandlerFunc {
	return func(ctx context.Context, record slog.Record) (context.Context, slog.Record) {
		newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
		record.Attrs(func(attr slog.Attr) bool {
			if err, ok := attr.Value.Any().(error); ok {
				errString := fmt.Sprintf("%+v", err)
				newRecord.AddAttrs(slog.String(
					attr.Key,
					strings.ReplaceAll(errString, "\t", "  ")),
				)
				return true
			}
			newRecord.AddAttrs(attr)
			return true
		})
		return ctx, newRecord
	}
}

func JSONErrorStackHandlerFunc() HandlerFunc {
	return func(ctx context.Context, record slog.Record) (context.Context, slog.Record) {
		newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
		record.Attrs(func(attr slog.Attr) bool {
			if err, ok := attr.Value.Any().(error); ok {
				data, _ := errors.BuildSentryReport(err)
				newRecord.AddAttrs(slog.Group(
					attr.Key,
					slog.String("message", err.Error()),
					slog.Any("exception", data.Exception),
				))
				return true
			}
			newRecord.AddAttrs(attr)
			return true
		})
		return ctx, newRecord
	}
}
