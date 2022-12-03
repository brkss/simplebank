package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error){
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", fmt.Errorf("faild to hash password : %w ", err)
  }
  return string(hashedPassword), nil 
}

func VerifyPassword(password string, hash string) error {
  return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
