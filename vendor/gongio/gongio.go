package gongio

import (
	"io/ioutil"
	"os"
	"errors"
	"fmt"
)

func WriteBytesToFile(filename string, buf []byte) {
	err := ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		panic(err)
	}
}

func Mkdir(filepath string) error {
	if fi, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.Mkdir(filepath, 0755); err != nil {
			return err
		}
	} else if !(fi.Mode().IsDir()) {
		panicText := fmt.Sprintf("%v path already exist", filepath)
		return errors.New(panicText)
	}
	return nil
}
