package cachemgr

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveJson(ranobeProvider RanobeProvider, uniqueName string, path string, structure any) error {
	ranobeDir, err := ConstructPath(ranobeProvider, uniqueName)

	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(ranobeDir, path))

	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(structure); err != nil {
		return err
	} else {
		return nil
	}
}
