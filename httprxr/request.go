package httprxr

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fidelfly/fxgo/logx"

	"github.com/gorilla/mux"
)

//export
func GetRequestVars(r *http.Request, keys ...string) map[string]string {
	vars := make(map[string]string, len(keys))
	muxVars := mux.Vars(r)
	logx.CaptureError(r.ParseForm())
	for _, key := range keys {
		value := muxVars[key]
		if len(value) == 0 {
			value = r.FormValue(key)
		}
		vars[key] = value
	}
	return vars
}

//export
func GetJSONRequestData(r *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	if isJSONRequest(r) {
		bodyData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return data
		}
		logx.CaptureError(json.Unmarshal(bodyData, &data))
	}
	return data
}

func isJSONRequest(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	if len(contentType) > 0 {
		return strings.Contains(strings.ToLower(contentType), "application/json")
	}
	return false
}
