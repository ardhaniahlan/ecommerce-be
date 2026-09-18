package utils

import "math"

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
}

func GeneratePagination(page, limit, totalItems int) Pagination {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return Pagination{
		CurrentPage: page,
		PageSize:    limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}
}