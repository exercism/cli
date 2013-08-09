package exercism

import (
	"os"
)

func Logout(dir string) {
	os.Remove(configFilename(dir))
}
