package main

import (
	"log/slog"

	"github.com/cockroachdb/errors"

	"github.com/tlipoca9/femirins/logx"
)

func foo() error {
	return errors.New("What Happen !?")
}

func bar() error {
	return errors.Wrapf(foo(), "call foo")
}

func baz() error {
	return errors.Wrapf(bar(), "call bar")
}

func main() {
	slog.SetDefault(slog.New(logx.DefaultConsoleHandler))
	slog.Debug("Hello, World!")
	slog.Info("Hello, World!")
	log := slog.With("answer", 42)
	log.Warn("Hello, World!")
	log.Error("Hello, World!", "err", baz())

	slog.SetDefault(slog.New(logx.DefaultJSONHandler))
	slog.Debug("Hello, World!")
	slog.Info("Hello, World!")
	log = slog.With("answer", 42)
	log.Warn("Hello, World!")
	log.Error("Hello, World!", "err", baz())
}
