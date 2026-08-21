package pagination

import "math"

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type Params struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func New(page, limit int) Params {
	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return Params{Page: page, Limit: limit}
}

func (p Params) Skip() int64 {
	return int64((p.Page - 1) * p.Limit)
}

func (p Params) Limit64() int64 {
	return int64(p.Limit)
}

func NewMeta(p Params, total int64) Meta {
	totalPages := 0
	if p.Limit > 0 && total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(p.Limit)))
	}
	return Meta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
