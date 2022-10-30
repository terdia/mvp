package main

import (
	"context"
	"github.com/terdia/mvp/internal/data"
	"net/http"
)

type contextKey string

const (
	userContextKey = contextKey("user")
)

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)

	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)

	if !ok {
		return nil
	}

	return user
}