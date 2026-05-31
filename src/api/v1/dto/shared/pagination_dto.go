package dto_shared

type PaginationDto struct {
	Page      int   `query:"page" json:"page"`
	PerPage   int   `query:"per_page" json:"per_page"`
	PageCount int   `json:"pageCount"`
	Total     int64 `json:"total"`
}
