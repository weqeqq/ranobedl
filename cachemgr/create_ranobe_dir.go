package cachemgr

import (
	"os"
)

func CreateRanobeDir(ranobeProvider RanobeProvider, uniqueName string) error {
	if ranobeDir, err := ConstructPath(ranobeProvider, uniqueName); err != nil {
		return err
	} else {
		os.MkdirAll(ranobeDir, 0777)
		return nil
	}
}
