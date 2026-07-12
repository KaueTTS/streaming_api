package dto_shared

type PaginationDto struct {
	Page      int   `query:"page" json:"page"`
	PerPage   int   `query:"per_page" json:"per_page"`
	PageCount int   `json:"page_count"`
	Total     int64 `json:"total"`
}
