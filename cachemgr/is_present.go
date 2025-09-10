package cachemgr

import (
	"os"
	"path/filepath"
)

func isPresent(ranobeProvider RanobeProvider, uniqueName string, filename string) (bool, error) {
	if ranobeDir, err := ConstructPath(ranobeProvider, uniqueName); err != nil {
		return false, err
	} else {
		if _, err := os.Stat(filepath.Join(ranobeDir, filename)); err != nil {
			if os.IsNotExist(err) {
				return false, nil
			} else {
				return false, err
			}
		} else {
			return true, nil
		}
	}
}
