package global

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

var LOG *slog.Logger

func SetupLogger() {
	if CFG.Debug {
		LOG = slog.New(&TraceHandler{
			slog.NewTextHandler(
				os.Stderr,
				&slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug},
			),
		})
	} else {
		LOG = slog.New(&TraceHandler{
			slog.NewJSONHandler(
				os.Stderr,
				&slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo},
			),
		})
	}
}

type TraceHandler struct {
	slog.Handler
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx).SpanContext()
	if span.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.TraceID().String()),
			slog.String("span_id", span.SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, r)
}
