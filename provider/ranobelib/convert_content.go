package ranobelib

import (
	"fmt"
	"path"
	api "ranobedl/api/ranobelib"
	"ranobedl/cachemgr"
	"ranobedl/schema"
	"strings"
)

type contentConvertor struct {
	UniqueName string

	Data api.ChapterContentData
}

func (cc *contentConvertor) fromHtml() (schema.Node, error) {
	if htmlContent, err := cc.Data.HtmlContent(); err != nil {
		return schema.Node{}, err
	} else {
		if htmlContent[len("<div>"):] == "<div>" {
			htmlContent = htmlContent[len("<div>") : len(htmlContent)-len("/<div>")]
		}
		node, err := schema.FromHtmlString(htmlContent)
		if err != nil {
			return schema.Node{}, err
		}
		for _, child := range node.Content {
			if child.Type == schema.NodeTypeImage {
				src, err := child.ImageSrc()
				if err != nil {
					return schema.Node{}, err
				}

				filename := path.Base(src)
				part := strings.TrimSuffix(filename, path.Ext(filename))
				child.Attrs["src"] = part
			}
		}
		return node, nil
	}
}
func (cc *contentConvertor) convertSchema(mlSchema schema.Node) (schema.Node, error) {
	output := mlSchema

	for index, child := range output.Content {
		if child.Type == schema.NodeTypeImage {
			output.Content[index] = schema.Node{
				Type:  schema.NodeTypeImage,
				Attrs: map[string]any{"src": child.Attrs["images"].([]any)[0].(map[string]any)["image"].(string)},
			}
		}
	}
	return output, nil
}
func (cc *contentConvertor) fromSchema() (schema.Node, error) {
	if schemaContent, err := cc.Data.SchemaContent(); err != nil {
		return schema.Node{}, err
	} else {
		return cc.convertSchema(schemaContent)
	}
}
func (cc *contentConvertor) downloadImage(index int, source string) (string, error) {
	for _, attachment := range cc.Data.Attachments {
		if attachment.Name == source {

			filename := fmt.Sprintf(
				"%s%simage%d.%s",
				cc.Data.Volume,
				cc.Data.Number,
				index,
				attachment.Extension,
			)
			return cachemgr.DownloadImage(
				provider,
				cc.UniqueName,
				"https://ranobelib.me"+attachment.Url,
				filename,
			)
		}
	}
	return "", fmt.Errorf("Image not found")
}
func (cc *contentConvertor) replaceImgSrc(node schema.Node) error {
	index := 0

	for _, child := range node.Content {
		if child.Type == schema.NodeTypeImage {
			if src, err := child.ImageSrc(); err != nil {
				return err
			} else {

				if path, err := cc.downloadImage(index, src); err != nil {
					return err

				} else {
					child.Attrs["src"] = path
					index++
				}
			}
		}
	}
	return nil
}
func (cc *contentConvertor) Convert() (schema.Node, error) {
	var output schema.Node

	if node, err := cc.fromHtml(); err != nil {

		if node, err := cc.fromSchema(); err != nil {
			return output, err
		} else {
			output = node
		}

	} else {
		output = node
	}
	if err := cc.replaceImgSrc(output); err != nil {
		return output, err
	} else {
		return output, nil
	}
}
func convertContent(uniqueName string, data api.ChapterContentData) (schema.Node, error) {
	return (&contentConvertor{uniqueName, data}).Convert()
}
