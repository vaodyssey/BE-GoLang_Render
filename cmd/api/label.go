package main

import (
	db "TechStore/db/sqlc"
	data "TechStore/internal/dto/data/paginated_result"
	"TechStore/internal/dto/payload"
	"math"
	"net/http"
	"strconv"
)

func (app *application) getAllLabelsPaginated(w http.ResponseWriter, r *http.Request) {
	pageNumber, _ := strconv.Atoi(app.readStrings(r.URL.Query(), "pageNumber", "1"))
	pageSize, _ := strconv.Atoi(app.readStrings(r.URL.Query(), "pageSize", "10"))

	params := db.GetAllLabelsPaginatedParams{
		Limit:  int32(pageSize),
		Offset: int32(pageSize * (pageNumber - 1)),
	}

	labels, err := app.queries.GetAllLabelsPaginated(r.Context(), params)
	if err != nil {
		labels = []db.Label{}
	}

	totalLabelsCount, err := app.queries.GetLabelTotalCounts(r.Context())
	if err != nil {
		totalLabelsCount = 0
	}

	totalPages := math.Ceil(float64(totalLabelsCount) / float64(pageSize))

	paginatedResult := data.PaginatedResult{
		Items:       labels,
		TotalCount:  int(totalLabelsCount),
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
