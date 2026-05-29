package validator

import (
	"net/mail"
	"strings"
	"unicode"
	"unicode/utf8"

	dto_auth "github.com/KaueTTS/streaming_api/src/api/v1/dto/auth"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
)

func ValidateRegisterRequest(request dto_auth.RegisterRequestDto) []dto_shared.DetailErrorDto {
	var details []dto_shared.DetailErrorDto

	name := strings.TrimSpace(request.Name)
	if name == "" {
		details = append(
			details,
			NewDetail(shared_constants.Name, name, shared_errors_auth.NameRequired),
		)
	} else if utf8.RuneCountInString(name) < shared_constants.MinNameLength {
		details = append(
			details,
			NewDetail(shared_constants.Name, name, shared_errors_auth.NameMustLeast2Character),
		)
	} else if utf8.RuneCountInString(name) > shared_constants.MaxNameLength {
		details = append(
			details,
			NewDetail(shared_constants.Name, name, shared_errors_auth.NameMustMaximum120Character),
		)
	}

	details = append(details, validateEmail(request.Email)...)

	password := request.Password
	if strings.TrimSpace(password) == "" {
		details = append(
			details,
			NewDetail(shared_constants.Password, "", shared_errors_auth.PasswordRequired),
		)
	} else if utf8.RuneCountInString(password) < shared_constants.MinPasswordLength {
		details = append(
			details,
			NewDetail(shared_constants.Password, shared_constants.Hidden, shared_errors_auth.PasswordMustLeast8Character),
		)
	} else if len([]byte(password)) > shared_constants.MaxPasswordBytes {
		details = append(
			details,
			NewDetail(shared_constants.Password, shared_constants.Hidden, shared_errors_auth.PasswordMustMaximum72Bytes),
		)
	} else if !hasLetterAndNumber(password) {
		details = append(
			details,
			NewDetail(shared_constants.Password, shared_constants.Hidden, shared_errors_auth.PasswordMustLettersAndNumbers),
		)
	}

	return details
}

func ValidateLoginRequest(request dto_auth.LoginRequestDto) []dto_shared.DetailErrorDto {
	var details []dto_shared.DetailErrorDto

	details = append(details, validateEmail(request.Email)...)

	if strings.TrimSpace(request.Password) == "" {
		details = append(
			details,
			NewDetail(shared_constants.Password, shared_constants.Hidden, shared_errors_auth.PasswordRequired),
		)
	}

	return details
}

func validateEmail(email string) []dto_shared.DetailErrorDto {
	trimmedEmail := strings.TrimSpace(email)
	if trimmedEmail == "" {
		return []dto_shared.DetailErrorDto{NewDetail(shared_constants.Email, email, shared_errors_auth.EmailRequired)}
	}

	parsedEmail, err := mail.ParseAddress(trimmedEmail)
	if err != nil || parsedEmail.Address != trimmedEmail {
		return []dto_shared.DetailErrorDto{NewDetail(shared_constants.Email, email, shared_errors_auth.EmailInvalid)}
	}

	return nil
}

func hasLetterAndNumber(password string) bool {
	var hasLetter, hasNumber bool

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}

		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}
