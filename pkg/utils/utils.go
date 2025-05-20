package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func ParseStringArray(input string) []string {
	// Удаляем квадратные скобки
	input = strings.TrimPrefix(input, "[")
	input = strings.TrimSuffix(input, "]")

	// Разделяем строку по запятым
	items := strings.Split(input, ",")

	// Удаляем кавычки и лишние пробелы у каждого элемента
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		item = strings.Trim(item, `"`)
		item = strings.Trim(item, `'`)
		if item != "" {
			result = append(result, item)
		}
	}

	return result
}

// GenerateHashString generate hash string
func GenerateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
