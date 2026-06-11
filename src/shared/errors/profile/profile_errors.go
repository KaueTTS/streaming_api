package shared_errors_profile

// Erros de operação
const (
	FailedToListProfiles  = "erro ao listar perfis"
	FailedToCreateProfile = "erro ao criar perfil"
	FailedToUpdateProfile = "erro ao atualizar perfil"
	FailedToDeleteProfile = "erro ao remover perfil"
)

// Erros de validação
const (
	InvalidCreateProfileData = "dados inválidos para criar perfil"
	InvalidUpdateProfileData = "dados inválidos para atualizar perfil"
	InvalidProfileID         = "id do perfil inválido"
)

// Erros de negócio
const (
	ProfileLimitReached = "o usuário pode ter no máximo 3 perfis"
	ProfileNotFound     = "perfil não encontrado"
)
