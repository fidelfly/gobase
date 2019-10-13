package httprxr

import (
	"net/http"

	"github.com/fidelfly/gox/errorx"
)

//export
func ResponseJSON(w http.ResponseWriter, statusCode int, data interface{}, settings ...RespSetting) {
	var setting ResponseSetting = JSONResponse
	if len(settings) > 0 {
		setting = settings[0]
	}

	Response(w, setting, NewJSONFiller(data), statusCode)
}

//export
func Response(w http.ResponseWriter, setting ResponseSetting, filler ResponseFiller, statusCode ...int) {
	w.Header().Set("Content-Type", setting.GetContentType())
	if !setting.IsCacheEnable() {
		w.Header().Set("cacheEnable-Control", "no-store")
		w.Header().Set("pragma", "no-cache")
	}

	header := setting.GetHeader()
	if header != nil {
		for key := range header {
			w.Header().Set(key, header.Get(key))
		}
	}

	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}

	w.WriteHeader(status)
	err := filler.FillResponse(&RespWriter{w})
	if err != nil {
		_, _ = w.Write([]byte("Error found during filling response"))
	}
}

//export
func NewErrorMessage(code, message string, data ...map[string]interface{}) ResponseMessage {
	return NewResponseMessage(RespMessageType.Error, code, message, data...)
}

//export
func MakeErrorMessage(code string, err error, data ...map[string]interface{}) ResponseMessage {
	return NewResponseMessage(RespMessageType.Error, code, err.Error(), data...)
}

//export
func ErrorMessage(err errorx.Error, data ...map[string]interface{}) ResponseMessage {
	return NewResponseMessage(RespMessageType.Error, err.Code(), err.Message(), data...)
}

//export
func ExceptionMessage(err error, codes ...string) ResponseMessage {
	code := http.StatusText(http.StatusInternalServerError)
	if len(codes) > 0 {
		code = codes[0]
	}
	return NewResponseMessage(RespMessageType.Error, code, err.Error())
}
