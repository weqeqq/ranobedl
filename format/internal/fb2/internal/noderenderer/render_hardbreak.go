package noderenderer

import (
	"fmt"
	"ranobedl/schema"
)

func renderHardBreak(node schema.Node) (string, error) {
	if node.Type != schema.NodeTypeHardBreak {
		return "", fmt.Errorf("Node is not hardbreak")
	}
	return "<empty-line/>", nil
}
