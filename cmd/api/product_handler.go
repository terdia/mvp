package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/pkg/dto"
	"github.com/terdia/mvp/pkg/validator"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.ProductRequest

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &data.Product{
		Name:            *input.Name,
		Cost:            input.Cost,
		AmountAvailable: input.Quantity,
		Seller:          *app.contextGetUser(r),
	}

	validationErrors, err := app.productService.Create(product)
	if validationErrors != nil {
		app.failedValidationResponse(w, r, validationErrors)
		return
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	result := dto.ResponseObject{
		StatusMsg: dto.Success,
		Data: dto.ProductResponse{
			Product: getProductResponse(product),
		},
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/products/%d", product.ID))

	if err = app.writeJson(w, http.StatusCreated, result, headers); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.extractIntParamFromContext(r, "id")
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	product, err := app.productService.GetOne(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	result := dto.ResponseObject{
		StatusMsg: dto.Success,
		Data: dto.ProductResponse{
			Product: getProductResponse(product),
		},
	}

	if err = app.writeJson(w, http.StatusOK, result, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listProductHandler(rw http.ResponseWriter, r *http.Request) {
	v := validator.New()

	qs := r.URL.Query()

	filters := data.Filters{
		Page:         app.readInt(qs, "page", 1, v),
		PageSize:     app.readInt(qs, "page_size", 10, v),
		Sort:         app.readString(qs, "sort", "id"),
		SortSafeList: []string{"id", "name", "-id", "-name"},
	}

	filters.ValidateFilters(v)
	if !v.Valid() {
		app.failedValidationResponse(rw, r, v.Errors)
		return
	}

	listRequest := dto.ListProductRequest{
		Name:    app.readString(qs, "name", ""),
		Filters: filters,
	}

	products, metadata, err := app.productService.List(listRequest)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	listProductResponse := dto.ListProductResponse{
		Products: []dto.APIProduct{},
	}

	for _, product := range products {
		listProductResponse.Products = append(listProductResponse.Products, getProductResponse(product))
	}

	if len(products) > 0 {
		listProductResponse.Metadata = &metadata
	}

	if err = app.writeJson(rw, http.StatusOK, dto.ResponseObject{
		StatusMsg: dto.Success,
		Data:      listProductResponse,
	}, nil); err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}
}

func (app *application) updateProductHandler(rw http.ResponseWriter, r *http.Request) {

	id, err := app.extractIntParamFromContext(r, "id")
	if err != nil {
		app.notFoundResponse(rw, r)
		return
	}

	var input struct {
		Name     string `json:"name"`
		Cost     int    `json:"cost"`
		Quantity int    `json:"amount_available"`
	}
	if err = app.readJson(rw, r, &input); err != nil {
		app.badRequestResponse(rw, r, err)
		return
	}

	product, validationErrors, err := app.productService.Update(data.Product{
		ID:              id,
		Cost:            input.Cost,
		Name:            input.Name,
		Seller:          *app.contextGetUser(r),
		CreatedAt:       time.Time{},
		AmountAvailable: input.Quantity,
	})

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(rw, r)
		case errors.Is(err, data.ErrNoPermission):
			app.notPermittedRResponse(rw, r)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	if validationErrors != nil {
		app.failedValidationResponse(rw, r, validationErrors)
		return
	}

	result := dto.ResponseObject{
		StatusMsg: dto.Success,
		Data: dto.ProductResponse{
			Product: getProductResponse(product),
		},
	}

	if err = app.writeJson(rw, http.StatusOK, result, nil); err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}
}

func (app *application) deleteProductHandler(rw http.ResponseWriter, r *http.Request) {

	id, err := app.extractIntParamFromContext(r, "id")
	if err != nil {
		app.notFoundResponse(rw, r)
		return
	}
	product := data.Product{
		ID:     id,
		Seller: *app.contextGetUser(r),
	}

	if err = app.productService.Remove(product); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(rw, r)
		case errors.Is(err, data.ErrNoPermission):
			app.notPermittedRResponse(rw, r)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	result := dto.ResponseObject{
		StatusMsg: dto.Success,
		Message:   "product successfully deleted",
	}

	if err = app.writeJson(rw, http.StatusOK, result, nil); err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}
}

func (app *application) buyProductHandler(rw http.ResponseWriter, r *http.Request) {

	id, err := app.extractIntParamFromContext(r, "id")
	if err != nil {
		app.notFoundResponse(rw, r)
		return
	}

	amount, err := app.extractIntParamFromContext(r, "amount")
	if err != nil {
		app.notFoundResponse(rw, r)
		return
	}

	product, err := app.productService.GetOne(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(rw, r)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	purchaseResponse, validationErrs, err := app.transactionService.BuyProduct(
		app.contextGetUser(r),
		product,
		int(amount),
	)

	if validationErrs != nil {
		app.failedValidationResponse(rw, r, validationErrs)
		return
	}

	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	app.logger.Printf("Purchase %+v", purchaseResponse)

	result := dto.ResponseObject{
		StatusMsg: dto.Success,
		Message:   "Buy product successful",
		Data: map[string]dto.BuyProductResponse{
			"purchase": *purchaseResponse,
		},
	}

	if err = app.writeJson(rw, http.StatusOK, result, nil); err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}
}

func getProductResponse(product *data.Product) dto.APIProduct {
	return dto.APIProduct{
		ID:              product.ID,
		Cost:            product.Cost,
		Name:            product.Name,
		CreatedAt:       product.CreatedAt,
		AmountAvailable: product.AmountAvailable,
	}
}
