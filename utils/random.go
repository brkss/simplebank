package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generate random number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generate random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// generate random owner name
func RandomOwner() string {
	return RandomString(6)
}

func RandomEmail() string {
  email := fmt.Sprintf("%s@email.com", RandomString(6))
  return email
} 

// generate random account balance
func RandomMoney() int64 {
	return RandomInt(0, 6)
}

// generate random curreny for account
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "MAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
