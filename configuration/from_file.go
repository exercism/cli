package configuration

import (
	"io/ioutil"
	"encoding/json"
)

func FromFile(dir string) (c Config, err error) {
	bytes, err := ioutil.ReadFile(Filename(dir))
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c)
	return
}

