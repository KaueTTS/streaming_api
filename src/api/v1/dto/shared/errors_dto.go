package dto_shared

// DTO de erro
type ErrorDto struct {
	Message     string           `json:"message"`
	CodeMessage string           `json:"codeMessage,omitempty"`
	Details     []DetailErrorDto `json:"details,omitempty"`
}

// Detalhes do erro
type DetailErrorDto struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// Códigos de erro
const (
	BadRequest          = "BAD_REQUEST"           // 400
	Unauthorized        = "UNAUTHORIZED"          // 401
	Forbidden           = "FORBIDDEN"             // 403
	NotFound            = "NOT_FOUND"             // 404
	MethodNotAllowed    = "METHOD_NOT_ALLOWED"    // 405
	Conflict            = "CONFLICT"              // 409
	TooManyRequests     = "TOO_MANY_REQUESTS"     // 429
	InternalServerError = "INTERNAL_SERVER_ERROR" // 500
	BadGateway          = "BAD_GATEWAY"           // 502
)
