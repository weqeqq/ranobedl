package schema

import (
	"encoding/json"
	"fmt"
)

type NodeType int

const (
	NodeTypeDoc NodeType = iota
	NodeTypeParagraph
	NodeTypeHeading
	NodeTypeBulletList
	NodeTypeOrderedList
	NodeTypeListItem
	NodeTypeBlockquote
	NodeTypeCodeBlock
	NodeTypeHorizontalRule

	// inline

	NodeTypeText
	NodeTypeHardBreak
	NodeTypeImage
)

type NodeGroup int

const (
	NodeGroupBlock NodeGroup = iota
	NodeGroupInline
)

func (nt NodeType) IsBlock() bool {
	return nt.Group() == NodeGroupBlock
}
func (nt NodeType) IsInline() bool {
	return nt.Group() == NodeGroupInline
}
func (nt NodeType) Group() NodeGroup {
	if nt == NodeTypeDoc ||
		nt == NodeTypeParagraph ||
		nt == NodeTypeHeading ||
		nt == NodeTypeBulletList ||
		nt == NodeTypeOrderedList ||
		nt == NodeTypeListItem ||
		nt == NodeTypeBlockquote ||
		nt == NodeTypeCodeBlock ||
		nt == NodeTypeHorizontalRule {
		return NodeGroupBlock
	}
	if nt == NodeTypeText ||
		nt == NodeTypeHardBreak ||
		nt == NodeTypeImage {
		return NodeGroupInline
	}
	panic(fmt.Sprintf("Undefined NodeGroup: %d", nt))
}
func NodeTypeFromString(str string) (NodeType, error) {
	switch str {
	case "doc":
		return NodeTypeDoc, nil
	case "paragraph":
		return NodeTypeParagraph, nil
	case "heading":
		return NodeTypeHeading, nil
	case "bulletList":
		return NodeTypeBulletList, nil
	case "orderedList":
		return NodeTypeOrderedList, nil
	case "listItem":
		return NodeTypeListItem, nil
	case "blockquote":
		return NodeTypeBlockquote, nil
	case "codeBlock":
		return NodeTypeCodeBlock, nil
	case "horizontalRule":
		return NodeTypeHorizontalRule, nil
	case "text":
		return NodeTypeText, nil
	case "hardBreak":
		return NodeTypeHardBreak, nil
	case "image":
		return NodeTypeImage, nil
	default:
		return -1, fmt.Errorf("Undefined NodeType: %s", str)
	}
}
func (nt *NodeType) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)

	if nodeType, err := NodeTypeFromString(str); err != nil {
		*nt = -1
		return err
	} else {
		*nt = nodeType
		return nil
	}
}
func (self NodeType) String() string {
	switch self {
	case NodeTypeDoc:
		return "doc"
	case NodeTypeParagraph:
		return "paragraph"
	case NodeTypeHeading:
		return "heading"
	case NodeTypeBulletList:
		return "bulletList"
	case NodeTypeOrderedList:
		return "orderedList"
	case NodeTypeListItem:
		return "listItem"
	case NodeTypeBlockquote:
		return "blockquote"
	case NodeTypeCodeBlock:
		return "codeBlock"
	case NodeTypeHorizontalRule:
		return "horizontalRule"
	case NodeTypeText:
		return "text"
	case NodeTypeHardBreak:
		return "hardBreak"
	case NodeTypeImage:
		return "image"
	default:
		panic(fmt.Sprintf("Undefined NodeType: %d", self))
	}
}
func (nt NodeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(nt.String())
}

func isParagraph(html string) bool {
	return html == "p"
}
func isHeading(html string) bool {
	return html == "h1" ||
		html == "h2" ||
		html == "h3" ||
		html == "h4" ||
		html == "h5" ||
		html == "h6"
}
func isBulletList(html string) bool {
	return html == "ul"
}
func isOrderedList(html string) bool {
	return html == "ol"
}
func isListItem(html string) bool {
	return html == "li"
}
func isBlockquote(html string) bool {
	return html == "blockquote"
}
func isCodeBlock(html string) bool {
	return html == "pre"
}
func isHorizontalRule(html string) bool {
	return html == "hr"
}
func isHardBreak(html string) bool {
	return html == "br"
}
func isImage(html string) bool {
	return html == "img"
}
func NodeTypeFromHTML(html string) (NodeType, error) {
	if isParagraph(html) {
		return NodeTypeParagraph, nil
	}
	if isHeading(html) {
		return NodeTypeHeading, nil
	}
	if isBulletList(html) {
		return NodeTypeBulletList, nil
	}
	if isOrderedList(html) {
		return NodeTypeOrderedList, nil
	}
	if isListItem(html) {
		return NodeTypeListItem, nil
	}
	if isBlockquote(html) {
		return NodeTypeBlockquote, nil
	}
	if isCodeBlock(html) {
		return NodeTypeCodeBlock, nil
	}
	if isHorizontalRule(html) {
		return NodeTypeHorizontalRule, nil
	}
	if isHardBreak(html) {
		return NodeTypeHardBreak, nil
	}
	if isImage(html) {
		return NodeTypeImage, nil
	}
	return -1, fmt.Errorf("Undefined HtmlTag: %s", html)
}
