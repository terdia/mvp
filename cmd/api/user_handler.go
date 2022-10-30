package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/pkg/dto"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserRequest
	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, validationErrors, err := app.userService.Create(input)
	if validationErrors != nil {
		app.failedValidationResponse(w, r, validationErrors)
		return
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))

	if err = app.writeJson(w, http.StatusCreated, dto.ResponseObject{
		StatusMsg: dto.Success,
		Data:      dto.UserResponse{User: getAPIUser(user)},
	}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAuthenticationToken(rw http.ResponseWriter, r *http.Request) {

	request := dto.AuthTokenRequest{}

	err := app.readJson(rw, r, &request)
	if err != nil {
		app.badRequestResponse(rw, r, err)
		return
	}

	token, validationErrors, err := app.userService.CreateAuthenticationToken(
		request, data.TokenScopeAuthentication,
	)
	if validationErrors != nil {
		app.failedValidationResponse(rw, r, validationErrors)
		return
	}

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(rw, r)
		case errors.Is(err, data.ErrInvalidCredentials):
			app.invalidCredentialsResponse(rw, r)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	tokenDto := dto.Token{
		PlainText: token.Plaintext,
		Expiry:    token.Expiry,
	}

	if err = app.writeJson(rw, http.StatusOK, dto.ResponseObject{
		StatusMsg: dto.Success,
		Data: dto.TokenResponse{
			Token: tokenDto,
		},
	}, nil); err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}
}

func (app *application) depositBalanceHandler(rw http.ResponseWriter, r *http.Request) {

	amount, err := app.extractIntParamFromContext(r, "amount")
	if err != nil {
		app.notFoundResponse(rw, r)
		return
	}

	user := app.contextGetUser(r)
	validationErrors, err := app.transactionService.DepositCoin(user, int(amount))
	if validationErrors != nil {
		app.failedValidationResponse(rw, r, validationErrors)
		return
	}

	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	if err = app.writeJson(rw, http.StatusOK, dto.ResponseObject{
		StatusMsg: dto.Success,
		Message:   "Deposit was successful",
		Data:      dto.UserResponse{User: getAPIUser(user)},
	}, nil); err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}
}

func (app *application) resetBalanceHandler(rw http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	validationErrors, err := app.transactionService.DepositReset(user)
	if validationErrors != nil {
		app.failedValidationResponse(rw, r, validationErrors)
		return
	}

	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	if err = app.writeJson(rw, http.StatusOK, dto.ResponseObject{
		StatusMsg: dto.Success,
		Message:   "Reset balance was successful",
		Data:      dto.UserResponse{User: getAPIUser(user)},
	}, nil); err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}
}

func getAPIUser(user *data.User) dto.APIUser {
	return dto.APIUser{
		ID:        user.ID,
		Role:      user.Role,
		Username:  user.Username,
		Deposit:   user.Deposit,
		CreatedAt: user.CreatedAt,
	}
}
