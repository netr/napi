package napi

type Paginater struct {
	CanPaginate
	CanOrder
	CanFilter
	CanSearch
}

type PaginaterData struct {
	Page     int
	Limit    int
	OrderBy  string
	OrderDir string
	FilterBy string
	Search   string
}

func (p *Paginater) Validate(data PaginaterData) {
	p.ValidatePagination(data.Page, data.Limit)
	p.ValidateOrder(data.OrderBy, data.OrderDir)
	p.ValidateFilter(data.FilterBy)
	p.ValidateSearch(data.Search)
}

type CanPaginate struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type CanOrder struct {
	OrderBy  string `query:"order_by"`
	OrderDir string `query:"order_dir"`
}

type CanFilter struct {
	FilterBy string `query:"filter_by"`
}

type CanSearch struct {
	Search string `query:"search"`
}

func (req *CanPaginate) Offset() int {
	return (req.Page - 1) * req.Limit
}

func (req *CanPaginate) ValidatePagination(page, limit int) {
	if req.Page < 1 {
		req.Page = page
	}
	if req.Limit < 1 {
		req.Limit = limit
	}
	if req.Limit > 100 {
		req.Limit = limit
	}
}

func (req *CanOrder) ValidateOrder(orderBy, orderDir string) {
	if req.OrderBy == "" {
		req.OrderBy = orderBy
	}
	if req.OrderDir != "asc" && req.OrderDir != "desc" {
		req.OrderDir = orderDir
	}
}

func (req *CanFilter) ValidateFilter(filterBy string) {
	if req.FilterBy == "" {
		req.FilterBy = filterBy
	}
}

func (req *CanSearch) ValidateSearch(search string) {
	if req.Search == "" {
		req.Search = search
	}
}
