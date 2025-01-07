package utils

import (
	"backend/internal/datasources/repositories"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ExtractQueryParams(query url.Values) (repositories.QueryParams, error) {
	filters := make(map[string]interface{})
	pagination := repositories.Pagination{
		Page:  1,
		Limit: 10,
	}

	for key, values := range query {
		if len(values) == 0 {
			continue
		}
		value := values[0]

		switch {
		case key == "page":
			page, err := strconv.Atoi(value)
			if err != nil || page < 1 {
				return repositories.QueryParams{}, fmt.Errorf("invalid page number")
			}
			pagination.Page = page
		case key == "limit":
			limit, err := strconv.Atoi(value)
			if err != nil || limit < 1 {
				return repositories.QueryParams{}, fmt.Errorf("invalid limit number")
			}
			pagination.Limit = limit
		case strings.HasPrefix(key, "min_"):
			filters[strings.TrimPrefix(key, "min_")+">="] = value
		case strings.HasPrefix(key, "max_"):
			filters[strings.TrimPrefix(key, "max_")+"<="] = value
		case strings.HasPrefix(key, "like_"):
			filters[strings.TrimPrefix(key, "like_")+"<>"] = fmt.Sprintf("%%%s%%", strings.ToLower(value))
		default:
			filters[key] = value
		}
	}

	return repositories.QueryParams{
		Filters:    filters,
		Pagination: pagination,
	}, nil
}
