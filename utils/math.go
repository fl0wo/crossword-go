package utils

import "math/rand"

func RandInt(i int, i2 int) int {
	return rand.Intn(i2-i) + i
}
