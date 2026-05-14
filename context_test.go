package way

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewContext(t *testing.T) {
	// Create a mock DB and Session
	mockDB := &DB{}
	mockSession := &Session{}

	// Create a mock http.ResponseWriter and http.Request
	mockResponseWriter := &MockResponseWriter{}
	mockRequest := &http.Request{}

	// Call the NewContext function
	context := NewContext(mockResponseWriter, mockRequest, mockDB, mockSession, nil)

	// Check if the ResponseWriter, Request, db, and Session fields are set correctly
	if context.Response != mockResponseWriter {
		t.Errorf("NewContext() Response field = %v; want %v", context.Response, mockResponseWriter)
	}
	if context.Request != mockRequest {
		t.Errorf("NewContext() Request field = %v; want %v", context.Request, mockRequest)
	}
	if context.db != mockDB {
		t.Errorf("NewContext() db field = %v; want %v", context.db, mockDB)
	}
	if context.Session != mockSession {
		t.Errorf("NewContext() Session field = %v; want %v", context.Session, mockSession)
	}
	if context.Logger == nil {
		t.Error("NewContext() Logger is nil")
	}
	if context.HTTPClient == nil {
		t.Error("NewContext() HTTPClient is nil")
	}
}

// MockResponseWriter is a mock implementation of http.ResponseWriter for testing purposes
type MockResponseWriter struct {
}

func (m *MockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *MockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *MockResponseWriter) WriteHeader(int) {
}

func TestContextJSONEncoderFailureWritesSingleError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := NewContext(rec, req, nil, nil, nil)

	ctx.JSON(http.StatusCreated, make(chan int))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType == "application/json" {
		t.Fatalf("content type = %q, should not be JSON after encode failure", contentType)
	}
}

func TestContextXMLFormatterFailureWritesSingleError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := NewContext(rec, req, nil, nil, nil)

	ctx.XML(http.StatusCreated, map[string]string{"bad": "xml"})

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestContextProxyMediaUsesConfiguredHTTPClient(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := NewContext(rec, req, nil, nil, nil)
	ctx.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://example.test/media" {
			t.Fatalf("URL = %q, want configured media URL", req.URL.String())
		}
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Header:     http.Header{"Content-Type": []string{"image/png"}},
			Body:       io.NopCloser(strings.NewReader("media")),
		}, nil
	})}

	ctx.ProxyMedia("https://example.test/media")

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if got := rec.Body.String(); got != "media" {
		t.Fatalf("body = %q, want media", got)
	}
}

func TestContextProxyMediaHandlesClientError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := NewContext(rec, req, nil, nil, nil)
	ctx.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("blocked")
	})}

	ctx.ProxyMedia("https://example.test/media")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
