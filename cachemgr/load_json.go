package cachemgr

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func loadJson(ranobeProvider RanobeProvider, unqiueName string, path string, structure any) error {
	ranobeDir, err := ConstructPath(ranobeProvider, unqiueName)

	if err != nil {
		return err
	}
	file, err := os.Open(filepath.Join(ranobeDir, path))

	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(structure); err != nil {
		return err
	} else {
		return nil
	}
}
