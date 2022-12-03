package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)


func TestHashPassword(t *testing.T){

  password := RandomString(6)

  hashedPassword, err := HashPassword(password)
  require.NoError(t, err)
  require.NotEmpty(t, hashedPassword)
  
  wrongPassword := RandomString(7)
  err = VerifyPassword(wrongPassword, hashedPassword);
  require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
