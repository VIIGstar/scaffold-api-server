package app_http

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type NOOPHandler struct {
}

func (h *NOOPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 2)
}

func TestPanicReporterHandler_ServeHTTP(t *testing.T) {
	const timeoutMessage = "Timeout"
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PanicReportTimeoutHandler(&NOOPHandler{}, time.Second, timeoutMessage).ServeHTTP)
	handler.ServeHTTP(rr, req)
	// Check the response body is what we expect.
	assert.Equal(t, timeoutMessage, rr.Body.String(), "Body is not ok")
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
}
