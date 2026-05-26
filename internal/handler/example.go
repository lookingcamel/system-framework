package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/lookingcamel/system-framework/internal/tracing"
)

type ExampleHandler struct{}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

func (h *ExampleHandler) GetExample(c *gin.Context) {
	ctx, span := startSpanSafely(c.Request.Context(), "GetExample")
	defer span.End()

	var req struct {
		ID string `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		Error(c, 400, "Invalid request")
		return
	}

	Success(c, gin.H{
		"id":      req.ID,
		"message": "Hello from system-framework",
	})
	c.Request = c.Request.WithContext(ctx)
}

func (h *ExampleHandler) ListExample(c *gin.Context) {
	ctx, span := startSpanSafely(c.Request.Context(), "ListExample")
	defer span.End()

	Success(c, gin.H{
		"items": []string{"item1", "item2", "item3"},
	})
	c.Request = c.Request.WithContext(ctx)
}

func startSpanSafely(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if tracing.Tracer() != nil {
		return tracing.StartSpan(ctx, name, opts...)
	}
	// Fallback to no-op tracer
	tr := otel.Tracer("fallback")
	return tr.Start(ctx, name, opts...)
}
