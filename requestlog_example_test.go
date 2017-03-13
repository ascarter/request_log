package requestlog_test

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

func Example() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/goodbye", goodbyeHandler)

	// Use request log middleware with default logger
	http.Handle("/", requestlog.RequestLogDefaultHandler(mux))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
