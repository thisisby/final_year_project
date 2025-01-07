package repositories

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"strings"
)

type Filters map[string]interface{}

type QueryParams struct {
	Filters    map[string]interface{}
	Pagination Pagination
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func ApplyFilters(queryBuilder squirrel.SelectBuilder, filters map[string]interface{}) squirrel.SelectBuilder {
	for key, value := range filters {
		switch {
		case strings.HasSuffix(key, ">="):
			column := strings.TrimSuffix(key, ">=")
			queryBuilder = queryBuilder.Where(fmt.Sprintf("%s >= ?", column), value)

		case strings.HasSuffix(key, "<="):
			column := strings.TrimSuffix(key, "<=")
			queryBuilder = queryBuilder.Where(fmt.Sprintf("%s <= ?", column), value)

		case strings.HasSuffix(key, "<>"):
			column := strings.TrimSuffix(key, "<>")
			queryBuilder = queryBuilder.Where(fmt.Sprintf("LOWER(%s) ILIKE ?", column), value)

		default:
			queryBuilder = queryBuilder.Where(squirrel.Eq{key: value})
		}
	}

	return queryBuilder
}

func ApplyPagination(builder squirrel.SelectBuilder, page, limit int) squirrel.SelectBuilder {
	if limit > 0 {
		offset := (page - 1) * limit
		builder = builder.Limit(uint64(limit)).Offset(uint64(offset))
	}
	return builder
}
