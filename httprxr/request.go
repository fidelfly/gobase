package httprxr

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fidelfly/gox/logx"

	"github.com/gorilla/mux"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

type RequestVar map[string][]string

func (rv RequestVar) Exist(key string) bool {
	_, ok := rv[key]
	return ok
}

func (rv RequestVar) set(key string, vals ...string) {
	rv[key] = vals
}

func (rv RequestVar) Get(key string) []string {
	return rv[key]
}

func (rv RequestVar) GetString(key string) string {
	if val, ok := rv[key]; ok {
		if len(val) > 0 {
			return val[0]
		}
		return ""
	}
	return ""
}

func (rv RequestVar) GetInt(key string) (int64, error) {
	val := rv.GetString(key)
	return strconv.ParseInt(val, 10, 64)
}

func (rv RequestVar) GetUint(key string) (uint64, error) {
	val := rv.GetString(key)
	return strconv.ParseUint(val, 10, 64)
}

func (rv RequestVar) GetFloat(key string) (float64, error) {
	val := rv.GetString(key)
	return strconv.ParseFloat(val, 64)
}

func (rv RequestVar) GetTime(key string, format string, loc ...*time.Location) (time.Time, error) {
	val := rv.GetString(key)
	if len(loc) > 0 {
		return time.ParseInLocation(format, val, loc[0])
	}
	return time.Parse(format, val)
}

func (rv RequestVar) GetBool(key string) (bool, error) {
	val := rv.GetString(key)
	return strconv.ParseBool(val)
}

func (rv RequestVar) GetInts(key string) ([]int64, error) {
	vals := rv.Get(key)
	if len(vals) == 0 {
		return make([]int64, 0), nil
	}
	arr := make([]int64, len(vals))
	for i, val := range vals {
		if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
			arr[i] = intVal
		} else {
			return nil, err
		}
	}
	return arr, nil
}

func (rv RequestVar) GetUints(key string) ([]uint64, error) {
	vals := rv.Get(key)
	if len(vals) == 0 {
		return make([]uint64, 0), nil
	}
	arr := make([]uint64, len(vals))
	for i, val := range vals {
		if uintVal, err := strconv.ParseUint(val, 10, 64); err == nil {
			arr[i] = uintVal
		} else {
			return nil, err
		}
	}
	return arr, nil
}

func (rv RequestVar) GetFloats(key string) ([]float64, error) {
	vals := rv.Get(key)
	if len(vals) == 0 {
		return make([]float64, 0), nil
	}
	arr := make([]float64, len(vals))
	for i, val := range vals {
		if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
			arr[i] = floatVal
		} else {
			return nil, err
		}
	}
	return arr, nil
}

func (rv RequestVar) GetTimes(key string, format string, loc ...*time.Location) ([]time.Time, error) {
	vals := rv.Get(key)
	if len(vals) == 0 {
		return make([]time.Time, 0), nil
	}
	arr := make([]time.Time, len(vals))
	for i, val := range vals {
		var timeVal time.Time
		var err error
		if len(loc) > 0 {
			timeVal, err = time.ParseInLocation(format, val, loc[0])
		} else {
			timeVal, err = time.Parse(format, val)
		}
		if err == nil {
			arr[i] = timeVal
		} else {
			return nil, err
		}
	}
	return arr, nil
}

func (rv RequestVar) GetBools(key string) ([]bool, error) {
	vals := rv.Get(key)
	if len(vals) == 0 {
		return make([]bool, 0), nil
	}
	arr := make([]bool, len(vals))
	for i, val := range vals {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			arr[i] = boolVal
		} else {
			return nil, err
		}
	}
	return arr, nil
}

func ParseRequestVars(r *http.Request, keys ...string) RequestVar {
	vars := make(RequestVar, len(keys))
	muxVars := mux.Vars(r)
	if r.Form == nil {
		_ = r.ParseMultipartForm(defaultMaxMemory)
	}
	if len(keys) > 0 {
		for _, key := range keys {
			value := muxVars[key]
			if len(value) > 0 {
				vars.set(key, value)
			} else if len(r.Form) > 0 {
				if fv, ok := r.Form[key]; ok {
					vars.set(key, fv...)
				}
			}
		}
	} else {
		for key, value := range muxVars {
			vars.set(key, value)
		}
		if len(r.Form) > 0 {
			for key, value := range r.Form {
				vars.set(key, value...)
			}
		}
	}

	return vars
}

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
func GetJSONRequestMap(r *http.Request) map[string]interface{} {
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

//export
func GetJSONRequestData(r *http.Request, data interface{}) error {
	if isJSONRequest(r) {
		bodyData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(bodyData, &data)
	}
	return nil
}

func isJSONRequest(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	if len(contentType) > 0 {
		return strings.Contains(strings.ToLower(contentType), "application/json")
	}
	return false
}
