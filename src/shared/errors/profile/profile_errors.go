package shared_errors_profile

const (
	FailedToListProfiles  = "erro ao listar perfis"
	FailedToCreateProfile = "erro ao criar perfil"
	FailedToUpdateProfile = "erro ao atualizar perfil"
)

const (
	InvalidCreateProfileData = "dados inválidos para criar perfil"
	InvalidUpdateProfileData = "dados inválidos para atualizar perfil"
	InvalidProfileID         = "id do perfil inválido"
)

const (
	ProfileLimitReached = "o usuário pode ter no máximo 3 perfis"
	ProfileNotFound     = "perfil não encontrado"
)
