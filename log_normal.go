//go:build !(js && wasm)

package scaff

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

const IsWASM = false

func newBaseHandler() slog.Handler {
	return tint.NewHandler(colorable.NewColorable(os.Stdout), &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})
}
