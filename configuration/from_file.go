package configuration

import (
	"encoding/json"
	"io/ioutil"
)

func FromFile(path string) (c Config, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c)
	return
}
