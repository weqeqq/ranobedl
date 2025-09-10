package cachemgr

import (
	"os"
	"path/filepath"
)

type RanobeProvider int

const (
	RanobeLib RanobeProvider = iota
	RanobeHub
)

func (self *RanobeProvider) String() string {
	switch *self {
	case RanobeLib:
		return "ranobelib"
	case RanobeHub:
		return "ranobehub"
	default:
		return ""
	}
}
func ConstructPath(ranobeProvider RanobeProvider, uniqueName string) (string, error) {
	if cacheDir, err := os.UserCacheDir(); err != nil {
		return "", err
	} else {
		return filepath.Join(cacheDir, "ranobedl", ranobeProvider.String(), uniqueName), nil
	}
}
