package ranobelib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type author struct {
	Name string `json:"name"`
}
type ranobeInfoData struct {
	Name    string   `json:"name"`
	Authors []author `json:"authors"`
}
type ranobeInfo struct {
	UniqueName string
}

func (self *ranobeInfo) constructUrl() string {
	return fmt.Sprintf("%s/manga/%s?fields[]=authors", apiUrl, self.UniqueName)
}
func (self *ranobeInfo) Parse() (ranobeInfoData, error) {
	output := struct {
		Data ranobeInfoData `json:"data"`
	}{}
	if response, err := http.Get(self.constructUrl()); err != nil {
		return output.Data, err
	} else {
		defer response.Body.Close()

		if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
			return output.Data, err
		} else {
			return output.Data, nil
		}
	}
}
func GetRanobeInfo(uniqueName string) (ranobeInfoData, error) {
	ranobeInfo := ranobeInfo{UniqueName: uniqueName}
	return ranobeInfo.Parse()
}
