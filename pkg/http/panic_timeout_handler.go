package app_http

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"runtime"
	"time"
)

// PanicReportTimeoutHandler Replace http.TimeoutHandler with PanicReportTimeoutHandler
func PanicReportTimeoutHandler(h http.Handler, dt time.Duration, msg string) http.Handler {
	return http.TimeoutHandler(&panicReporterHandler{handler: h}, dt, msg)
}

type panicReporterHandler struct {
	handler http.Handler
}

func (h *panicReporterHandler) logf(r *http.Request, format string, args ...interface{}) {
	s, _ := r.Context().Value(http.ServerContextKey).(*http.Server)
	if s != nil && s.ErrorLog != nil {
		s.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (h *panicReporterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			if err != http.ErrAbortHandler {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				logrus.Error(fmt.Sprintf("http: panic serving from %v: \ndetail: %s", r.RemoteAddr, err))
				logrus.Error(fmt.Sprintf("http: panic serving, stacktrace: %s", buf))
			}
			panic(http.ErrAbortHandler)
		}
	}()
	h.handler.ServeHTTP(w, r)
}
