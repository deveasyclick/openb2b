package pagination

import (
	"math"
	"net/url"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"totalPages"`
	Total      int64 `json:"total"`
}

// ParseLimit returns the limit or default 20
const (
	DefaultLimit = 20
	MaxLimit     = 100
	DefaultPage  = 1
)

// FilterCondition represents a single filter on a field
type FilterCondition struct {
	Field    string
	Operator string // "=", "LIKE", "IN", ">=", "<=", etc.
	Value    interface{}
}

// Options holds pagination and filtering params
type Options struct {
	Page            int
	Limit           int
	Filters         []FilterCondition
	SortBy          string // e.g. "created_at desc"
	Preloads        []string
	SearchFields    []string
	SearchJoinQuery string // Optional: JOIN clauses if needed
}

var skipKeys map[string]bool = map[string]bool{
	"page":            true,
	"limit":           true,
	"sort":            true,
	"preloads":        true,
	"searchFields":    true,
	"searchJoinQuery": true,
	"search_fields":   true,
}

func BuildPagination(total int64, opts Options) Pagination {
	totalPages := 0
	if opts.Limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(opts.Limit)))
	}

	return Pagination{
		Total:      total,
		Limit:      opts.Limit,
		Page:       opts.Page,
		TotalPages: totalPages,
	}
}

func ParsePaginationOptions(query url.Values, allowedFields map[string]bool) (Options, error) {
	filters, err := parseFiltersFromQuery(query)
	if err != nil {
		return Options{}, err
	}

	return Options{
		Page:         parsePage(query.Get("page")),
		Limit:        parseLimit(query.Get("limit")),
		Filters:      filters,
		SortBy:       query.Get("sort"),
		Preloads:     parsePreloads(query.Get("preloads")),
		SearchFields: parseSearchFields(query, allowedFields),
	}, nil
}

// Sample pagination query: ?manufacturer_id_in=0,1&code=33000&manufacturer_Id=1&sort=code desc&preloads=Manufacturer&type_like=floor
func Paginate[T any](db *gorm.DB, opts Options) (items []T, total int64, err error) {
	// Count total matching records with filters
	countDB := db.Model(new(T))
	countDB = applyFilters(countDB, opts.Filters, opts.SearchFields, opts.SearchJoinQuery)
	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply filters, sorting, limit & offset
	query := db.Model(new(T))
	query = applyFilters(query, opts.Filters, opts.SearchFields, opts.SearchJoinQuery)

	if opts.SortBy != "" {
		query = query.Order(opts.SortBy)
	}

	if opts.Preloads != nil {
		for _, preload := range opts.Preloads {
			query = query.Preload(preload)
		}
	}

	offset := (opts.Page - 1) * opts.Limit
	query = query.Offset(offset).Limit(opts.Limit)

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ParseFiltersFromQuery parses URL query params into FilterConditions
func parseFiltersFromQuery(query url.Values) ([]FilterCondition, error) {
	var filters []FilterCondition

	for key, values := range query {
		if skipKeys[key] || len(values) == 0 {
			continue
		}

		value := values[0]

		field, op := parseFieldAndOperator(key)
		var filterValue interface{} = value

		switch op {
		case "IN":
			parts := strings.Split(value, ",")
			intVals := []int{}
			allInts := true
			for _, p := range parts {
				if iv, err := strconv.Atoi(p); err == nil {
					intVals = append(intVals, iv)
				} else {
					allInts = false
					break
				}
			}
			if allInts {
				filterValue = intVals
			} else {
				filterValue = parts
			}
		case "LIKE":
			if !strings.HasPrefix(value, "%") {
				value = "%" + value
			}
			if !strings.HasSuffix(value, "%") {
				value = value + "%"
			}
			filterValue = value
		// For GT, LT, GTE, LTE keep as string, let GORM convert
		default:
			// no change
		}

		filters = append(filters, FilterCondition{
			Field:    field,
			Operator: opToSymbol(op),
			Value:    filterValue,
		})
	}

	return filters, nil
}

func parseFieldAndOperator(param string) (string, string) {
	ops := []string{"_like", "_in", "_gte", "_lte", "_gt", "_lt"}
	for _, op := range ops {
		if strings.HasSuffix(param, op) {
			field := strings.TrimSuffix(param, op)
			return field, strings.ToUpper(strings.TrimPrefix(op, "_"))
		}
	}
	return param, "="
}

func opToSymbol(op string) string {
	switch op {
	case "LIKE":
		return "LIKE"
	case "IN":
		return "IN"
	case "GTE":
		return ">="
	case "LTE":
		return "<="
	case "GT":
		return ">"
	case "LT":
		return "<"
	case "=":
		return "="
	default:
		return "="
	}
}

// ParsePage returns the page number or default 1
func parsePage(pageStr string) int {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return DefaultPage
	}
	return page
}

func parseLimit(limitStr string) int {
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}
func parsePreloads(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var preloads []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			preloads = append(preloads, trimmed)
		}
	}
	return preloads
}

// applyFilters adds WHERE clauses for each filter condition
func applyFilters(db *gorm.DB, filters []FilterCondition, searchFields []string, joinQuery string) *gorm.DB {
	for _, f := range filters {
		if f.Field == "search" && f.Value.(string) != "" {
			value := "%" + f.Value.(string) + "%"
			if joinQuery != "" {
				db = db.Joins(joinQuery)
			}
			// Build the OR conditions dynamically
			if len(searchFields) > 0 {
				conditions := db
				for i, field := range searchFields {
					if i == 0 {
						conditions = conditions.Where(field+" ILIKE ?", value)
					} else {
						conditions = conditions.Or(field+" ILIKE ?", value)
					}
				}
				db = conditions
			}
			continue
		}

		switch f.Operator {
		case "IN":
			db = db.Where(f.Field+" IN ?", f.Value)
		case "LIKE":
			db = db.Where(f.Field+" LIKE ?", f.Value)
		case "=", ">", "<", ">=", "<=":
			db = db.Where(f.Field+" "+f.Operator+" ?", f.Value)
		default:
			db = db.Where(f.Field+" = ?", f.Value)
		}
	}
	return db
}

func parseSearchFields(query url.Values, allowedFields map[string]bool) []string {
	if len(allowedFields) == 0 {
		return nil
	}

	raw := query.Get("search_fields")
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var fields []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" && allowedFields[trimmed] {
			fields = append(fields, trimmed)
		}
	}

	return fields
}
