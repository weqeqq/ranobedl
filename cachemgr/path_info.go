package cachemgr

type Chapter struct {
	Path   string
	Number string
	Volume string
}

type PathInfo struct {
	Data []Chapter
}

const pathInfoFilename = "PathInfo.json"

func pathInfoIsPresent(RanobeProvider RanobeProvider, uniqueName string) (bool, error) {
	return isPresent(RanobeProvider, uniqueName, pathInfoFilename)
}

func LoadPathInfo(ranobeProvider RanobeProvider, uniqueName string) (PathInfo, error) {
	var pathInfo PathInfo

	return pathInfo, loadJson(
		ranobeProvider,
		uniqueName,
		pathInfoFilename,
		&pathInfo,
	)
}
func (self *PathInfo) Save(ranobeProvider RanobeProvider, uniqueName string) error {
	return SaveJson(
		ranobeProvider,
		uniqueName,
		pathInfoFilename,
		self,
	)
}
