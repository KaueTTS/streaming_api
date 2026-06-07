package dto_auth

type RegisterRequestDto struct {
	Name     string `query:"name" json:"name"`
	Email    string `query:"email" json:"email"`
	Password string `query:"password" json:"password"`
}

type LoginRequestDto struct {
	Email    string `query:"email" json:"email"`
	Password string `query:"password" json:"password"`
}
