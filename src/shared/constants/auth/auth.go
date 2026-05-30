package shared_constants_auth

import "time"

const (
	AuthRateLimitMax        = 10
	AuthRateLimitExpiration = time.Minute

	MinPasswordLength = 8
	MaxPasswordBytes  = 72

	JWTAlgorithm = "HS256"
	JWTType      = "JWT"
)
