package shared_normalizers

// NormalizePage garante que a página seja sempre um número inteiro positivo, retornando 1 se o valor for menor ou igual a zero.
func NormalizePage(page int) int {
	if page <= 0 {
		return 1
	}

	return page
}
