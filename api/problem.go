package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Problem struct {
	ID       string            `json:"id"`
	TrackID  string            `json:"track_id"`
	Language string            `json:"language"`
	Slug     string            `json:"slug"`
	Name     string            `json:"name"`
	IsFresh  bool              `json:"fresh"`
	Files    map[string]string `json:"files"`
}

func (p *Problem) String() string {
	return fmt.Sprintf("%s - %s in %s", p.ID, p.Name, p.Language)
}

func (p *Problem) ExistsIn(dir string) bool {
	path := fmt.Sprintf("%s/%s", dir, p.ID)

	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func (p *Problem) Save(dir string) error {
	if p.ExistsIn(dir) {
		return nil
	}

	for name, text := range p.Files {
		file := fmt.Sprintf("%s/%s/%s", dir, p.ID, name)

		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(file, []byte(text), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
