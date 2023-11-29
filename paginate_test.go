package napi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	ID int `json:"id" validate:"required"`
	Paginater
}

func TestTableManager_Validate_ShouldSetDefaultsProperly(t *testing.T) {
	cp := new(TestRequest)
	// validate here
	cp.ID = 2

	cp.Validate(PaginaterData{
		Page:     3,
		Limit:    25,
		OrderBy:  "test",
		OrderDir: "desc",
		FilterBy: "all",
		Search:   "test",
	})

	assert.Equal(t, 3, cp.Page)
}

func TestTableManager_ValidateSearch_ShouldSetDefaultsProperly(t *testing.T) {
	cp := new(TestRequest)
	// validate here
	cp.ID = 2

	cp.ValidateSearch("test")
	assert.Equal(t, "test", cp.Search)
}

func TestTableManager_ValidatePagination_ShouldSetDefaultsProperly(t *testing.T) {
	cp := new(TestRequest)
	// validate here
	cp.ID = 2

	cp.ValidatePagination(3, 25)

	assert.Equal(t, 3, cp.Page)
	assert.Equal(t, 25, cp.Limit)
}

func TestTableManager_ValidateFilter_ShouldSetDefaultsProperly(t *testing.T) {
	cp := new(TestRequest)
	// validate here
	cp.ID = 2

	cp.ValidateFilter("all")

	assert.Equal(t, "all", cp.FilterBy)
}

func TestTableManager_ValidateOrder_ShouldSetDefaultsProperly(t *testing.T) {
	cp := new(TestRequest)
	// validate here
	cp.ID = 2

	cp.ValidateOrder("test", "desc")

	assert.Equal(t, "test", cp.OrderBy)
	assert.Equal(t, "desc", cp.OrderDir)
}
