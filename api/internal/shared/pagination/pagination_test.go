package pagination

import (
	"net/url"
	"regexp"
	"sort"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	ID    int
	Name  string
	Price float64
}

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock: %v", err)
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening gorm DB: %v", err)
	}
	return gdb, mock
}

func TestParsePageAndLimit(t *testing.T) {
	assert.Equal(t, DefaultPage, parsePage(""))
	assert.Equal(t, 2, parsePage("2"))
	assert.Equal(t, DefaultLimit, parseLimit(""))
	assert.Equal(t, MaxLimit, parseLimit("1000"))
}

func TestBuildPagination(t *testing.T) {
	opts := Options{Page: 1, Limit: 10}
	result := BuildPagination(50, opts)
	assert.Equal(t, 5, result.TotalPages)
	assert.Equal(t, int64(50), result.Total)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 1, result.Page)
}

func TestParseFiltersFromQuery(t *testing.T) {
	values := url.Values{}
	values.Set("price_gte", "100")
	values.Set("name_like", "chair")
	values.Set("id_in", "1,2,3")
	filters, err := parseFiltersFromQuery(values)
	assert.NoError(t, err)
	assert.Len(t, filters, 3)

	// Sort by field
	sort.Slice(filters, func(i, j int) bool {
		return filters[i].Field < filters[j].Field
	})

	assert.Equal(t, ">=", filters[0].Operator)   // e.g., "id"
	assert.Equal(t, "LIKE", filters[1].Operator) // e.g., "name"
	assert.Equal(t, "IN", filters[2].Operator)   // e.g., "price"
}

func TestPaginate(t *testing.T) {
	db, mock := setupMockDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" ORDER BY id LIMIT $1`)).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(1, "Chair", 50.0).
			AddRow(2, "Table", 100.0))

	opts := Options{
		Page:   1,
		Limit:  10,
		SortBy: "id",
	}

	var items []Product
	items, total, err := Paginate[Product](db, opts)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, items, 2)
	assert.Equal(t, "Chair", items[0].Name)
}

func TestParsePreloads(t *testing.T) {
	preloads := parsePreloads("Manufacturer,Category")
	assert.Equal(t, []string{"Manufacturer", "Category"}, preloads)
}

func TestParseSearchFields(t *testing.T) {
	values := url.Values{}
	values.Set("search_fields", "name,price")
	allowed := map[string]bool{"name": true, "price": true}
	fields := parseSearchFields(values, allowed)
	assert.Equal(t, []string{"name", "price"}, fields)
}
