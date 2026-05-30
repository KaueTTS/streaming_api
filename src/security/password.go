package security

import (
	"unicode"
	"unicode/utf8"

	shared_constants_auth "github.com/KaueTTS/streaming_api/src/shared/constants/auth"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"golang.org/x/crypto/bcrypt"
)

func Validate(password string) error {
	if utf8.RuneCountInString(password) < shared_constants_auth.MinPasswordLength {
		return shared_errors.ErrPasswordMustLeast8Character
	}

	if len([]byte(password)) > shared_constants_auth.MaxPasswordBytes {
		return shared_errors.ErrPasswordMustMaximum72Bytes
	}

	var hasLetter, hasNumber bool

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}

		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	if !hasLetter || !hasNumber {
		return shared_errors.ErrPasswordMustLettersAndNumbers
	}

	return nil
}

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func Compare(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
