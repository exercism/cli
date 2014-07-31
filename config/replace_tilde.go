package config

import (
	"strings"
)

func ReplaceTilde(oldPath string) string {
	return strings.Replace(oldPath, "~/", HomeDir()+"/", 1)
}
