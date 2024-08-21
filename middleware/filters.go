package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Kelado/DeviceService/models"
)

type filterCtxKey string

const filtersKey filterCtxKey = "filters"

const SearchKeyword = "s"

func Filters() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var filters []models.Filter

			searchParam := r.URL.Query().Get(SearchKeyword)
			if searchParam != "" {
				// This will work only for one search term
				// If more terms are wanted, first split on comma
				searchPair := strings.Split(searchParam, ":")
				if len(searchPair) == 2 {
					term := searchPair[0]
					value := searchPair[1]
					if term == "brand" {
						filters = append(filters, models.NewBrandFilter(value))
					}
				}
			}

			ctx := context.WithValue(r.Context(), filtersKey, filters)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetFilterFromCtx(r *http.Request) []models.Filter {
	filters, ok := r.Context().Value(filtersKey).([]models.Filter)
	if !ok {
		return []models.Filter{}
	}
	return filters
}
