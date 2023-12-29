package way

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewContext(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseWriter mock
	res := httptest.NewRecorder()

	// Create a new Context
	ctx := NewContext(res, req)

	// Test if the Context is correctly initialized
	if ctx.Request != req {
		t.Errorf("Expected Request to be %v, got %v", req, ctx.Request)
	}

	if ctx.Response != res {
		t.Errorf("Expected Response to be %v, got %v", res, ctx.Response)
	}
}
