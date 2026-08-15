package tmdb_dto

type ContentFiltersDto struct {
	Type       string
	Page       int
	SortBy     string
	WithGenres string
	Language   string
	Year       int
	Query      string
	IsKids     bool
}
