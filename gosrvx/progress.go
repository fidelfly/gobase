package gosrvx

import (
	"net/http"

	"github.com/fidelfly/gox/cachex/mcache"
	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/pkg/randx"
	"github.com/fidelfly/gox/progx"

	"github.com/sirupsen/logrus"
)

var socketCache = mcache.NewCache(mcache.DefaultExpiration, 0)

//export
func GetProgress(key string, code string) *progx.Progress {
	if conn, ok := socketCache.Get(key); ok {
		if wsconn, ok := conn.(*httprxr.WsConnect); ok {
			return progx.NewProgress((*httprxr.WsProgressHandler)(wsconn), code)
		}
		//return progx.NewProgress(*httprxr.WsProgressHandler(conn.(*httprxr.WsConnect)), code)
	}
	return progx.NewProgress(nil, code)
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
