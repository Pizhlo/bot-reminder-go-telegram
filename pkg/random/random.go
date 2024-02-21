package random

import (
	"math/rand"
	"strconv"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890123456789"

func Int(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func String(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func Phone() string {
	phoneStr := "+7"

	phone := Int(9000000000, 9999999999)

	phoneStr += strconv.Itoa(phone)

	return phoneStr
}

func Email(n int) string {
	login := String(n)

	login += "@mail.ru"

	return login
}
