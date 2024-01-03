package way_test

import (
	"net/http"
	"testing"

	"github.com/swayedev/way"
)

// MockResponseWriter is a mock implementation of http.ResponseWriter
type MockResponseWriter struct {
}

func (m *MockResponseWriter) Header() http.Header {
	return nil
}

func (m *MockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *MockResponseWriter) WriteHeader(int) {
}

func TestNewContext(t *testing.T) {
	// Create a mock DB object
	mockDB := &way.DB{}

	// Create a mock http.ResponseWriter
	mockResponseWriter := &MockResponseWriter{}

	// Create a mock http.Request
	mockRequest := &http.Request{}

	// Call the NewContext function
	context := way.NewContext(mockDB, mockResponseWriter, mockRequest)

	// Check if the ResponseWriter is set correctly
	if context.Response != mockResponseWriter {
		t.Errorf("NewContext() ResponseWriter = %v; want %v", context.Response, mockResponseWriter)
	}

	// Check if the Request is set correctly
	if context.Request != mockRequest {
		t.Errorf("NewContext() Request = %v; want %v", context.Request, mockRequest)
	}

	// Check if the db is set correctly
	if context.GetDB() != mockDB {
		t.Errorf("NewContext() DB = %v; want %v", context.GetDB(), mockDB)
	}
}
