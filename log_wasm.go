//go:build js && wasm

package scaff

import (
	"encoding/json"
	"log/slog"
	"strings"
	"syscall/js"
)

const IsWASM = true

// jsConsoleWriter converts slog JSON payloads to formatted console API calls.
type jsConsoleWriter struct{}

func (w *jsConsoleWriter) Write(p []byte) (n int, err error) {
	var data map[string]any
	if err := json.Unmarshal(p, &data); err != nil {
		// Fallback to basic logging if JSON is somehow invalid
		js.Global().Get("console").Call("log", string(p))
		return len(p), nil
	}

	msg, _ := data["msg"].(string)
	levelStr, _ := data["level"].(string)

	// Clean up fields that the browser handles or we already processed
	delete(data, "msg")
	delete(data, "level")
	delete(data, "time") // Browser developer consoles already timestamp logs

	// Route to correct console API
	var method string
	switch levelStr {
	case "DEBUG":
		method = "debug"
	case "INFO":
		method = "info"
	case "WARN":
		method = "warn"
	case "ERROR":
		method = "error"
	default:
		method = "log"
	}

	// Prepare arguments for console.Call
	var jsArgs []any

	// Extract the prefix to style it with CSS in the browser console
	if idx := strings.Index(msg, " > "); idx != -1 {
		prefix := msg[:idx]
		actualMsg := msg[idx+3:]
		jsArgs = append(jsArgs,
			"%c"+prefix+" > %c"+actualMsg,
			"color: #5588aa; font-weight: bold;", // Muted blue style for prefix
			"color: inherit;",                    // Default style for remaining message
		)
	} else {
		jsArgs = append(jsArgs, msg)
	}

	// If there are contextual structured attributes (including nested groups),
	// pass them as a native JavaScript object so developers can expand/inspect them.
	if len(data) > 0 {
		jsArgs = append(jsArgs, js.ValueOf(data))
	}

	console := js.Global().Get("console")
	if console.Truthy() {
		console.Call(method, jsArgs...)
	}

	return len(p), nil
}

func newBaseHandler() slog.Handler {
	return slog.NewJSONHandler(&jsConsoleWriter{}, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
}
