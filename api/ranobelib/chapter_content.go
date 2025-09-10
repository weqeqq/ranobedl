package ranobelib

import (
	"encoding/json"
	"fmt"
	"ranobedl/schema"
	"ranobedl/util"
)

type Attachment struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Url       string `json:"url"`
}
type ChapterContentData struct {
	Volume      string          `json:"volume"`
	Number      string          `json:"number"`
	Name        string          `json:"name"`
	Content     json.RawMessage `json:"content"`
	Attachments []Attachment    `json:"attachments"`
}

func (self *ChapterContentData) HtmlContent() (string, error) {
	var output string
	return output, json.Unmarshal(self.Content, &output)
}
func (self *ChapterContentData) SchemaContent() (schema.Node, error) {
	var output schema.Node
	return output, json.Unmarshal(self.Content, &output)
}

type chapterContent struct {
	UniqueName string
	Number     string
	Volume     string
}

func (self *chapterContent) constructUrl() string {
	return fmt.Sprintf(
		"%s/manga/%s/chapter?number=%s&volume=%s",
		apiUrl,
		self.UniqueName,
		self.Number,
		self.Volume,
	)
}
func (self *chapterContent) Parse() (ChapterContentData, error) {
	output := struct {
		Data ChapterContentData `json:"data"`
	}{}
	if response, err := util.SendRequest(self.constructUrl()); err != nil {
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
func GetChapterContent(uniqueName string, number string, volume string) (ChapterContentData, error) {
	return (&chapterContent{
		UniqueName: uniqueName,
		Number:     number,
		Volume:     volume,
	}).Parse()
}
