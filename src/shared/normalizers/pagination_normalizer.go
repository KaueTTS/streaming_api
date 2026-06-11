package shared_normalizers

import (
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
)

// NormalizePage garante que a página seja sempre um número inteiro positivo, retornando 1 se o valor for menor ou igual a zero.
func NormalizePage(page int) int {
	if page <= 0 {
		return 1
	}

	return page
}

// NormalizePerPage garante que o número de itens por página seja um número inteiro positivo, retornando 10 se o valor for menor ou igual a zero, e limitando a 100 para evitar sobrecarga.
func NormalizePerPage(perPage int) int {
	if perPage <= 0 {
		return 10
	}

	if perPage > 100 {
		return 100
	}

	return perPage
}

// NormalizePagination é uma função de conveniência que normaliza tanto a página quanto o número de itens por página usando as funções NormalizePage e NormalizePerPage, respectivamente.
func NormalizePagination(pagination dto_shared.PaginationDto) (page int, perPage int) {
	page = NormalizePage(pagination.Page)
	perPage = NormalizePerPage(pagination.PerPage)

	return page, perPage
}
