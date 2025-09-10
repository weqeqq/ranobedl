package fb2

import (
	"ranobedl/format/internal/fb2/internal/noderenderer"
	"ranobedl/schema"
)

func RenderInline(node []schema.Node) (string, error) {
	return noderenderer.RenderInlineChildren(node)
}
