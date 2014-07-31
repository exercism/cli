package config

import (
	"fmt"
)

const FILENAME = ".exercism.go"

func Filename(dir string) string {
	return fmt.Sprintf("%s/%s", dir, FILENAME)
}
