package scaff

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

type prefixHandler struct {
	h      slog.Handler
	prefix string
}

func (p *prefixHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return p.h.Enabled(ctx, level)
}

const blue = "\x1b[38;5;67m" // slightly brighter muted blue
const reset = "\x1b[0m"

func (p *prefixHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Message = blue + p.prefix + reset + r.Message
	return p.h.Handle(ctx, r)
}

func (p *prefixHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &prefixHandler{h: p.h.WithAttrs(attrs), prefix: p.prefix}
}

func (p *prefixHandler) WithGroup(name string) slog.Handler {
	return &prefixHandler{h: p.h.WithGroup(name), prefix: p.prefix}
}

// TODO: Expand to be able to set to JSON logging in production for potentially connecting loki and stuff
func NewLogger(part string) *slog.Logger {
	h := tint.NewHandler(colorable.NewColorable(os.Stdout), &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})

	return slog.New(&prefixHandler{
		h:      h,
		prefix: part + " > ",
	})
}
