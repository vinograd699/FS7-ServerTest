package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCafeNegative checks invalid /cafe requests and error responses
func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()                // Create mock response recorder
		req := httptest.NewRequest("GET", v.request, nil) // Create mock request
		handler.ServeHTTP(response, req)                  // Execute handler

		assert.Equal(t, v.status, response.Code)                              // Check HTTP status
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String())) // Check error message
	}
}

// TestCafeWhenOk checks valid /cafe requests and success responses
func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()        // Create mock response recorder
		req := httptest.NewRequest("GET", v, nil) // Create mock request
		handler.ServeHTTP(response, req)          // Execute handler

		assert.Equal(t, http.StatusOK, response.Code) // Verify success status
	}
}
