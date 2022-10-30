package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/pkg/dto"
)

func (app *application) routes() http.Handler {

	router := chi.NewRouter()

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Use(app.recoverPanic, app.logRequest, app.enableCors)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_ = app.writeJson(w, http.StatusOK, dto.ResponseObject{
			StatusMsg: dto.Success,
			Message:   "MVP Vending machine",
		}, nil)
	})

	router.Get("/v1/healthcheck", app.healthcheckHandler)

	router.Route("/v1/products", func(r chi.Router) {
		r.Post("/", app.requirePermission(data.PermissionProductsWrite, app.createProductHandler))
		r.Get("/", app.listProductHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.showProductHandler)
			r.Put("/", app.requirePermission(data.PermissionProductsWrite, app.updateProductHandler))
			r.Delete("/", app.requirePermission(data.PermissionProductsWrite, app.deleteProductHandler))

			r.Get("/buy/{amount}", app.requirePermission(data.PermissionProductsBuy, app.buyProductHandler))
		})

	})

	router.Route("/v1/users", func(r chi.Router) {
		r.Post("/", app.registerUserHandler)
		r.Get("/deposit/{amount}", app.requirePermission(data.PermissionProductsBuy, app.depositBalanceHandler))
		r.Get("/deposit/reset", app.requirePermission(data.PermissionProductsBuy, app.resetBalanceHandler))
	})

	router.Post("/v1/auth/tokens", app.getAuthenticationToken)

	return router
}
