package ranobe

import (
	"ranobedl/cachemgr"
	"ranobedl/provider/ranobelib"
)

func Download(provider cachemgr.RanobeProvider, uniqueName string, callback func(current, total int)) error {
	if inCache, err := cachemgr.InCache(provider, uniqueName); err != nil {
		return err
	} else {

		if inCache {
			return nil
		}
	}
	switch provider {

	case cachemgr.RanobeLib:
		return ranobelib.DownloadRanobe(uniqueName, callback)

	}
	return nil
}
