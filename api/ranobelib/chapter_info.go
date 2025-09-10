package ranobelib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type user struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
}
type branch struct {
	Id          int    `json:"id"`
	BranchId    int    `json:"branch_id"`
	CreatedAt   string `json:"created_at"`
	Teams       []team `json:"teams"`
	ExpiredType int    `json:"expired_type"`
	User        user   `json:"user"`
}
type chapterInfoData struct {
	Id              int      `json:"id"`
	Index           int      `json:"index"`
	ItemNumber      int      `json:"item_number"`
	Volume          string   `json:"volume"`
	Number          string   `json:"number"`
	NumberSecondary string   `json:"number_secondary"`
	Name            string   `json:"name"`
	BranchesCount   int      `json:"branches_count"`
	Branches        []branch `json:"branches"`
}
type chapterInfo struct {
	uniqueName string
}

func (self *chapterInfo) constructUrl() string {
	return fmt.Sprintf("%s/manga/%s/chapters", apiUrl, self.uniqueName)
}
func (self *chapterInfo) Parse() ([]chapterInfoData, error) {
	output := struct {
		Data []chapterInfoData `json:"data"`
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
func GetChapterInfo(uniqueName string) ([]chapterInfoData, error) {
	chapterInfo := chapterInfo{uniqueName: uniqueName}
	return chapterInfo.Parse()
}
