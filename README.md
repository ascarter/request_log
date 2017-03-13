# request_log
Request logging middleware for Go

# requestlog [![GoDoc](https://godoc.org/github.com/ascarter/requestlog?status.svg)](http://godoc.org/github.com/ascarter/requestlog)[![Go Report Card](https://goreportcard.com/badge/github.com/ascarter/requestlog)](https://goreportcard.com/report/github.com/ascarter/requestlog)

Request logging middleware for Go. Requestlog wraps any handler with start and end information. The end line will include the total time for the request. If the header `X-Request-ID` is present, the request id will also be logged.

`RequestLogHandler` uses a `log.Logger` for output. If the default `log.Logger` is in use, calling `RequestLogDefaultHandler` will use the same configuration as the default logger.

# Usage Example

```go

package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/ascarter/requestlog"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Running hello handler")
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Running goodbye handler")
	fmt.Fprintf(w, "Goodbye, %q", html.EscapeString(r.URL.Path))
}

func main() {
	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/goodbye", goodbyeHandler)

	// Use request log middleware with default logger
	http.Handle("/", requestlog.RequestLogDefaultHandler(mux))

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

```

# Output Examples

Log output with `X-Request-ID` set:

```
2017/03/13 14:20:57 [dc6efe7f-cfe7-418c-baa3-7c0f80334572] Started GET /goodbye for [192.168.0.10]:62966
2017/03/13 14:20:57 Running goodbye handler
2017/03/13 14:20:57 [dc6efe7f-cfe7-418c-baa3-7c0f80334572] Completed 200 OK in 237.884µs
```