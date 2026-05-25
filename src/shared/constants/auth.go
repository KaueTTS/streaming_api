package shared_constants

import "time"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

const (
	AuthRateLimitMax        = 10
	AuthRateLimitExpiration = time.Minute

	MinNameLength = 2
	MaxNameLength = 120

	MinPasswordLength = 8
	MaxPasswordBytes  = 72

	JWTAlgorithm = "HS256"
	JWTType      = "JWT"

	Hidden = "hidden"
)
