package gosrvx

import (
	"fmt"
	"log"
	"net/http"
)

//export
func ListenAndServe(handler http.Handler, port int64) {
	server := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", port),
	}

	log.Fatal(server.ListenAndServe())
}

//export
func ListenAndServeTLS(certificate string, key string, handler http.Handler, port int64) {
	server := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", port),
	}

	log.Fatal(server.ListenAndServeTLS(certificate, key))
}
