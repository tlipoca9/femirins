package logx

import (
	"context"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

var (
	DefaultConsoleHandler = New(
		log.NewWithOptions(
			os.Stdout,
			log.Options{
				Level:           log.DebugLevel,
				ReportTimestamp: true,
				ReportCaller:    true,
			},
		),
		ContextHandlerFunc(),
		ConsoleErrorStackHandlerFunc(),
	)

	DefaultJSONHandler = New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
		ContextHandlerFunc(),
		JSONErrorStackHandlerFunc(),
	)
)

type HandlerFunc func(ctx context.Context, record slog.Record) (context.Context, slog.Record)

var _ slog.Handler = (*RecordHandler)(nil)

type RecordHandler struct {
	chains    []HandlerFunc
	delegator slog.Handler
}

func New(h slog.Handler, chains ...HandlerFunc) slog.Handler {
	return RecordHandler{
		chains:    chains,
		delegator: h,
	}
}

func (m RecordHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return m.delegator.Enabled(ctx, level)
}

func (m RecordHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, fn := range m.chains {
		ctx, record = fn(ctx, record)
	}
	return m.delegator.Handle(ctx, record)
}

func (m RecordHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return RecordHandler{
		chains:    m.chains,
		delegator: m.delegator.WithAttrs(attrs),
	}
}

func (m RecordHandler) WithGroup(name string) slog.Handler {
	return RecordHandler{
		chains:    m.chains,
		delegator: m.delegator.WithGroup(name),
	}
}
