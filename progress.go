package fxgo

import (
	"net/http"

	"github.com/fidelfly/fxgo/logx"

	"github.com/fidelfly/fxgo/pkg/randx"

	"github.com/sirupsen/logrus"

	"github.com/fidelfly/fxgo/cachex"
	"github.com/fidelfly/fxgo/httprxr"
)

var socketCache = cachex.CreateCache(cachex.DefaultExpiration, 0)

//export
func GetProgress(key string, code string) *httprxr.WsProgress {
	if conn, ok := socketCache.Get(key); ok {
		return httprxr.NewWsProgress(conn.(*httprxr.WsConnect), code)
	}
	return httprxr.NewWsProgress(nil, code)
}

//export
func SetupProgressRoute(wsPath string, restricted bool) {
	AttchProgressRoute(Router(), wsPath, restricted)
}

//export
func AttchProgressRoute(router *RootRouter, wsPath string, restricted bool) {
	router.HandleFunc(wsPath, ProgressSetupHandler).Restricted(restricted)
}

func ProgressSetupHandler(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "code")
	code := params["code"]

	wsc := &httprxr.WsConnect{Code: code, Duration: 100}

	err := httprxr.SetupWebsocket(wsc, w, r)
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}

	progressKey := randx.GenUUID(code)

	socketCache.Set(progressKey, wsc)
	defer socketCache.Remove(progressKey)

	logx.CaptureError(wsc.Conn.WriteJSON(map[string]string{"progressKey": progressKey}))

	wsc.ListenAndServe()

	logrus.Infof("WebSocket %s is Closed", r.RequestURI)
}
