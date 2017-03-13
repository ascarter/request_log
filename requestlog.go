package requestlog

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (r *responseLogger) Header() http.Header {
	return r.w.Header()
}

func (r *responseLogger) Write(b []byte) (int, error) {
	if r.status == 0 {
		// Status will be StatusOK if WriteHeader not called yet
		r.status = http.StatusOK
	}
	size, err := r.w.Write(b)
	r.size += size
	return size, err
}

func (r *responseLogger) WriteHeader(s int) {
	r.w.WriteHeader(s)
	r.status = s
}

func (r *responseLogger) Status() int {
	return r.status
}

func (r *responseLogger) Size() int {
	return r.size
}

// RequestLogHandler logs request to output logger with start and end information.
// The start log line includes information on the request.
// The end log line includes result of request and time elapsed.
// If `X-Request-ID` header is present, includes the request id.
//
// Example log line:
//	2017/03/13 14:20:57 [dc6efe7f-cfe7-418c-baa3-7c0f80334572] Started GET /goodbye for [192.168.0.10]:62966
//	2017/03/13 14:20:57 Running goodbye handler
//	2017/03/13 14:20:57 [dc6efe7f-cfe7-418c-baa3-7c0f80334572] Completed 200 OK in 237.884µs

func RequestLogHandler(h http.Handler, out *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseLogger{w: w}

		raddr := r.Header.Get("X-Forwarded-For")
		if raddr == "" {
			raddr = r.RemoteAddr
		}

		rid := r.Header.Get("X-Request-ID")
		if rid != "" {
			rid = fmt.Sprintf("[%s] ", rid)
		}

		out.Printf("%sStarted %s %s for %s", rid, r.Method, r.URL.Path, raddr)
		h.ServeHTTP(rw, r)
		out.Printf("%sCompleted %v %s in %v\n", rid, rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	})
}

// RequestLogDefaultHandler logs requests to the standard logger.
func RequestLogDefaultHandler(h http.Handler) http.Handler {
	// Standard logger configuration
	logger := log.New(os.Stderr, "", log.LstdFlags)
	return RequestLogHandler(h, logger)
}
