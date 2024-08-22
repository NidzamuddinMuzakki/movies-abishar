package util

import "math"

func Pagination(totalData int64, limit int64, page int64) (totalPages, previousPage, nextPage int64) {
	totalPages = int64(math.Ceil(float64(totalData) / float64(limit)))
	previousPage = page - 1
	if page > totalPages {
		previousPage = 0
	}
	nextPage = page + 1
	if nextPage > totalPages {
		nextPage = 0
	}

	return
}
