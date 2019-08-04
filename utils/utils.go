package utils

import "strings"

// Redact masks the given token by replacing part of the string with *
func Redact(token string) string {
	str := token[4 : len(token)-3]
	redaction := strings.Repeat("*", len(str))
	return string(token[:4]) + redaction + string(token[len(token)-3:])
}
