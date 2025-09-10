package cachemgr

import (
	"os"
	"path/filepath"
)

func ClearCache() error {
	if cacheDir, err := os.UserCacheDir(); err != nil {
		return err
	} else {
		return os.RemoveAll(filepath.Join(cacheDir, "ranobedl"))
	}
}
