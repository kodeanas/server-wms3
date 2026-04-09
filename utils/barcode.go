package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateWarehouseBarcode() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	digits := []rune("0123456789")

	var code []rune
	for i := 0; i < 4; i++ {
		code = append(code, letters[rand.Intn(len(letters))])
	}
	var num []rune
	for i := 0; i < 8; i++ {
		num = append(num, digits[rand.Intn(len(digits))])
	}
	return fmt.Sprintf("LQD%s-%s", string(code), string(num))
}
