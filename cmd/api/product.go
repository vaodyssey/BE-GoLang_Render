package main

import (
	"TechStore/internal/dto/payload"
	"net/http"
)

func (app *application) getProductDetailHandler(w http.ResponseWriter, r *http.Request) {
	productId, err := app.readIDParam(r, "productId")
	if err != nil {
		app.badRequestErrorResponse(w, "Can not get product Id from params.")
		return
	}

	product, err := app.queries.GetProductById(r.Context(), productId.String())
	if err != nil {
		app.notFoundErrorResponse(w, "Product does not exist!")
		return
	}

	err = app.writeJson(w, http.StatusOK, payload.BaseResponse{
		ResultCode:    SuccessCode,
		ResultMessage: SuccessMessage,
		Data:          product,
	}, nil)

	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
