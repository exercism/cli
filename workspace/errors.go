package workspace

import "fmt"

// ErrNotInWorkspace signals that the target directory is outside the configured workspace.
type ErrNotInWorkspace string

// ErrNotExist signals that the target directory could not be located.
type ErrNotExist string

func (err ErrNotInWorkspace) Error() string {
	return fmt.Sprintf("%s not within workspace", string(err))
}

func (err ErrNotExist) Error() string {
	return fmt.Sprintf("%s not found", string(err))
}

// IsNotInWorkspace checks if this is an ErrNotInWorkspace error.
func IsNotInWorkspace(err error) bool {
	_, ok := err.(ErrNotInWorkspace)
	return ok
}

// IsNotExist checks if this is an ErrNotExist error.
func IsNotExist(err error) bool {
	_, ok := err.(ErrNotExist)
	return ok
}
