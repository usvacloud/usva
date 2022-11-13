package utils

import "strings"

func StandardTypedToBytes(typed string) {
	trimmed := strings.ReplaceAll(typed, " ", "")
	trimmedLower := strings.ToLower(trimmed)
	_ = trimmedLower
}
