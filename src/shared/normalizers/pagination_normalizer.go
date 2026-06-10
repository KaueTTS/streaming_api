package shared_normalizers

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
