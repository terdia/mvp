package main

import (
	"fmt"
	"net/http"

	"github.com/terdia/mvp/pkg/dto"
)

func (app *application) logErrorWithHttpRequestContext(r *http.Request, err error) {
	app.logger.Err(err).Msgf("%+v", map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) serverErrorResponse(rw http.ResponseWriter, r *http.Request, err error) {
	app.logger.Err(err).Msg("")

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(rw, r, http.StatusInternalServerError, dto.ResponseObject{
		Message: message,
	})
}

func (app *application) notFoundResponse(rw http.ResponseWriter, r *http.Request) {

	app.errorResponse(rw, r, http.StatusNotFound, dto.ResponseObject{
		Message: "the requested resource could not be found",
	})
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, dto.ResponseObject{
		Message: message,
	})
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, dto.ResponseObject{
		Message: err.Error(),
	})
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, dto.ResponseObject{
		StatusMsg: dto.Fail,
		Data: dto.ValidationError{
			Errors: errors,
		},
	})
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {

	app.errorResponse(w, r, http.StatusUnauthorized, dto.ResponseObject{
		StatusMsg: dto.Fail,
		Message:   "invalid authentication credentials",
	})
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	app.errorResponse(w, r, http.StatusUnauthorized, dto.ResponseObject{
		StatusMsg: dto.Fail,
		Message:   "invalid or missing token",
	})
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusUnauthorized, dto.ResponseObject{
		Message: "you must be authenticated to access this resource",
	})
}

func (app *application) notPermittedRResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusForbidden, dto.ResponseObject{
		Message: "your user account doesn't have the necessary permissions to perform this operation",
	})
}
