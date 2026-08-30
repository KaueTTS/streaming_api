package dto_favorite

type FavoriteListRequestDto struct {
	ProfileID uint   `query:"profile_id" json:"profile_id"`
	Page      int    `query:"page" json:"page"`
	PerPage   int    `query:"per_page" json:"per_page"`
	Language  string `query:"language" json:"language"`
}

type FavoriteRequestDto struct {
	ProfileID         uint   `json:"profile_id"`
	ContentExternalID int    `json:"content_external_id"`
	Type              string `json:"type"`
}
