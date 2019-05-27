package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	myRouter *mux.Router
}

type Route struct {
	myRoute *mux.Route
}

func (r *Router) Get(name string) *Route {
	return &Route{r.myRouter.Get(name)}
}

func (r *Router) GetRoute(name string) *Route {
	return &Route{r.myRouter.GetRoute(name)}
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
	return &Route{r.myRouter.NewRoute()}
}

// Name registers a new route with a name.
// See Route.Name().
func (r *Router) Name(name string) *Route {
	return &Route{r.myRouter.Name(name)}
}

// Handle registers a new route with a matcher for the URL path.
// See Route.Path() and Route.Handler().
func (r *Router) Handle(path string, handler http.Handler) *Route {
	return &Route{r.myRouter.Handle(path, handler)}
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.Path() and Route.HandlerFunc().
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return &Route{r.myRouter.HandleFunc(path, f)}
}

// Headers registers a new route with a matcher for request header values.
// See Route.Headers().
func (r *Router) Headers(pairs ...string) *Route {
	return &Route{r.myRouter.Headers(pairs...)}
}

// Host registers a new route with a matcher for the URL host.
// See Route.Host().
func (r *Router) Host(tpl string) *Route {
	return &Route{r.myRouter.Host(tpl)}
}


// Methods registers a new route with a matcher for HTTP methods.
// See Route.Methods().
func (r *Router) Methods(methods ...string) *Route {
	return &Route{r.myRouter.Methods(methods...)}
}

// Path registers a new route with a matcher for the URL path.
// See Route.Path().
func (r *Router) Path(tpl string) *Route {
	return &Route{r.myRouter.Path(tpl)}
}

// PathPrefix registers a new route with a matcher for the URL path prefix.
// See Route.PathPrefix().
func (r *Router) PathPrefix(tpl string) *Route {
	return &Route{r.myRouter.PathPrefix(tpl)}
}

// Queries registers a new route with a matcher for URL query values.
// See Route.Queries().
func (r *Router) Queries(pairs ...string) *Route {
	return &Route{r.myRouter.Queries(pairs...)}
}

// Schemes registers a new route with a matcher for URL schemes.
// See Route.Schemes().
func (r *Router) Schemes(schemes ...string) *Route {
	return &Route{r.myRouter.Schemes(schemes...)}
}



// SkipClean reports whether path cleaning is enabled for this route via
// Router.SkipClean.
func (r *Route) SkipClean() bool {
	return r.myRoute.SkipClean()
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

