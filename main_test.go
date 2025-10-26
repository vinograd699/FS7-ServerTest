package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCafeNegative checks invalid requests and error responses
func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	tests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}

	for _, tc := range tests {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", tc.request, nil)
		handler.ServeHTTP(resp, req)

		assert.Equal(t, tc.status, resp.Code)                              // Validate HTTP status code
		assert.Equal(t, tc.message, strings.TrimSpace(resp.Body.String())) // Validate error message
	}
}

// TestCafeWhenOk checks valid requests and success responses
func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	tests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}

	for _, tc := range tests {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", tc, nil)
		handler.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code) // Validate successful response
	}
}

// TestCafeCount validates the 'count' parameter behavior
func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	tests := []struct {
		count int // Input count value
		want  int // Expected number of cafes in response
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])}, // cafeList must be defined in the code
	}

	for _, tc := range tests {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?count="+strconv.Itoa(tc.count)+"&city=moscow", nil)
		handler.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code) // Fatal check for HTTP status

		body := strings.TrimSpace(resp.Body.String())
		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",") // Assumes response is comma-separated list
		}
		assert.Equal(t, tc.want, len(cafes)) // Validate number of results
	}
}

// TestCafeSearch validates the 'search' parameter behavior
func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	tests := []struct {
		search    string // Search query
		wantCount int    // Expected number of matching results
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
		{"ко", 3},
	}

	for _, tc := range tests {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+tc.search, nil)
		handler.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code) // Fatal check for HTTP status

		body := strings.TrimSpace(resp.Body.String())
		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",")
		}
		assert.Equal(t, tc.wantCount, len(cafes)) // Validate number of results

		// Ensure each result contains the search term
		for _, cafe := range cafes {
			assert.True(t, strings.Contains(strings.ToLower(cafe), strings.ToLower(tc.search)))
		}
	}
}
