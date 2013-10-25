package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ToFile(dir string, c Config) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	filename := Filename(dir)
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Your credentials have been written to %s\n", filename)
	return nil
}
