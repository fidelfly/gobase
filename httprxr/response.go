package httprxr

import (
	"encoding/json"
	"net/http"
)

type ResponseData map[string]interface{}

type ResponseSetting interface {
	GetContentType() string
	IsCacheEnable() bool
	GetHeader() http.Header
}

type RespSetting struct {
	contentType string
	cacheEnable bool
	header      http.Header
}

//export
func NewRespSetting(contentType string, cacheEnable bool, header http.Header) *RespSetting {
	return &RespSetting{
		contentType,
		cacheEnable,
		header,
	}
}

func (rs RespSetting) GetContentType() string {
	return rs.contentType
}

func (rs RespSetting) IsCacheEnable() bool {
	return rs.cacheEnable
}

func (rs *RespSetting) SetCacheEnable(cache bool) {
	rs.cacheEnable = cache
}

func (rs RespSetting) GetHeader() http.Header {
	return rs.header
}

func (rs *RespSetting) SetHeader(header http.Header) {
	rs.header = header
}

func (rs RespSetting) New(cacheEnable bool, header http.Header) *RespSetting {
	rs.SetCacheEnable(cacheEnable)
	rs.SetHeader(header)
	return &rs
}

type RespWriter struct {
	w http.ResponseWriter
}

func (rw *RespWriter) Write(p []byte) (n int, err error) {
	n, err = rw.w.Write(p)
	return
}

type ResponseFiller interface {
	FillResponse(w *RespWriter) error
}

type messageType string
type responseMsgType struct {
	Error   messageType
	Info    messageType
	Warning messageType
}

var RespMessageType = responseMsgType{
	Error:   messageType("error"),
	Info:    messageType("info"),
	Warning: messageType("warning"),
}

type ResponseMessage struct {
	Code    string                 `json:"error_code"`
	Message string                 `json:"error_message"`
	Data    map[string]interface{} `json:"data"`
	Type    messageType            `json:"type"`
}

//export
func NewResponseMessage(msgType messageType, code, message string, data ...map[string]interface{}) ResponseMessage {
	return ResponseMessage{
		Type:    msgType,
		Code:    code,
		Message: message,
		Data:    combineData(data...),
	}
}

func combineData(data ...map[string]interface{}) map[string]interface{} {
	switch len(data) {
	case 0:
		return nil
	case 1:
		return data[0]
	default:
		mapData := make(map[string]interface{})
		for _, m := range data {
			for key, value := range m {
				mapData[key] = value
			}
		}
		return mapData
	}

}

var JSONResponse = &RespSetting{
	contentType: "application/json;charset=UTF-8",
	cacheEnable: false,
}

type JSONFiller struct {
	data interface{}
}

func (jf *JSONFiller) FillResponse(w *RespWriter) (err error) {
	if jf.data != nil {
		if byteData, ok := jf.data.([]byte); ok {
			_, err = w.Write(byteData)
		} else {
			err = json.NewEncoder(w).Encode(jf.data)
		}
	}
	return
}

func NewJSONFiller(data interface{}) *JSONFiller {
	return &JSONFiller{data}
}
