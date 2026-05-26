package utils

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetLocalIP(t *testing.T) {
	ip := GetLocalIP()
	if ip == "" {
		t.Error("GetLocalIP should return a non-empty string")
	}
}

func TestGetClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1")

	ip := GetClientIP(req)
	if ip != "192.168.1.1" {
		t.Errorf("Expected '192.168.1.1', got '%s'", ip)
	}
}

func TestGetClientIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Real-IP", "10.0.0.2")

	ip := GetClientIP(req)
	if ip != "10.0.0.2" {
		t.Errorf("Expected '10.0.0.2', got '%s'", ip)
	}
}

func TestGetClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	ip := GetClientIP(req)
	if ip == "" {
		t.Error("IP should not be empty")
	}
}

func TestGenerateID(t *testing.T) {
	id := GenerateID()
	if id == "" {
		t.Error("GenerateID should return a non-empty string")
	}
	if len(id) != 32 {
		t.Errorf("Expected ID length 32, got %d", len(id))
	}
}

func TestGenerateID_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := GenerateID()
		if ids[id] {
			t.Error("GenerateID should generate unique IDs")
		}
		ids[id] = true
	}
}

func TestFormatDuration_Milliseconds(t *testing.T) {
	d := 500 * time.Millisecond
	result := FormatDuration(d)
	if result != "500ms" {
		t.Errorf("Expected '500ms', got '%s'", result)
	}
}

func TestFormatDuration_Seconds(t *testing.T) {
	d := 2 * time.Second
	result := FormatDuration(d)
	if result != "2.00s" {
		t.Errorf("Expected '2.00s', got '%s'", result)
	}
}

func TestFormatDuration_Minutes(t *testing.T) {
	d := 90 * time.Second
	result := FormatDuration(d)
	if result != "1.50m" {
		t.Errorf("Expected '1.50m', got '%s'", result)
	}
}

func TestFormatDuration_Hours(t *testing.T) {
	d := 2 * time.Hour
	result := FormatDuration(d)
	if result != "2.00h" {
		t.Errorf("Expected '2.00h', got '%s'", result)
	}
}

func TestContains_True(t *testing.T) {
	slice := []string{"a", "b", "c"}
	result := Contains(slice, "b")
	if !result {
		t.Error("Contains should return true")
	}
}

func TestContains_False(t *testing.T) {
	slice := []string{"a", "b", "c"}
	result := Contains(slice, "d")
	if result {
		t.Error("Contains should return false")
	}
}

func TestContains_EmptySlice(t *testing.T) {
	slice := []string{}
	result := Contains(slice, "a")
	if result {
		t.Error("Contains should return false for empty slice")
	}
}

func TestInSlice(t *testing.T) {
	slice := []string{"x", "y", "z"}
	result := InSlice(slice, "y")
	if !result {
		t.Error("InSlice should return true")
	}
}

func TestRemoveDuplicates(t *testing.T) {
	slice := []string{"a", "b", "a", "c", "b", "d"}
	result := RemoveDuplicates(slice)

	if len(result) != 4 {
		t.Errorf("Expected 4 unique elements, got %d", len(result))
	}

	expected := map[string]bool{"a": true, "b": true, "c": true, "d": true}
	for _, v := range result {
		if !expected[v] {
			t.Errorf("Unexpected element: %s", v)
		}
	}
}

func TestRemoveDuplicates_EmptySlice(t *testing.T) {
	slice := []string{}
	result := RemoveDuplicates(slice)
	if len(result) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(result))
	}
}

func TestRemoveDuplicates_NoDuplicates(t *testing.T) {
	slice := []string{"a", "b", "c"}
	result := RemoveDuplicates(slice)
	if len(result) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(result))
	}
}
