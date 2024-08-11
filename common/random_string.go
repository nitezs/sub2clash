package common

import "math/rand"

func RandomString(length int) string {

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, charset[rand.Intn(len(charset))])
	}
	return string(result)
}
