package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
)

type Claims struct {
	UserID uint   `json:"sub"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"` // expires at
	Iat    int64  `json:"iat"` // issued at
	Iss    string `json:"iss"` // issuer
}

type tokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func GenerateToken(userID uint, email string) (string, error) {
	secret := strings.TrimSpace(env.JWTSecret)
	if secret == "" {
		return "", errors.New("JWT_SECRET não informado")
	}

	now := time.Now()
	expiresAt := now.Add(GetExpirationDuration())

	header := tokenHeader{
		Alg: shared_constants.JWTAlgorithm,
		Typ: shared_constants.JWTType,
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		Iat:    now.Unix(),
		Exp:    expiresAt.Unix(),
		Iss:    env.AppName,
	}

	headerEncoded, err := encodeJSON(header)
	if err != nil {
		return "", err
	}

	payloadEncoded, err := encodeJSON(claims)
	if err != nil {
		return "", err
	}

	unsignedToken := headerEncoded + "." + payloadEncoded
	signature := sign(unsignedToken, secret)

	return unsignedToken + "." + signature, nil
}

func ValidateToken(token string) (*Claims, error) {
	secret := strings.TrimSpace(env.JWTSecret)
	if secret == "" {
		return nil, errors.New("JWT_SECRET não informado")
	}

	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 {
		return nil, errors.New(shared_errors.InvalidToken)
	}

	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return nil, errors.New(shared_errors.InvalidToken)
		}
	}

	headerPayload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("header do token inválido")
	}

	var header tokenHeader
	if err := json.Unmarshal(headerPayload, &header); err != nil {
		return nil, errors.New("header do token inválido")
	}

	if header.Alg != shared_constants.JWTAlgorithm || !strings.EqualFold(header.Typ, shared_constants.JWTType) {
		return nil, errors.New("algoritmo do token inválido")
	}

	unsignedToken := parts[0] + "." + parts[1]
	expectedSignature := sign(unsignedToken, secret)

	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, errors.New("assinatura do token inválida")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("payload do token inválido")
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("claims do token inválidas")
	}

	now := time.Now().Unix()
	if claims.Iat <= 0 || claims.Exp <= 0 || claims.Exp <= claims.Iat {
		return nil, errors.New("claims do token inválidas")
	}

	if claims.Exp <= now {
		return nil, errors.New("token expirado")
	}

	if claims.Iat > now+60 {
		return nil, errors.New("token emitido em data inválida")
	}

	if claims.Iss != env.AppName || claims.UserID == 0 || strings.TrimSpace(claims.Email) == "" {
		return nil, errors.New("claims do token inválidas")
	}

	return &claims, nil
}

func ExtractBearerToken(authorizationHeader string) (string, error) {
	authorizationHeader = strings.TrimSpace(authorizationHeader)
	if authorizationHeader == "" {
		return "", errors.New("header Authorization não informado")
	}

	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("header Authorization deve estar no formato Bearer token")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("token não informado")
	}

	return token, nil
}

func encodeJSON(value any) (string, error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar json do token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(jsonBytes), nil
}

func sign(value string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(value))

	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func GetExpirationDuration() time.Duration {
	if env.AuthTokenExpirationTimeInHours <= 0 {
		return 8 * time.Hour
	}

	return time.Duration(env.AuthTokenExpirationTimeInHours * float64(time.Hour))
}
