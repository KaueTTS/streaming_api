package security_test

import (
	"strings"
	"testing"

	security "github.com/KaueTTS/streaming_api/src/security"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "should return nil when password is valid",
			password:    "abc12345",
			expectedErr: nil,
		},
		{
			name:        "should return error when password has less than minimum characters",
			password:    "abc123",
			expectedErr: shared_errors.ErrPasswordMustLeast8Character,
		},
		{
			name:        "should return error when password has more than maximum bytes",
			password:    strings.Repeat("a", 72) + "1",
			expectedErr: shared_errors.ErrPasswordMustMaximum72Bytes,
		},
		{
			name:        "should return error when password has no numbers",
			password:    "abcdefgh",
			expectedErr: shared_errors.ErrPasswordMustLettersAndNumbers,
		},
		{
			name:        "should return error when password has no letters",
			password:    "12345678",
			expectedErr: shared_errors.ErrPasswordMustLettersAndNumbers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := security.Validate(tt.password)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestHash(t *testing.T) {
	t.Run("should return hash for password", func(t *testing.T) {
		password := "abc12345"

		hash, err := security.Hash(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
		assert.True(t, security.Compare(hash, password))
	})
}

func TestCompare(t *testing.T) {
	t.Run("should return true when password matches hash", func(t *testing.T) {
		password := "abc12345"

		hash, err := security.Hash(password)
		require.NoError(t, err)

		assert.True(t, security.Compare(hash, password))
	})

	t.Run("should return false when password does not match hash", func(t *testing.T) {
		password := "abc12345"

		hash, err := security.Hash(password)
		require.NoError(t, err)

		assert.False(t, security.Compare(hash, "wrong12345"))
	})

	t.Run("should return false when hash is invalid", func(t *testing.T) {
		assert.False(t, security.Compare("invalid-hash", "abc12345"))
	})
}
