package workspace

import (
	"errors"
	"path/filepath"
)

// Transmission is the data necessary to submit a solution.
type Transmission struct {
	Files   []string
	Dir     string
	ArgDirs []string
}

// NewTransmission processes the arguments to the submit command to prep a submission.
func NewTransmission(root string, args []string) (*Transmission, error) {
	tx := &Transmission{}
	for _, arg := range args {
		pt, err := DetectPathType(arg)
		if err != nil {
			return nil, err
		}
		if pt == TypeFile {
			arg, err = filepath.Abs(arg)
			if err != nil {
				return nil, err
			}
			tx.Files = append(tx.Files, arg)
			continue
		}
		// For our purposes, if it's not a file then it's a directory.
		tx.ArgDirs = append(tx.ArgDirs, arg)
	}
	if len(tx.ArgDirs) > 1 {
		return nil, errors.New("more than one dir")
	}
	if len(tx.ArgDirs) > 0 && len(tx.Files) > 0 {
		return nil, errors.New("mixing files and dirs")
	}
	if len(tx.Files) > 0 {
		ws := New(root)
		parents := map[string]bool{}
		for _, file := range tx.Files {
			dir, err := ws.SolutionDir(file)
			if err != nil {
				return nil, err
			}
			parents[dir] = true
			tx.Dir = dir
		}
		if len(parents) > 1 {
			return nil, errors.New("files are from more than one solution")
		}
	}
	if len(tx.ArgDirs) == 1 {
		tx.Dir = tx.ArgDirs[0]
	}
	return tx, nil
}
