package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/terdia/mvp/pkg/dto"
	"github.com/terdia/mvp/pkg/validator"
)

func (app *application) extractIntParamFromContext(r *http.Request, key string) (int64, error) {

	stringId := chi.URLParam(r, key)
	if stringId == "" {
		return 0, fmt.Errorf("invalid parameter, missing key %s", key)
	}

	id, err := strconv.ParseInt(chi.URLParam(r, key), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s, expected int", key)
	}

	return id, nil
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {

	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (app *application) readString(qs url.Values, key, defaultValue string) string {

	str := qs.Get(key)
	if str == "" {
		return defaultValue
	}

	return str
}

func (app *application) writeJson(rw http.ResponseWriter, status int, envelop dto.ResponseObject, headers http.Header) error {

	js, err := json.MarshalIndent(envelop, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		rw.Header()[key] = value
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(js)

	return nil
}

func (app *application) readJson(rw http.ResponseWriter, r *http.Request, dst interface{}) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(rw, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("request boby must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err

		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) errorResponse(rw http.ResponseWriter, r *http.Request, status int, envelop dto.ResponseObject) {

	if envelop.StatusMsg == 0 {
		envelop.SetStatus(dto.Error)
	}

	err := app.writeJson(rw, status, envelop, nil)
	if err != nil {
		app.logErrorWithHttpRequestContext(r, err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
