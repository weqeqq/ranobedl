package fb2

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type builder struct {
	document       document
	currentSection *section
}
type document struct {
	XMLName          xml.Name    `xml:"FictionBook"`
	XMLNamespace     string      `xml:"xmlns,attr"`
	XMLNamespaceLink string      `xml:"xmlns:l,attr"`
	Description      description `xml:"description"`
	Body             body        `xml:"body"`
	Binary           []binary    `xml:"binary"`
}
type description struct {
	TitleInfo titleInfo `xml:"title-info"`
}
type titleInfo struct {
	BookTitle string   `xml:"book-title"`
	Author    []author `xml:"author"`
	Date      string   `xml:"date"`
}
type author struct {
	FirstName string `xml:"first-name"`
	LastName  string `xml:"last-name"`
}
type body struct {
	Sections []section `xml:"section"`
}
type section struct {
	Title      title       `xml:"title"`
	Paragraphs []paragraph `xml:"p"`
}
type title struct {
	Paragraph string `xml:"p"`
}
type paragraph struct {
	Text string `xml:",innerxml"`
}
type binary struct {
	ID          string `xml:"id,attr"`
	ContentType string `xml:"content-type,attr"`
	Data        string `xml:",chardata"`
}

func NewBuilder() *builder {
	fb2 := document{
		Description: description{
			TitleInfo: titleInfo{
				BookTitle: "",
				Author: []author{{
					FirstName: "",
					LastName:  "",
				}},
				Date: time.Now().Format("2006-01-02"),
			},
		},
		Body: body{
			Sections: []section{},
		},
	}
	return &builder{
		document: fb2,
	}
}
func (self *builder) SetTitle(name string) {
	self.document.Description.TitleInfo.BookTitle = name
}
func (self *builder) SetAuthor(author string) {
}
func (self *builder) PushChapter(chapterTitle string) error {
	section := section{
		Title: title{Paragraph: chapterTitle},
	}
	self.document.Body.Sections = append(self.document.Body.Sections, section)
	self.currentSection = &self.document.Body.Sections[len(self.document.Body.Sections)-1]
	return nil
}
func (self *builder) PushParagraph(text string) error {
	if self.currentSection == nil {
		return errors.New("Chapter is not created")
	}
	paragraph := paragraph{
		Text: text,
	}
	self.currentSection.Paragraphs = append(self.currentSection.Paragraphs, paragraph)
	return nil
}
func (self *builder) PushImage(imagePath string) error {
	if self.currentSection == nil {
		return errors.New("Chapter is not created")
	}

	filename := filepath.Base(imagePath)
	imageID := strings.TrimSuffix(filename, path.Ext(filename))

	data, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}
	base64Data := base64.StdEncoding.EncodeToString(data)

	var contentType string
	if strings.HasSuffix(strings.ToLower(imagePath), ".png") {
		contentType = "image/png"
	} else {
		contentType = "image/jpeg"
	}
	self.PushParagraph(fmt.Sprintf("<image l:href=\"#%s\"/>", imageID))

	binary := binary{
		ID:          imageID,
		ContentType: contentType,
		Data:        base64Data,
	}
	self.document.Binary = append(self.document.Binary, binary)
	return nil
}
func (c *builder) Build(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	c.document.XMLNamespace = "http://www.gribuser.ru/xml/fictionbook/2.0"
	c.document.XMLNamespaceLink = "http://www.w3.org/1999/xlink"

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(c.document); err != nil {
		return err
	}
	return nil
}
