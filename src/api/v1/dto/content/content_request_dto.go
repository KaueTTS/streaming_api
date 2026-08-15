package dto_content

type ContentListRequestDto struct {
	Type      string `query:"type" json:"type"`
	Page      int    `query:"page" json:"page"`
	SortBy    string `query:"sort_by" json:"sort_by"`
	Genre     string `query:"genre" json:"genre"`
	Language  string `query:"language" json:"language"`
	Year      int    `query:"year" json:"year"`
	ProfileID uint   `query:"profile_id" json:"profile_id"`
}

type ContentSearchRequestDto struct {
	Type      string `query:"type" json:"type"`
	Page      int    `query:"page" json:"page"`
	Language  string `query:"language" json:"language"`
	Query     string `query:"query" json:"query"`
	ProfileID uint   `query:"profile_id" json:"profile_id"`
}
