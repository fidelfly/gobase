package gosrvx

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"gopkg.in/oauth2.v3"

	"github.com/fidelfly/gox"
	"github.com/fidelfly/gox/pkg/randx"
	"github.com/fidelfly/gox/routex"

	"github.com/gorilla/mux"

	"github.com/fidelfly/gox/errorx"
	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/logx"

	"github.com/fidelfly/gox/authx"
)

type tokenKey struct {
}
type userKey struct {
}

func GetUserKey(r *http.Request) string {
	if v := r.Context().Value(userKey{}); v != nil {
		if key, ok := v.(string); ok {
			return key
		}
	}
	return ""
}

type requestIdKey struct {
}

func GetRequestId(r *http.Request) string {
	if v := r.Context().Value(requestIdKey{}); v != nil {
		if key, ok := v.(string); ok {
			return key
		}
	}
	panic("there is no request id if audit for router is not enabled")
}

const (
	//contextKey
	//ContextUserKey   = "context.user.id"
	//ContextTokenKey  = "context.token"
	//ContextRequestId = "context.request.id"

	//routerName
	defaultRouterKey = "router.default"
)

var routerMap = make(map[string]*RootRouter)

//var defaultRouter = NewRouter()

type RouterHook func()

type RootRouter struct {
	*routex.Router
	//authServer  *authx.Server
	auditLogger logx.StdLog
	authFilter  func(w http.ResponseWriter, req *http.Request, next http.Handler)
}

type RouterPlugin interface {
	Inject(*RootRouter)
}

func (rr *RootRouter) AttachPlugins(plugins ...RouterPlugin) {
	for _, plugin := range plugins {
		plugin.Inject(rr)
	}
}

type TokenIssuer struct {
	*authx.Server
	tokenPath string
	//clearPath string
}

func (t *TokenIssuer) Setup(server *authx.Server, tokenPath string) {
	t.Server = server
	t.tokenPath = tokenPath
}

//export
func NewTokenIssuer(server *authx.Server, tokenPath string) *TokenIssuer {
	return &TokenIssuer{server, tokenPath}
}

func (t *TokenIssuer) Inject(rr *RootRouter) {
	//rr.SetAuthFilter(t.AuthFilter)
	rr.EnableAuthFilter(t.AuthFilter)
	rr.Path(t.tokenPath).Methods(http.MethodPost).HandlerFunc(t.HandleTokenRequest)
}

func (t *TokenIssuer) AuthFilter(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if ti, err := t.ValidateToken(w, r); err != nil {
		if codeError, ok := err.(errorx.Error); ok {
			httprxr.ResponseJSON(w, http.StatusUnauthorized, httprxr.ErrorMessage(codeError))
		} else {
			httprxr.ResponseJSON(w, http.StatusUnauthorized, httprxr.MakeErrorMessage(authx.UnauthorizedErrorCode, err))
		}
		return
	} else if ti != nil {
		r = httprxr.ContextSet(r, userKey{}, ti.GetUserID(), tokenKey{}, ti)
	}
	next.ServeHTTP(w, r)

}

func (t *TokenIssuer) AuthorizeDisposeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		obj := httprxr.ContextGet(r, userKey{})
		if obj != nil {
			if ti, ok := obj.(oauth2.TokenInfo); ok {
				logx.CaptureError(t.RemoveAccessToken(ti.GetAccess()))
				logx.CaptureError(t.RemoveRefreshToken(ti.GetRefresh()))
			}
		}
	})
}

func (t *TokenIssuer) AuthorizeDisposeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	obj := httprxr.ContextGet(r, tokenKey{})
	if obj != nil {
		if ti, ok := obj.(oauth2.TokenInfo); ok {
			logx.CaptureError(t.RemoveAccessToken(ti.GetAccess()))
			logx.CaptureError(t.RemoveRefreshToken(ti.GetRefresh()))
			httprxr.ResponseJSON(w, http.StatusOK, nil)
			return
		}
		//should never come to here
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(errors.New("token is not right")))
		return
	}
	httprxr.ResponseJSON(w, http.StatusNotFound, nil)
}

func (rr *RootRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer gox.CapturePanicAndRecover(fmt.Sprintf("Panic found and recovered during %s:%s", req.Method, req.URL.Path))
	rr.Router.ServeHTTP(w, req)
}

func (rr *RootRouter) EnableAudit(loggers ...logx.StdLog) {
	if len(loggers) > 0 {
		rr.auditLogger = loggers[0]
	} else {
		rr.auditLogger = gox.ConsoleOutput{}
	}
	rr.Router.Use(rr.AuditMiddleware)
}

func (rr *RootRouter) EnableRecover() {
	rr.Router.Use(rr.RecoverMiddleware)
}

func (rr *RootRouter) SetAuthFilter(filter func(w http.ResponseWriter, req *http.Request, next http.Handler)) {
	rr.authFilter = filter
}

func (rr *RootRouter) EnableAuthFilter(filter ...func(w http.ResponseWriter, req *http.Request, next http.Handler)) {
	if len(filter) > 0 {
		rr.SetAuthFilter(filter[0])
	}
	rr.Router.Use(rr.AuthorizeMiddleware)
}

func (rr *RootRouter) ProtectPrefix(pathPrefix string) *routex.Router {
	myRouter := rr.PathPrefix(pathPrefix).Restricted(true).Subrouter()
	//myRouter.Use(rr.AuthorizeMiddleware)
	return myRouter
}

/*func (rr *RootRouter) SetAuthorizer(server *authx.Server) {
	rr.authServer = server
}*/

func (rr *RootRouter) CurrentRouteConfig(r *http.Request) (routex.RouteConfig, bool) {
	if route := mux.CurrentRoute(r); route != nil {
		config := rr.GetRouteConfig(route)
		if config != nil {
			return config.GetCopy(true), true
		}
	}
	//should never come to here
	logx.Errorf("can't find route confx for %s : %s", r.Method, r.URL.Path)
	return routex.NewConfig(), false
}

func (rr *RootRouter) AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rr.authFilter == nil {
			next.ServeHTTP(w, r)
			return
		}
		restricted := false
		if config, ok := rr.CurrentRouteConfig(r); ok {
			restricted = config.IsRestricted()
		}
		if restricted {
			rr.authFilter(w, r, next)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (rr *RootRouter) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				if rr.auditLogger != nil {
					rr.auditLogger.Panicf("Panic occurs when handle %s %s", r.Method, r.URL.Path)
					rr.auditLogger.Error(debug.Stack())
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (rr *RootRouter) AuditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		audit := false
		if route := mux.CurrentRoute(r); route != nil {
			config := rr.GetRouteConfig(route)
			if config != nil {
				audit = config.IsAuditEnable()
			}
		}
		auditStart := time.Now()
		w = httprxr.MakeStatusResponse(w)
		reqId := randx.GenUUID(r.URL.Path)
		r = httprxr.ContextSet(r, requestIdKey{}, reqId)
		next.ServeHTTP(w, r)

		auditEnd := time.Now()
		if audit {
			statusCode := 0
			if sr, ok := w.(*httprxr.StatusResponse); ok {
				statusCode = sr.GetStatusCode()
			}
			duration := auditEnd.Sub(auditStart) / time.Millisecond
			user := httprxr.ContextGet(r, userKey{})
			if user != nil {
				rr.auditLogger.Infof("[Router Audit] %s %s [Duration=%dms, User=%s, Status=%s, RequestId=%s]",
					r.Method, r.URL.Path, duration, user, http.StatusText(statusCode), reqId)
			} else {
				rr.auditLogger.Infof("[Router Audit] %s %s [Duration=%dms, Status=%s, RequestId=%s]",
					r.Method, r.URL.Path, duration, http.StatusText(statusCode), reqId)
			}
		}
	})
}

var routerHooks = make([]RouterHook, 0)

//export
func AddRouterHook(hook RouterHook) {
	routerHooks = append(routerHooks, hook)
}

//export
func AttachHookRoute() {
	for _, hook := range routerHooks {
		hook()
	}
	routerHooks = nil
}

//export
func Router(routerKey ...string) *RootRouter {
	//var myRouter *RootRouter
	if len(routerKey) > 0 {
		for _, key := range routerKey {
			if myRouter, ok := routerMap[key]; ok {
				return myRouter
			}
		}
		return nil
	}

	if myRouter, ok := routerMap[defaultRouterKey]; ok {
		return myRouter
	}
	myRouter := &RootRouter{
		Router: routex.New(),
	}
	routerMap[defaultRouterKey] = myRouter
	return myRouter

}

//export
func NewRouter(routerKey ...string) *RootRouter {
	myRouter := &RootRouter{
		Router: routex.New(),
	}
	if len(routerKey) > 0 {
		for _, key := range routerKey {
			if len(key) > 0 {
				routerMap[key] = myRouter
			}
		}
	}

	return myRouter
}

//export
func NewAuditRouter(logger logx.StdLog, routerKey ...string) *RootRouter {
	myRouter := NewRouter(routerKey...)
	myRouter.EnableAudit(logger)
	return myRouter
}

//export
func GetRouter(routerKey string) *RootRouter {
	return routerMap[routerKey]
}

//export
func EnableRouterAudit(loggers ...logx.StdLog) *RootRouter {
	myRouter := Router()
	if len(loggers) > 0 {
		myRouter.EnableAudit(loggers...)
	} else {
		myRouter.EnableAudit(logx.StandardLogger())
	}
	return myRouter
}

//export
func AttachRouterPlugin(plugins ...RouterPlugin) {
	Router().AttachPlugins(plugins...)
}

//export
func ProtectPrefix(pathPrefix string) *routex.Router {
	return Router().ProtectPrefix(pathPrefix)
}

//export
func AttachMiddleware(handler http.Handler, middlewares ...func(handler http.Handler) http.Handler) http.Handler {
	if len(middlewares) == 0 {
		return handler
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

//export
func AttachFuncMiddleware(handlerFunc http.HandlerFunc, middlewares ...func(handler http.Handler) http.Handler) http.Handler {
	return AttachMiddleware(handlerFunc, middlewares...)
}
