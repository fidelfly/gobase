package routex

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	myRouter    *mux.Router
	routeConfig map[*mux.Route]*RouteConfig
	config      RouteConfig
}

type RouteRegister func(r *Router)

func (r *Router) Register(regs ...RouteRegister) {
	for _, reg := range regs {
		reg(r)
	}
}

func (r *Router) GetRouteConfig(route *mux.Route) *RouteConfig {
	return r.routeConfig[route]
}

func (r *Router) Restricted(restricted bool) {
	r.config.restricted = restricted
}

func (r *Router) Audit(audit bool) {
	r.config.audit = audit
}

type Route struct {
	myRoute     *mux.Route
	routeConfig map[*mux.Route]*RouteConfig
}

func New() *Router {
	return &Router{myRouter: mux.NewRouter(), routeConfig: make(map[*mux.Route]*RouteConfig), config: RouteConfig{audit: true}}
}

func (r *Route) SetConfig(config RouteConfig) *Route {
	r.Restricted(config.restricted)
	r.Audit(config.audit)
	return r
}

func (r *Route) getConfig() *RouteConfig {
	return r.routeConfig[r.myRoute]
}

type PropSetter func(props *RouteProps)

func (r *Route) ApplyProps(props ...PropSetter) {
	config := r.getConfig()
	if config.props == nil {
		config.props = make(map[string]interface{})
	}
	rp := RouteProps{config.props}
	for _, prop := range props {
		prop(&rp)
	}
}

func Props(key string, val interface{}) PropSetter {
	return func(props *RouteProps) {
		_ = props.Set(key, val)
	}
}

// Deprecated: use ApplyProps and routex.Props(key, val) instead
func (r *Route) SetProps(key string, prop interface{}) *Route {
	if config := r.getConfig(); config != nil {
		config.SetProps(key, prop)
	}
	return r
}

func (r *Route) GetProps(key string) interface{} {
	if config := r.getConfig(); config != nil {
		return config.GetProps(key)
	}
	return nil
}

func (r *Route) Restricted(restricted bool) *Route {
	if config := r.getConfig(); config != nil {
		config.restricted = restricted
	}
	return r
}

func (r *Route) Audit(audit bool) *Route {
	if config := r.getConfig(); config != nil {
		config.audit = audit
	}
	return r
}
func (r *Router) Get(name string) *Route {
	return r.makeRoute(r.myRouter.Get(name))
}

func (r *Router) makeRoute(route *mux.Route) *Route {
	newRoute := &Route{myRoute: route, routeConfig: r.routeConfig}
	config := r.config.GetCopy()
	r.routeConfig[route] = &config
	return newRoute
}

func (r *Router) GetRoute(name string) *Route {
	return r.makeRoute(r.myRouter.GetRoute(name))
}

func (r *Router) StrictSlash(value bool) *Router {
	r.myRouter.StrictSlash(value)
	return r
}

func (r *Router) SkipClean(value bool) *Router {
	r.myRouter.SkipClean(value)
	return r
}

func (r *Router) UseEncodedPath() *Router {
	r.myRouter.UseEncodedPath()
	return r
}

// ----------------------------------------------------------------------------
// Route factories
// ----------------------------------------------------------------------------

// NewRoute registers an empty route.
func (r *Router) NewRoute() *Route {
	return r.makeRoute(r.myRouter.NewRoute())
}

// Name registers a new route with a name.
// See Route.Name().
func (r *Router) Name(name string) *Route {
	return r.makeRoute(r.myRouter.Name(name))
}

// Handle registers a new route with a matcher for the URL path.
// See Route.Path() and Route.Handler().
func (r *Router) Handle(path string, handler http.Handler) *Route {
	return r.makeRoute(r.myRouter.Handle(path, handler))
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.Path() and Route.HandlerFunc().
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.makeRoute(r.myRouter.HandleFunc(path, f))
}

// Headers registers a new route with a matcher for request header values.
// See Route.Headers().
func (r *Router) Headers(pairs ...string) *Route {
	return r.makeRoute(r.myRouter.Headers(pairs...))
}

// Host registers a new route with a matcher for the URL host.
// See Route.Host().
func (r *Router) Host(tpl string) *Route {
	return r.makeRoute(r.myRouter.Host(tpl))
}

// Methods registers a new route with a matcher for HTTP methods.
// See Route.Methods().
func (r *Router) Methods(methods ...string) *Route {
	return r.makeRoute(r.myRouter.Methods(methods...))
}

// Path registers a new route with a matcher for the URL path.
// See Route.Path().
func (r *Router) Path(tpl string) *Route {
	return r.makeRoute(r.myRouter.Path(tpl))
}

// PathPrefix registers a new route with a matcher for the URL path prefix.
// See Route.PathPrefix().
func (r *Router) PathPrefix(tpl string) *Route {
	return r.makeRoute(r.myRouter.PathPrefix(tpl))
}

// Queries registers a new route with a matcher for URL query values.
// See Route.Queries().
func (r *Router) Queries(pairs ...string) *Route {
	return r.makeRoute(r.myRouter.Queries(pairs...))
}

// Schemes registers a new route with a matcher for URL schemes.
// See Route.Schemes().
func (r *Router) Schemes(schemes ...string) *Route {
	return r.makeRoute(r.myRouter.Schemes(schemes...))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.myRouter.ServeHTTP(w, req)
}

func (r *Router) Use(mwf ...func(http.Handler) http.Handler) {
	middlewares := make([]mux.MiddlewareFunc, len(mwf))
	for i := 0; i < len(mwf); i++ {
		middlewares[i] = mwf[i]
	}
	r.myRouter.Use(middlewares...)
}

// SkipClean reports whether path cleaning is enabled for this route via
// Router.SkipClean.
func (r *Route) SkipClean() bool {
	return r.myRoute.SkipClean()
}

func (r *Route) GetConfig() *RouteConfig {
	return r.routeConfig[r.myRoute]
}

// ----------------------------------------------------------------------------
// Route attributes
// ----------------------------------------------------------------------------

// GetError returns an error resulted from building the route, if any.
func (r *Route) GetError() error {
	return r.myRoute.GetError()
}

// Handler --------------------------------------------------------------------

// Handler sets a handler for the route.
func (r *Route) Handler(handler http.Handler) *Route {
	r.myRoute.Handler(handler)
	return r
}

// HandlerFunc sets a handler function for the route.
func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
	return r.Handler(http.HandlerFunc(f))
}

// GetHandler returns the handler for the route, if any.
func (r *Route) GetHandler() http.Handler {
	return r.myRoute.GetHandler()
}

// Name -----------------------------------------------------------------------

// Name sets the name for the route, used to build URLs.
// It is an error to call Name more than once on a route.
func (r *Route) Name(name string) *Route {
	r.myRoute.Name(name)
	return r
}

// GetName returns the name for the route, if any.
func (r *Route) GetName() string {
	return r.myRoute.GetName()
}

func (r *Route) Headers(pairs ...string) *Route {
	r.myRoute.Headers(pairs...)
	return r
}

// HeadersRegexp accepts a sequence of key/value pairs, where the value has regex
// support. For example:
//
//     r := mux.NewRouter()
//     r.HeadersRegexp("Content-Type", "application/(text|json)",
//               "X-Requested-With", "XMLHttpRequest")
//
// The above route will only match if both the request header matches both regular expressions.
// If the value is an empty string, it will match any value if the key is set.
// Use the start and end of string anchors (^ and $) to match an exact value.
func (r *Route) HeadersRegexp(pairs ...string) *Route {
	r.myRoute.HeadersRegexp(pairs...)
	return r
}

// Host -----------------------------------------------------------------------

// Host adds a matcher for the URL host.
// It accepts a template with zero or more URL variables enclosed by {}.
// Variables can define an optional regexp pattern to be matched:
//
// - {name} matches anything until the next dot.
//
// - {name:pattern} matches the given regexp pattern.
//
// For example:
//
//     r := mux.NewRouter()
//     r.Host("www.example.com")
//     r.Host("{subdomain}.domain.com")
//     r.Host("{subdomain:[a-z]+}.domain.com")
//
// Variable names must be unique in a given route. They can be retrieved
// calling mux.Vars(request).
func (r *Route) Host(tpl string) *Route {
	r.myRoute.Host(tpl)
	return r
}

// Methods adds a matcher for HTTP methods.
// It accepts a sequence of one or more methods to be matched, e.g.:
// "GET", "POST", "PUT".
func (r *Route) Methods(methods ...string) *Route {
	r.myRoute.Methods(methods...)
	return r
}

// Path -----------------------------------------------------------------------

// Path adds a matcher for the URL path.
// It accepts a template with zero or more URL variables enclosed by {}. The
// template must start with a "/".
// Variables can define an optional regexp pattern to be matched:
//
// - {name} matches anything until the next slash.
//
// - {name:pattern} matches the given regexp pattern.
//
// For example:
//
//     r := mux.NewRouter()
//     r.Path("/products/").Handler(ProductsHandler)
//     r.Path("/products/{key}").Handler(ProductsHandler)
//     r.Path("/articles/{category}/{id:[0-9]+}").
//       Handler(ArticleHandler)
//
// Variable names must be unique in a given route. They can be retrieved
// calling mux.Vars(request).
func (r *Route) Path(tpl string) *Route {
	r.myRoute.Path(tpl)
	return r
}

// PathPrefix -----------------------------------------------------------------

// PathPrefix adds a matcher for the URL path prefix. This matches if the given
// template is a prefix of the full URL path. See Route.Path() for details on
// the tpl argument.
//
// Note that it does not treat slashes specially ("/foobar/" will be matched by
// the prefix "/foo") so you may want to use a trailing slash here.
//
// Also note that the setting of Router.StrictSlash() has no effect on routes
// with a PathPrefix matcher.
func (r *Route) PathPrefix(tpl string) *Route {
	r.myRoute.PathPrefix(tpl)
	return r
}

// Query ----------------------------------------------------------------------

// Queries adds a matcher for URL query values.
// It accepts a sequence of key/value pairs. Values may define variables.
// For example:
//
//     r := mux.NewRouter()
//     r.Queries("foo", "bar", "id", "{id:[0-9]+}")
//
// The above route will only match if the URL contains the defined queries
// values, e.g.: ?foo=bar&id=42.
//
// If the value is an empty string, it will match any value if the key is set.
//
// Variables can define an optional regexp pattern to be matched:
//
// - {name} matches anything until the next slash.
//
// - {name:pattern} matches the given regexp pattern.
func (r *Route) Queries(pairs ...string) *Route {
	r.myRoute.Queries(pairs...)

	return r
}

// Schemes --------------------------------------------------------------------

// Schemes adds a matcher for URL schemes.
// It accepts a sequence of schemes to be matched, e.g.: "http", "https".
func (r *Route) Schemes(schemes ...string) *Route {
	r.myRoute.Schemes(schemes...)
	return r
}

func (r *Route) Subrouter() *Router {
	return &Router{myRouter: r.myRoute.Subrouter(), routeConfig: r.routeConfig, config: r.getConfig().GetCopy()}
}
