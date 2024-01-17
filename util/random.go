package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Randommobile generates a random phone number
func RandomMobile(n int) int {
	return rand.Intn(n)
}

const alphabet = "abcdefghijklmnopqrstuvwz"

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	alphabetLen := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(alphabetLen)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomUser generates a random user name
func RandomUser() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomUser())
}
