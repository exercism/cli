package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ToFile(path string, c Config) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Your credentials have been written to %s\n", path)
	return nil
}
