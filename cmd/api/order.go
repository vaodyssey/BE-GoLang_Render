package main

import (
	"TechStore/cache"
	db "TechStore/db/sqlc"
	data "TechStore/internal/dto/data/paginated_result"
	"TechStore/internal/dto/payload"
	"TechStore/internal/pkg/validator"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	OrderPending = 0
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var request payload.CreateOrderRequest
	err := app.readJson(w, r, &request)
	if err != nil {
		app.badRequestErrorResponse(w, err.Error())
		return
	}

	v := validator.New()
	if payload.ValidateCreateOrderRequest(v, request); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		app.badRequestErrorResponse(w, "Can not get token from header !")
		return
	}

	inMemoryUserId, ok := cache.Store.Get(token)
	if !ok {
		app.badRequestErrorResponse(w, "Invalid token !")
		return
	}

	if inMemoryUserId != request.Data.UserId {
		app.badRequestErrorResponse(w, "Invalid token !")
		return
	}

	var productIds []string
	for _, product := range request.Data.Products {
		productIds = append(productIds, product.Id)
	}

	productsCount, err := app.queries.CountProductByIds(r.Context(), productIds)
	if err != nil || int(productsCount) != len(productIds) {
		app.badRequestErrorResponse(w, "Product not found!")
		return
	}

	tx, err := app.db.Begin()
	if err != nil {
		app.serverErrorResponse(w, fmt.Errorf("cannot begin transaction"))
		return
	}
	defer func(tx *sql.Tx, ctx context.Context) {
		_ = tx.Rollback()
	}(tx, r.Context())

	qtx := app.queries.WithTx(tx)
	createOrderParam := db.CreateOrderParams{
		ID:        uuid.NewString(),
		Amount:    request.Data.Amount,
		UserID:    request.Data.UserId,
		Status:    OrderPending,
		CreatedAt: time.Now(),
	}

	err = qtx.CreateOrder(r.Context(), createOrderParam)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	for _, product := range request.Data.Products {
		createOrderDetailsParams := db.CreateOrderDetailsParams{
			OrderID:   createOrderParam.ID,
			ProductID: product.Id,
			Quantity:  int32(product.Quantity),
		}

		err = qtx.CreateOrderDetails(r.Context(), createOrderDetailsParams)
		if err != nil {
			app.serverErrorResponse(w, err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, fmt.Errorf("failed to commit transaction"))
	}
	err = app.writeJson(w, http.StatusOK, payload.BaseResponse{
		ResultCode:    SuccessCode,
		ResultMessage: SuccessMessage,
		Data: map[string]string{
			"orderId": createOrderParam.ID,
		},
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) getOrdersPaginated(w http.ResponseWriter, r *http.Request) {
	pageNumber, _ := strconv.Atoi(app.readStrings(r.URL.Query(), "pageNumber", "1"))
	pageSize, _ := strconv.Atoi(app.readStrings(r.URL.Query(), "pageSize", "10"))
	sortBy := app.readStrings(r.URL.Query(), "sortBy", "createdAt")
	sortOrder := app.readStrings(r.URL.Query(), "sortOrder", "DESC")
	userId := app.readStrings(r.URL.Query(), "userId", uuid.Nil.String())

	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		app.badRequestErrorResponse(w, "Can not get token from header !")
		return
	}

	inMemoryUserId, ok := cache.Store.Get(token)
	if !ok {
		app.badRequestErrorResponse(w, "Invalid token !")
		return
	}

	if inMemoryUserId != userId {
		app.badRequestErrorResponse(w, "Invalid token !")
		return
	}

	params := db.GetOrdersPaginatedParams{
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     int32(pageSize),
		Offset:    int32(pageSize * (pageNumber - 1)),
		UserID:    userId,
	}

	orders, err := app.queries.GetOrdersPaginated(r.Context(), params)
	if err != nil {
		orders = []db.Order{}
	}

	totalCount, err := app.queries.GetOrderTotalCount(r.Context(), userId)
	if err != nil {
		totalCount = 0
	}

	totalPages := math.Ceil(float64(totalCount) / float64(pageSize))

	paginatedResult := data.PaginatedResult{
		Items:       orders,
		TotalCount:  int(totalCount),
		HasNextPage: pageNumber < int(totalPages),
		HasPrevPage: pageNumber > 1,
		PageNumber:  pageNumber,
		PageSize:    pageSize,
	}

	err = app.writeJson(w, http.StatusOK, payload.BaseResponse{
		ResultCode:    SuccessCode,
		ResultMessage: SuccessMessage,
		Data:          paginatedResult,
	}, nil)

	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) getOrderById(w http.ResponseWriter, r *http.Request) {
	orderId, err := app.readIDParam(r, "orderId")
	if err != nil {
		app.badRequestErrorResponse(w, "orderId must be provide!")
		return
	}

	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		app.badRequestErrorResponse(w, "Can not get token from header !")
		return
	}

	userId, ok := cache.Store.Get(token)
	if !ok {
		app.badRequestErrorResponse(w, "Invalid token !")
		return
	}

	params := db.GetOrderByIdParams{
		UserID: userId,
		ID:     orderId.String(),
	}

	order, err := app.queries.GetOrderById(r.Context(), params)
	if err != nil {
		app.notFoundErrorResponse(w, "Order does not exist!")
		return
	}

	orderDetails, err := app.queries.GetOrderDetails(r.Context(), order.ID)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	result := payload.GetOrderByIdResponse{
		Id:           order.ID,
		Amount:       int(order.Amount),
		Status:       int(order.Status),
		CreatedAt:    order.CreatedAt,
		OrderDetails: orderDetails,
	}

	err = app.writeJson(w, http.StatusOK, payload.BaseResponse{
		ResultCode:    SuccessCode,
		ResultMessage: SuccessMessage,
		Data:          result,
	}, nil)

	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderId, err := app.readIDParam(r, "orderId")
	if err != nil {
		app.badRequestErrorResponse(w, "orderId must be provide!")
		return
	}

	var request payload.UpdateOrderRequest
	err = app.readJson(w, r, &request)
	if err != nil {
		app.badRequestErrorResponse(w, err.Error())
		return
	}

	v := validator.New()
	if payload.ValidateUpdateOrderRequest(v, request); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		app.badRequestErrorResponse(w, "Can not get token from header !")
		return
	}

	userId, ok := cache.Store.Get(token)
	if !ok {
		app.badRequestErrorResponse(w, "Invalid token !")
		return
	}

	getOrderParams := db.GetOrderByIdParams{
		UserID: userId,
		ID:     orderId.String(),
	}

	order, err := app.queries.GetOrderById(r.Context(), getOrderParams)
	if err != nil {
		app.notFoundErrorResponse(w, "Order does not exist!")
		return
	}

	if order.UserID != userId {
		app.unauthorizedErrorResponse(w, "Unauthorized user to access this order!")
		return
	}

	updateOrderParams := db.UpdateOrderStatusParams{
		Status: int32(request.Data.Status),
		ID:     order.ID,
	}

	err = app.queries.UpdateOrderStatus(r.Context(), updateOrderParams)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, payload.BaseResponse{
		ResultCode:    SuccessCode,
		ResultMessage: SuccessMessage,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
