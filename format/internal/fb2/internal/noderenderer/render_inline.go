package noderenderer

import (
	"ranobedl/schema"
)

func renderInline(node schema.Node) (string, error) {
	switch node.Type {
	case schema.NodeTypeText:
		return renderText(node)
	case schema.NodeTypeHardBreak:
		return renderHardBreak(node)
	case schema.NodeTypeImage:
		return renderImage(node)
	default:
		panic("Unreachable code")
	}
}

func RenderInlineChildren(node []schema.Node) (string, error) {
	output := ""

	for _, child := range node {
		if rendered, err := renderInline(child); err != nil {
			return "", err
		} else {
			output += rendered
		}
	}
	return output, nil
}
