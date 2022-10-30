package main

import (
	"net/http"

	"github.com/terdia/mvp/pkg/dto"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.writeJson(w, http.StatusOK, dto.ResponseObject{
		StatusMsg: dto.Success,
		Message:   "available",
	}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
