package nodehandler

import (
	"fmt"
	"ranobedl/format/internal/builder"
	"ranobedl/schema"
)

type RenderInline = func(node []schema.Node) (string, error)

func pushDoc(builder builder.Builder, renderInline RenderInline, node schema.Node) error {
	for _, child := range node.Content {
		if err := PushBlock(builder, renderInline, child); err != nil {
			return err
		}
	}
	return nil
}
func pushParagraph(builder builder.Builder, renderInline RenderInline, node schema.Node) error {
	if rendered, err := renderInline(node.Content); err != nil {
		return err
	} else {
		return builder.PushParagraph(rendered)
	}
}
func pushHeading(builder builder.Builder, renderInline RenderInline, node schema.Node) error {
	return pushParagraph(
		builder,
		renderInline,
		node,
	)
}
func pushBulletList(builder builder.Builder, renderInline RenderInline, node schema.Node) error {
	for _, child := range node.Content {
		if err := pushBulletListItem(builder, renderInline, child); err != nil {
			return err
		}
	}
	return nil
}
func pushOrderedList(builder builder.Builder, renderInline RenderInline, node schema.Node) error {
	for index, child := range node.Content {
		if err := pushOrderedListItem(builder, renderInline, child, index+1); err != nil {
			return err
		}
	}
	return nil
}
func pushBulletListItem(builder builder.Builder, renderInline RenderInline, node schema.Node) error {
	if rendered, err := renderInline(node.Content[0].Content); err != nil {
		return err
	} else {
		return builder.PushParagraph(fmt.Sprintf("	* %s", rendered))
	}
}
func pushOrderedListItem(builder builder.Builder, renderInline RenderInline, node schema.Node, index int) error {
	if rendered, err := renderInline(node.Content[0].Content); err != nil {
		return err
	} else {
		return builder.PushParagraph(fmt.Sprintf("	%d. %s", index, rendered))
	}
}
func pushBlockquote(builder builder.Builder, renderInline RenderInline, node schema.Node) error {
	for _, child := range node.Content {
		if err := PushBlock(builder, renderInline, child); err != nil {
			return err
		}
	}
	return nil
}
func pushCodeBlock(builder builder.Builder, renderInline RenderInline, node schema.Node) error {
	return pushParagraph(
		builder,
		renderInline,
		node,
	)
}
func pushHorizontalRule(builder builder.Builder, _ RenderInline, _ schema.Node) error {
	return builder.PushParagraph("***")
}
func pushImage(builder builder.Builder, _ RenderInline, node schema.Node) error {
	if src, err := node.ImageSrc(); err != nil {
		return err
	} else {
		return builder.PushImage(src)
	}
}

func PushBlock(builder builder.Builder, renderInline RenderInline, node schema.Node) error {
	switch node.Type {
	case schema.NodeTypeDoc:
		return pushDoc(builder, renderInline, node)
	case schema.NodeTypeParagraph:
		return pushParagraph(builder, renderInline, node)
	case schema.NodeTypeHeading:
		return pushHeading(builder, renderInline, node)
	case schema.NodeTypeBulletList:
		return pushBulletList(builder, renderInline, node)
	case schema.NodeTypeOrderedList:
		return pushOrderedList(builder, renderInline, node)
	case schema.NodeTypeBlockquote:
		return pushBlockquote(builder, renderInline, node)
	case schema.NodeTypeCodeBlock:
		return pushCodeBlock(builder, renderInline, node)
	case schema.NodeTypeHorizontalRule:
		return pushHorizontalRule(builder, renderInline, node)
	case schema.NodeTypeImage:
		return pushImage(builder, renderInline, node)
	default:
		panic("Unreachable")
	}
}
