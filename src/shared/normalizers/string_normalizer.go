package shared_normalizers

import "strings"

// NormalizeString normaliza uma string, convertendo-a para minúsculas e removendo espaços em branco no início e no final.
func NormalizeString(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

// TrimString remove espaços em branco no início e no final de uma string.
func TrimString(value string) string {
	return strings.TrimSpace(value)
}
