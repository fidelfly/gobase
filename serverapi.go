package fxgo

import (
	"fmt"
	"log"
	"net/http"
)

//export
func ListenAndServe(router *RootRouter, port int64) {
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", port),
	}

	log.Fatal(server.ListenAndServe())
}

//export
func ListenAndServeTLS(certificate string, key string, router *RootRouter, port int64) {
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", port),
	}

	log.Fatal(server.ListenAndServeTLS(certificate, key))
}
