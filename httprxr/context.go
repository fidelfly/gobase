package httprxr

import (
	"context"
	"net/http"
)

func ContextGet(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

func ContextSet(r *http.Request, paris ...interface{}) *http.Request {
	if len(paris) == 0 {
		return r
	}
	ct := r.Context()
	for i := 0; i < len(paris)-1; i += 2 {
		key := paris[i]
		val := paris[i+1]
		if key != nil && val != nil {
			ct = context.WithValue(ct, key, val)
		}
	}

	return r.WithContext(ct)
}
