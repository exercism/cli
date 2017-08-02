package browser

import (
	"os/exec"
	"runtime"
	"strings"
)

// Open opens a browser to the given URL.
// The terminal's open command is operating system dependent.
func Open(url string) error {
	// Escape characters are not allowed by cmd/bash.
	switch runtime.GOOS {
	case "windows":
		url = strings.Replace(url, "&", `^&`, -1)
	default:
		url = strings.Replace(url, "&", `\&`, -1)
	}

	// The command to open the browser is OS-dependent.
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "freebsd", "linux", "netbsd", "openbsd":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}

	return cmd.Run()
}
