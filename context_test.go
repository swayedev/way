package way

import (
	"net/http"
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
	context := NewContext(mockResponseWriter, mockRequest, mockDB, mockSession)

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
