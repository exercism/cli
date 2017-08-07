package visibility

import "syscall"

// HideFile sets a Windows file's 'hidden' attribute.
// This is the equivalent of giving a filename on
// Linux or MacOS a leading dot (e.g. .bash_rc).
func HideFile(path string) error {
	// This is based on the discussion in
	// https://www.reddit.com/r/golang/comments/5t3ezd/hidden_files_directories/
	// but instead of duplicating all the effort to write the file, this takes
	// the path of a written file and then flips the bit on the relevant attribute.
	// The attributes are a bitmask (uint32), so we can't call
	// SetFileAttributes(ptr, syscall.File_ATTRIBUTE_HIDDEN) as suggested, since
	// that would wipe out any existing attributes.
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	attributes, err := syscall.GetFileAttributes(ptr)
	if err != nil {
		return err
	}

	attributes |= syscall.FILE_ATTRIBUTE_HIDDEN
	return syscall.SetFileAttributes(ptr, attributes)
}
