package utils

import (
	"strconv"
	"strings"
	"unicode"
)

// ParseIntDefault tries to parse string to int, returns 0 if fail
// Menghapus koma/titik ribuan dan spasi
func ParseIntDefault(s string) int {
	cleaned := cleanNumberString(s)
	i, _ := strconv.Atoi(cleaned)
	return i
}

// ParseFloatDefault tries to parse string to float64, returns 0 if fail
// Menghapus koma/titik ribuan dan spasi
func ParseFloatDefault(s string) float64 {
	cleaned := cleanNumberString(s)
	f, _ := strconv.ParseFloat(cleaned, 64)
	return f
}

// cleanNumberString menghapus semua koma, titik, dan spasi kecuali satu titik desimal terakhir
func cleanNumberString(s string) string {
	s = strings.TrimSpace(s)
	// Hapus semua spasi
	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
	// Hapus semua koma dan titik, kecuali titik terakhir (anggap sebagai desimal)
	var result []rune
	dotSeen := false
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' && !dotSeen {
			dotSeen = true
			result = append(result, '.')
		} else if s[i] == ',' || s[i] == '.' {
			// skip
		} else {
			result = append(result, rune(s[i]))
		}
	}
	// Balikkan hasilnya
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}
