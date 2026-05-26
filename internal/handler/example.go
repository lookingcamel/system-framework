package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/lookingcamel/system-framework/internal/tracing"
)

type ExampleHandler struct{}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

func (h *ExampleHandler) GetExample(c *gin.Context) {
	ctx, span := tracing.StartSpan(c.Request.Context(), "GetExample")
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
	ctx, span := tracing.StartSpan(c.Request.Context(), "ListExample")
	defer span.End()

	Success(c, gin.H{
		"items": []string{"item1", "item2", "item3"},
	})
	c.Request = c.Request.WithContext(ctx)
}
