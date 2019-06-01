package httprxr

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

type StatusResponse struct {
	http.ResponseWriter
	statusCode int
}

func (sr *StatusResponse) WriteHeader(statusCode int) {
	sr.ResponseWriter.WriteHeader(statusCode)
	sr.statusCode = statusCode
}

func (sr *StatusResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := sr.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}

	return nil, nil, errors.New("not hijacker response")
}

func (sr *StatusResponse) GetStatusCode() int {
	return sr.statusCode
}

func MakeStatusResponse(writer http.ResponseWriter) *StatusResponse {
	return &StatusResponse{writer, http.StatusProcessing}
}
