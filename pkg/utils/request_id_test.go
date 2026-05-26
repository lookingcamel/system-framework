package utils

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGenerateRequestID(t *testing.T) {
	id := GenerateRequestID()
	if id == "" {
		t.Error("GenerateRequestID should return a non-empty string")
	}
}

func TestGenerateRequestID_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := GenerateRequestID()
		if ids[id] {
			t.Error("GenerateRequestID should generate unique IDs")
		}
		ids[id] = true
	}
}

func TestSetRequestID(t *testing.T) {
	ctx := context.Background()
	newCtx := SetRequestID(ctx, "test-request-id")

	requestID := GetRequestID(newCtx)
	if requestID != "test-request-id" {
		t.Errorf("Expected 'test-request-id', got '%s'", requestID)
	}
}

func TestGetRequestID_NotSet(t *testing.T) {
	ctx := context.Background()
	requestID := GetRequestID(ctx)
	if requestID != "" {
		t.Errorf("Expected empty string, got '%s'", requestID)
	}
}

func TestGinSetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	GinSetRequestID(c, "gin-request-id")

	requestID := GinGetRequestID(c)
	if requestID != "gin-request-id" {
		t.Errorf("Expected 'gin-request-id', got '%s'", requestID)
	}
}

func TestGinGetRequestID_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestID := GinGetRequestID(c)
	if requestID != "" {
		t.Errorf("Expected empty string, got '%s'", requestID)
	}
}
