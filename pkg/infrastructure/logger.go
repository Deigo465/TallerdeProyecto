package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

const (
	black   = "\033[1;30m"
	red     = "\033[1;31m"
	green   = "\033[1;32m"
	yellow  = "\033[1;33m"
	blue    = "\033[1;34m"
	magenta = "\033[1;35m"
	cyan    = "\033[1;36m"
	white   = "\033[1;37m"
	reset   = "\033[0m"
)

func ColorMsg(msg any, color string) string {
	return fmt.Sprintf("%s%v%s", color, msg, reset)
}

type color struct {
}

func (c *color) CyanString(s string) string {
	return ColorMsg(s, cyan)
}

func (c *color) MagentaString(s string) string {
	return ColorMsg(s, magenta)
}
func (c *color) BlueString(s string) string {
	return ColorMsg(s, blue)
}
func (c *color) YellowString(s string) string {
	return ColorMsg(s, yellow)
}
func (c *color) RedString(s string) string {
	return ColorMsg(s, red)
}
func (c *color) WhiteString(s string) string {
	return ColorMsg(s, white)
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"
	color := color{}

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

func NewSLogger(name string) *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)
	logger.With("app", name)
	return logger
}
