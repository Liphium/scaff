package scaff

import (
	"context"
	"log/slog"
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
	if IsWASM {
		// Avoid ANSI color codes in WASM. The browser writer
		// will parse this plain prefix and apply native CSS styles instead.
		r.Message = p.prefix + r.Message
	} else {
		r.Message = blue + p.prefix + reset + r.Message
	}
	return p.h.Handle(ctx, r)
}

func (p *prefixHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &prefixHandler{h: p.h.WithAttrs(attrs), prefix: p.prefix}
}

func (p *prefixHandler) WithGroup(name string) slog.Handler {
	return &prefixHandler{h: p.h.WithGroup(name), prefix: p.prefix}
}

// NewLogger automatically yields a native browser console logger
// when compiled for WASM, and a tinted terminal logger otherwise.
func NewLogger(part string) *slog.Logger {
	h := newBaseHandler()

	return slog.New(&prefixHandler{
		h:      h,
		prefix: part + " > ",
	})
}
