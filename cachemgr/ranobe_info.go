package cachemgr

type RanobeInfo struct {
	Name   string
	Author string
}

const ranobeInfoFilename string = "RanobeInfo.json"

func ranobeInfoIsPresent(ranobeProvider RanobeProvider, uniqueName string) (bool, error) {
	return isPresent(ranobeProvider, uniqueName, ranobeInfoFilename)
}
func LoadRanobeInfo(ranobeProvider RanobeProvider, uniqueName string) (RanobeInfo, error) {
	var ranobeInfo RanobeInfo

	return ranobeInfo, loadJson(
		ranobeProvider,
		uniqueName,
		ranobeInfoFilename,
		&ranobeInfo,
	)
}
func (self *RanobeInfo) Save(ranobeProvider RanobeProvider, uniqueName string) error {
	return SaveJson(
		ranobeProvider,
		uniqueName,
		ranobeInfoFilename,
		self,
	)
}
