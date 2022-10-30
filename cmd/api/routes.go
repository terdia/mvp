package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {

	router := chi.NewRouter()

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Use(app.recoverPanic, app.logRequest, app.enableCors)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("MVP Vending machine"))
	})

	router.Get("/v1/healthcheck", app.healthcheckHandler)

	router.Route("/v1/products", func(r chi.Router) {

		r.Post("/", app.requirePermission("products:write", app.createProductHandler))
		r.Get("/", app.listProductHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.showProductHandler)
			r.Put("/", app.requirePermission("products:write", app.updateProductHandler))
			r.Delete("/", app.requirePermission("products:write", app.deleteProductHandler))
		})

	})

	router.Route("/v1/users", func(r chi.Router) {
		r.Post("/", app.registerUserHandler)
	})

	router.Post("/v1/auth/tokens", app.getAuthenticationToken)

	return router
}
