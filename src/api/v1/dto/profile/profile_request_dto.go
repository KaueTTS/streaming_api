package dto_profile

type ProfileRequestDto struct {
	Name      string  `query:"name" json:"name"`
	AvatarURL *string `query:"avatar_url" json:"avatar_url"`
	IsKids    bool    `query:"is_kids" json:"is_kids"`
}
