package fxgo

import (
	"errors"
	"net/http"
	"time"

	"github.com/fidelfly/fxgo/routex"

	"github.com/gorilla/mux"

	"gopkg.in/oauth2.v3"

	"github.com/fidelfly/fxgo/errorx"
	"github.com/fidelfly/fxgo/httprxr"
	"github.com/fidelfly/fxgo/logx"

	"github.com/fidelfly/fxgo/authx"
)

const (
	//contextKey
	ContextUserKey  = "context.user.id"
	ContextTokenKey = "context.token"

	//routerName
	defaultRouterKey = "router.default"
)

var routerMap = map[string]*RootRouter{
	defaultRouterKey: {
		Router: routex.New(),
	},
}

//var defaultRouter = NewRouter()

type RouterHook func()

type RootRouter struct {
	*routex.Router
	authServer  *authx.Server
	auditLogger logx.StdLog
}

func (rr *RootRouter) SetLogger(logger logx.StdLog) {
	rr.auditLogger = logger
}

func (rr *RootRouter) EnableAudit(loggers ...logx.StdLog) {
	if len(loggers) > 0 {
		rr.auditLogger = loggers[0]
	} else {
		rr.auditLogger = ConsoleOutput{}
	}
	rr.Router.Use(rr.AuditMiddleware)
}

func (rr *RootRouter) ProtectPrefix(pathPrefix string) *routex.Router {
	myRouter := rr.PathPrefix(pathPrefix).Restricted(true).Subrouter()
	//myRouter.Use(rr.AuthorizeMiddleware)
	return myRouter
}

func (rr *RootRouter) SetAuthorizer(server *authx.Server) {
	rr.authServer = server
}

func (rr *RootRouter) HandleTokenRequest(w http.ResponseWriter, r *http.Request) {
	logx.CaptureError(rr.authServer.HandleTokenRequest(w, r))
}

func (rr *RootRouter) AuthorizeDisposeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		obj := httprxr.ContextGet(r, ContextTokenKey)
		if obj != nil {
			if ti, ok := obj.(oauth2.TokenInfo); ok {
				logx.CaptureError(rr.authServer.RemoveAccessToken(ti.GetAccess()))
				logx.CaptureError(rr.authServer.RemoveRefreshToken(ti.GetRefresh()))
			}
		}
	})
}

func (rr *RootRouter) AuthorizeDisposeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	obj := httprxr.ContextGet(r, ContextTokenKey)
	if obj != nil {
		if ti, ok := obj.(oauth2.TokenInfo); ok {
			logx.CaptureError(rr.authServer.RemoveAccessToken(ti.GetAccess()))
			logx.CaptureError(rr.authServer.RemoveRefreshToken(ti.GetRefresh()))
			httprxr.ResponseJSON(w, http.StatusOK, nil)
			return
		}
		//should never come to here
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(errors.New("token is not right")))
		return
	}
	httprxr.ResponseJSON(w, http.StatusNotFound, nil)
}

func (rr *RootRouter) CurrentRouteConfig(r *http.Request) (routex.RouteConfig, bool) {
	if route := mux.CurrentRoute(r); route != nil {
		config := rr.GetRouteConfig(route)
		if config != nil {
			return config.GetCopy(), true
		}
	}
	//should never come to here
	logx.Errorf("can't find route confx for %s : %s", r.Method, r.URL.Path)
	return routex.NewConfig(), false
}

func (rr *RootRouter) AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rr.authServer == nil {
			next.ServeHTTP(w, r)
			return
		}
		restricted := false
		if config, ok := rr.CurrentRouteConfig(r); ok {
			restricted = config.IsRestricted()
		}
		if restricted {
			if ti, err := rr.authServer.ValidateToken(w, r); err != nil {
				if codeError, ok := err.(errorx.Error); ok {
					httprxr.ResponseJSON(w, http.StatusUnauthorized, httprxr.ErrorMessage(codeError))
				} else {
					httprxr.ResponseJSON(w, http.StatusUnauthorized, httprxr.MakeErrorMessage(authx.UnauthorizedErrorCode, err))
				}
				return
			} else if ti != nil {
				r = httprxr.ContextSet(r, ContextUserKey, ti.GetUserID(), ContextTokenKey, ti)
			}
		}

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
		next.ServeHTTP(w, r)

		auditEnd := time.Now()
		if audit {
			statusCode := 0
			if sr, ok := w.(*httprxr.StatusResponse); ok {
				statusCode = sr.GetStatusCode()
			}
			duration := auditEnd.Sub(auditStart) / time.Millisecond
			user := httprxr.ContextGet(r, ContextUserKey)
			if user != nil {
				rr.auditLogger.Infof("[Router Audit] %s %s [Duration=%dms, User=%s, Status=%s]", r.Method, r.URL.Path, duration, user, http.StatusText(statusCode))
			} else {
				rr.auditLogger.Infof("[Router Audit] %s %s [Duration=%dms, Status=%s]", r.Method, r.URL.Path, duration, http.StatusText(statusCode))
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
	if len(routerKey) == 0 {
		return routerMap[defaultRouterKey]
	}
	for _, key := range routerKey {
		if myRouter, ok := routerMap[key]; ok {
			return myRouter
		}
	}
	return nil

}

//export
func NewRouter(routerKey ...string) *RootRouter {
	myRouter := &RootRouter{
		Router: routex.New(),
	}
	if len(routerKey) > 0 {
		for _, key := range routerKey {
			routerMap[key] = myRouter
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
func InitRouter(audit bool, logger logx.StdLog) {
	if logger != nil {
		Router().SetLogger(logger)
	}
	if audit {
		Router().EnableAudit()
	}
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
func SetupAuthorizeRoute(tokenPath string, authServer *authx.Server) *RootRouter {
	return AttchAuthorizeRoute(Router(), tokenPath, authServer)
}

//export
func AttchAuthorizeRoute(router *RootRouter, tokenPath string, authServer *authx.Server, middlewares ...func(handler http.Handler) http.Handler) *RootRouter {
	router.SetAuthorizer(authServer)
	router.Path(tokenPath).Methods(http.MethodPost).Handler(AttachFuncMiddleware(router.HandleTokenRequest, middlewares...))
	router.Use(router.AuthorizeMiddleware)
	return router
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
	return AttachMiddleware(http.HandlerFunc(handlerFunc), middlewares...)
}
