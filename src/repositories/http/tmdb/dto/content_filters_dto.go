package tmdb_dto

type ContentListFiltersDto struct {
	Type       string
	Page       int
	SortBy     string
	WithGenres string
	Language   string
	Year       int
}

type ContentSearchFiltersDto struct {
	Type     string
	Page     int
	Language string
	Query    string
}
