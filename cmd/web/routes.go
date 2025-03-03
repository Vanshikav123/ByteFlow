package main

import (
	"net/http"

	"github.com/Vanshikav123/ByteFlow.git/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
Handler vs. HandleFunc
Use Handler when:

You need to implement a complex handler with state or additional methods.

You want to reuse the same handler logic across multiple routes.

Use HandleFunc when:

You have a simple handler that doesn’t require state.

You want to quickly register a function as a handler.

Counter vs. CounterVec
Use Counter when:

You only need to track a single metric without any categorization.

Example: Total number of requests.

Use CounterVec when:

You need to track a metric with multiple dimensions (e.g., by endpoint, method, or status).

Example: Number of requests for each endpoint, grouped by HTTP method.
*/
func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	// Unprotected application routes using the "dynamic" middleware chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.withMetrics(app.home)))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.withMetrics(app.snippetView)))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.withMetrics(app.userSignup)))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.withMetrics(app.userSignupPost)))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.withMetrics(app.userLogin)))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.withMetrics(app.userLoginPost)))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.withMetrics(app.snippetCreate)))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.withMetrics(app.snippetCreatePost)))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.withMetrics(app.userLogoutPost)))

	// Metrics endpoint
	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
