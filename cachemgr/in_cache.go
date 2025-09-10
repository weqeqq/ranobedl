package cachemgr

func InCache(ranobeProvider RanobeProvider, uniqueName string) (bool, error) {

	ranobeInfo, err := ranobeInfoIsPresent(ranobeProvider, uniqueName)
	if err != nil {
		return false, err
	}

	pathInfo, err := pathInfoIsPresent(ranobeProvider, uniqueName)
	if err != nil {
		return false, err
	}
	return ranobeInfo && pathInfo, nil
}
