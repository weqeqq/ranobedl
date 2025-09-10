package noderenderer

import (
	"fmt"
	"html"
	"ranobedl/schema"
)

type textRenderer struct {
	schema.Node
}

func newTextRenderer(node schema.Node) *textRenderer {
	return &textRenderer{node}
}

func (tr *textRenderer) handleBold(text string) (string, error) {
	return fmt.Sprintf("<strong>%s</strong>", text), nil
}
func (tr *textRenderer) handleItalic(text string) (string, error) {
	return fmt.Sprintf("<emphasis>%s</emphasis>", text), nil
}
func (tr *textRenderer) handleUnderline(text string) (string, error) {
	return fmt.Sprintf("<span class=\"underline\">%s</span>", text), nil
}
func (tr *textRenderer) handleStrike(text string) (string, error) {
	return fmt.Sprintf("<span class=\"strike\">%s</span>", text), nil
}
func (tr *textRenderer) handleCode(text string) (string, error) {
	return fmt.Sprintf("<code>%s</code>", text), nil
}
func (tr *textRenderer) handleLink(text string, mark schema.Mark) (string, error) {
	if href, err := mark.LinkHref(); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("<a href=\"%s\">%s</a>", text, href), nil
	}
}
func (tr *textRenderer) renderMark(text string, mark schema.Mark) (string, error) {
	text = html.EscapeString(text)

	switch mark.Type {
	case schema.MarkTypeBold:
		return tr.handleBold(text)
	case schema.MarkTypeItalic:
		return tr.handleItalic(text)
	case schema.MarkTypeUnderline:
		return tr.handleUnderline(text)
	case schema.MarkTypeStrike:
		return tr.handleStrike(text)
	case schema.MarkTypeCode:
		return tr.handleCode(text)
	case schema.MarkTypeLink:
		return tr.handleLink(text, mark)

	default:
		panic(fmt.Sprintf("Undefined MarkType: %d", mark.Type))
	}
}
func (tr *textRenderer) Render() (string, error) {
	if tr.Node.Type != schema.NodeTypeText {
		return "", fmt.Errorf("Expected text node, but got %v", tr.Node.Type)
	}
	output := html.EscapeString(tr.Node.Text)

	for _, mark := range tr.Node.Marks {
		if rendered, err := tr.renderMark(output, mark); err != nil {
			return "", err
		} else {
			output = rendered
		}
	}
	return output, nil
}
func renderText(node schema.Node) (string, error) {
	return newTextRenderer(node).Render()
}
