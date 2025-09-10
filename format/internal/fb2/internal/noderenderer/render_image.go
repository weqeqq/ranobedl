package noderenderer

import (
	"fmt"
	"ranobedl/schema"
)

func renderImage(node schema.Node) (string, error) {
	if src, err := node.ImageSrc(); err != nil {
		return "", err
	} else {
		return fmt.Sprintf(`<image l:href="%s"/>`, src), nil
	}
}
