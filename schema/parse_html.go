package schema

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type htmlParser struct {
	HtmlNode *html.Node
}

func newHtmlParser(htmlNode *html.Node) *htmlParser {
	return &htmlParser{htmlNode}
}

func (hp *htmlParser) findAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
func (hp *htmlParser) parseBlockChildren(node *html.Node) ([]Node, error) {
	var output []Node
	var paragraphContent []Node

	pContentIsEmpty := func() bool {
		return len(paragraphContent) == 0
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		switch child.Type {
		case html.TextNode:
			{
				if strings.TrimSpace(child.Data) != "" {
					paragraphContent = append(
						paragraphContent,
						Node{
							Type: NodeTypeText,
							Text: child.Data,
						},
					)
				}
			}
		case html.ElementNode:
			{
				nodeType, err := NodeTypeFromHTML(child.Data)
				if err != nil {
					return []Node{}, err
				}
				switch nodeType.Group() {
				case NodeGroupBlock:
					{

						if !pContentIsEmpty() {
							output = append(output, Node{
								Type:    NodeTypeParagraph,
								Content: paragraphContent,
							})
							paragraphContent = nil
						}
						if blockNode, err := hp.parseBlock(child); err != nil {
							return []Node{}, err
						} else {
							output = append(output, blockNode)
						}
					}
				case NodeGroupInline:
					{

						if nodeType == NodeTypeImage {
							output = append(output, Node{
								Type:  NodeTypeImage,
								Attrs: map[string]any{"src": hp.findAttr(child, "src")},
							})
						} else if nodeType == NodeTypeHardBreak {
							output = append(output, Node{
								Type: NodeTypeHardBreak,
							})
						} else if inlineNodes, err := hp.parseInline(child, []Mark{}); err != nil {
							return []Node{}, err
						} else {
							paragraphContent = append(paragraphContent, inlineNodes...)
						}
					}
				}
			}
		}
	}
	if !pContentIsEmpty() {
		output = append(output, Node{
			Type:    NodeTypeParagraph,
			Content: paragraphContent,
		})
	}
	return output, nil
}
func (hp *htmlParser) handleParagraph(node *html.Node) (Node, error) {
	if content, err := hp.parseInlineChildren(node, []Mark{}); err != nil {
		return Node{}, err
	} else {
		return Node{
			Type:    NodeTypeParagraph,
			Content: content,
		}, nil
	}
}
func (hp *htmlParser) handleHeading(node *html.Node) (Node, error) {
	level, _ := strconv.Atoi(node.Data[1:])

	if content, err := hp.parseInlineChildren(node, []Mark{}); err != nil {
		return Node{}, err
	} else {
		return Node{
			Type:    NodeTypeHeading,
			Attrs:   map[string]any{"level": level},
			Content: content,
		}, nil
	}
}
func (hp *htmlParser) handleBulletList(node *html.Node) (Node, error) {
	if content, err := hp.parseBlockChildren(node); err != nil {
		return Node{}, err
	} else {
		return Node{
			Type:    NodeTypeBulletList,
			Content: content,
		}, nil
	}
}
func (hp *htmlParser) handleOrderedList(node *html.Node) (Node, error) {
	if content, err := hp.parseBlockChildren(node); err != nil {
		return Node{}, err
	} else {
		return Node{
			Type:    NodeTypeOrderedList,
			Content: content,
		}, nil
	}
}
func (hp *htmlParser) handleListItem(node *html.Node) (Node, error) {
	if content, err := hp.parseBlockChildren(node); err != nil {
		return Node{}, err
	} else {
		return Node{
			Type:    NodeTypeListItem,
			Content: content,
		}, nil
	}
}
func (hp *htmlParser) handleBlockquote(node *html.Node) (Node, error) {
	if content, err := hp.parseInlineChildren(node, []Mark{}); err != nil {
		return Node{}, err
	} else {
		return Node{
			Type:    NodeTypeBlockquote,
			Content: content,
		}, nil
	}
}
func (hp *htmlParser) handleCodeBlock(node *html.Node) (Node, error) {
	if content, err := hp.parseInlineChildren(node.FirstChild, []Mark{}); err != nil {
		return Node{}, err
	} else {
		return Node{
			Type:    NodeTypeCodeBlock,
			Content: content,
		}, nil
	}
}
func (hp *htmlParser) handleHorizontalRule(_ *html.Node) (Node, error) {
	return Node{
		Type: NodeTypeHorizontalRule,
	}, nil
}
func (hp *htmlParser) handleHardBreak(_ *html.Node) (Node, error) {
	return Node{
		Type: NodeTypeHardBreak,
	}, nil
}
func (hp *htmlParser) parseBlock(node *html.Node) (Node, error) {
	tag := node.Data

	if isParagraph(tag) {
		return hp.handleParagraph(node)
	}
	if isHeading(tag) {
		return hp.handleHeading(node)
	}
	if isBulletList(tag) {
		return hp.handleBulletList(node)
	}
	if isOrderedList(tag) {
		return hp.handleOrderedList(node)
	}
	if isListItem(tag) {
		return hp.handleListItem(node)
	}
	if isBlockquote(tag) {
		return hp.handleBlockquote(node)
	}
	if isCodeBlock(tag) {
		return hp.handleCodeBlock(node)
	}
	if isHorizontalRule(tag) {
		return hp.handleHorizontalRule(node)
	}
	if isHardBreak(tag) {
		return hp.handleHardBreak(node)
	}
	return Node{}, fmt.Errorf("Undefined block tag, %s", tag)
}

func (hp *htmlParser) parseInlineChildren(node *html.Node, marks []Mark) ([]Node, error) {
	var output []Node

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if content, err := hp.parseInline(child, marks); err != nil {
			return []Node{}, err
		} else {
			output = append(output, content...)
		}
	}
	return output, nil
}
func (hp *htmlParser) parseInline(node *html.Node, marks []Mark) ([]Node, error) {
	switch node.Type {
	case html.TextNode:
		{
			if strings.TrimSpace(node.Data) != "" {
				return []Node{{
					Type:  NodeTypeText,
					Text:  node.Data,
					Marks: marks,
				}}, nil
			}
			return []Node{}, nil
		}
	case html.ElementNode:
		{
			var tag string = node.Data
			var mark Mark

			if isBold(tag) {
				mark.Type = MarkTypeBold

			} else if isItalic(tag) {
				mark.Type = MarkTypeItalic

			} else if isUnderline(tag) {
				mark.Type = MarkTypeUnderline

			} else if isStrike(tag) {
				mark.Type = MarkTypeStrike

			} else if isCode(tag) {
				mark.Type = MarkTypeCode

			} else if isLink(tag) {
				mark.Type = MarkTypeLink
				mark.Attrs = map[string]any{"href": hp.findAttr(node, "href")}
			} else {
				return []Node{}, fmt.Errorf("Undefined inline tag, %s", tag)
			}
			return hp.parseInlineChildren(node, append(marks, mark))
		}
	default:
		return []Node{}, fmt.Errorf("Unsupported HTML Node type: %v", node.Type)
	}
}
func (hp *htmlParser) findBody(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "body" {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if body := hp.findBody(child); body != nil {
			return body
		}
	}
	return nil
}
func (hp *htmlParser) Parse() (Node, error) {
	if content, err := hp.parseBlockChildren(hp.findBody(hp.HtmlNode)); err != nil {
		return Node{}, err
	} else {
		return Node{
			Type:    NodeTypeDoc,
			Content: content,
		}, nil
	}
}
func parseHtml(node *html.Node) (Node, error) {
	return newHtmlParser(node).Parse()
}
func FromHtmlString(htmlString string) (Node, error) {
	if node, err := html.Parse(strings.NewReader(htmlString)); err != nil {
		return Node{}, err
	} else {
		return parseHtml(node)
	}
}
func FromHtmlStream(stream io.Reader) (Node, error) {
	if node, err := html.Parse(stream); err != nil {
		return Node{}, err
	} else {
		return parseHtml(node)
	}
}
func FromHtmlFile(filename string) (Node, error) {
	if file, err := os.Open(filename); err != nil {
		return Node{}, err
	} else {
		defer file.Close()
		return FromHtmlStream(file)
	}
}
