package dto_profile

type CreateProfileRequestDto struct {
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`
	IsKids    bool    `json:"is_kids"`
}
