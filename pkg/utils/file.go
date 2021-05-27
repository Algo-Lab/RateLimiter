package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// WriteFileSafety trys to over write a file safety.
func WriteFileSafety(filename string, data []byte, perm os.FileMode) (err error) {
	tempFile := filename + ".tmp"
Try:
	for i := 0; i < 5; i++ {
		err = ioutil.WriteFile(tempFile, data, perm)
		if err == nil {
			break Try
		}
	}
	if err != nil {
		return err
	}
	err = os.Rename(tempFile, filename)
	return
}

const JsonExt = ".json"

var ErrIgnore = errors.New("error ignore")

// ReadJsonFile reads a json into interface
// If a file is not .json ignore it
// If a file is empty, ignore it
func ReadJsonFile(file string, v interface{}) error {
	if path.Ext(file) != JsonExt {
		return ErrIgnore
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read file %s got error %v", file, err)
	}
	if len(b) == 0 {
		return ErrIgnore
	}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("unmarshal json file %s got error %v", file, err)
	}
	return nil
}
