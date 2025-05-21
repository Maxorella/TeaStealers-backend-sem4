package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func ParseStringArray(input string) []string {
	input = strings.TrimPrefix(input, "[")
	input = strings.TrimSuffix(input, "]")

	items := strings.Split(input, ",")

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

func GenerateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
