// +build !windows

package visibility

// HideFile is a no-op for non-Windows systems.
func HideFile(string) error {
	return nil
}

// ShowFile is a no-op for non-Windows systems.
func ShowFile(path string) error {
	return nil
}
