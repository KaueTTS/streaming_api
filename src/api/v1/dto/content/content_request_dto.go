package dto_content

type ContentListRequestDto struct {
	Type       string `query:"type" json:"type"`
	Page       int    `query:"page" json:"page"`
	SortBy     string `query:"sort_by" json:"sort_by"`
	WithGenres string `query:"with_genres" json:"with_genres"`
	Language   string `query:"language" json:"language"`
	Year       int    `query:"year" json:"year"`
}

type ContentSearchRequestDto struct {
	Type     string `query:"type" json:"type"`
	Page     int    `query:"page" json:"page"`
	Language string `query:"language" json:"language"`
	Query    string `query:"query" json:"query"`
}
