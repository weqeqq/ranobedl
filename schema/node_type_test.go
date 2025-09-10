package schema

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNodeTypeFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected NodeType
		err      error
	}{
		{"doc", NodeTypeDoc, nil},
		{"paragraph", NodeTypeParagraph, nil},
		{"heading", NodeTypeHeading, nil},
		{"bulletList", NodeTypeBulletList, nil},
		{"orderedList", NodeTypeOrderedList, nil},
		{"listItem", NodeTypeListItem, nil},
		{"blockquote", NodeTypeBlockquote, nil},
		{"codeBlock", NodeTypeCodeBlock, nil},
		{"horizontalRule", NodeTypeHorizontalRule, nil},
		{"text", NodeTypeText, nil},
		{"hardBreak", NodeTypeHardBreak, nil},
		{"image", NodeTypeImage, nil},
		{"invalid", -1, fmt.Errorf("Undefined NodeType: invalid")},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := NodeTypeFromString(tt.input)

			if result != tt.expected {
				t.Errorf("NodeTypeFromString(%q) = %v; want %v", tt.input, result, tt.expected)
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("NodeTypeFromString(%q) error = %v; want %v", tt.input, err, tt.err)
			}
			if err == nil && tt.err != nil {
				t.Errorf("NodeTypeFromString(%q) expected error %v, got nil", tt.input, tt.err)
			}
		})
	}
}
func TestNodeTypeString(t *testing.T) {
	tests := []struct {
		nodeType NodeType
		expected string
	}{
		{NodeTypeDoc, "doc"},
		{NodeTypeParagraph, "paragraph"},
		{NodeTypeHeading, "heading"},
		{NodeTypeBulletList, "bulletList"},
		{NodeTypeOrderedList, "orderedList"},
		{NodeTypeListItem, "listItem"},
		{NodeTypeBlockquote, "blockquote"},
		{NodeTypeCodeBlock, "codeBlock"},
		{NodeTypeHorizontalRule, "horizontalRule"},
		{NodeTypeText, "text"},
		{NodeTypeHardBreak, "hardBreak"},
		{NodeTypeImage, "image"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.nodeType.String()

			if result != tt.expected {
				t.Errorf("NodeType(%v).String() = %q; want %q", tt.nodeType, result, tt.expected)
			}
		})
	}
	t.Run("undefined", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for undefined NodeType, but none occurred")
			}
		}()
		_ = NodeType(-1).String()
	})
}
func TestNodeTypeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected NodeType
		err      error
	}{
		{`"doc"`, NodeTypeDoc, nil},
		{`"paragraph"`, NodeTypeParagraph, nil},
		{`"heading"`, NodeTypeHeading, nil},
		{`"bulletList"`, NodeTypeBulletList, nil},
		{`"orderedList"`, NodeTypeOrderedList, nil},
		{`"listItem"`, NodeTypeListItem, nil},
		{`"blockquote"`, NodeTypeBlockquote, nil},
		{`"codeBlock"`, NodeTypeCodeBlock, nil},
		{`"horizontalRule"`, NodeTypeHorizontalRule, nil},
		{`"text"`, NodeTypeText, nil},
		{`"hardBreak"`, NodeTypeHardBreak, nil},
		{`"image"`, NodeTypeImage, nil},
		{`"invalid"`, -1, fmt.Errorf("Undefined NodeType: invalid")},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var result NodeType
			err := json.Unmarshal([]byte(tt.input), &result)

			if result != tt.expected {
				t.Errorf("UnmarshalJSON(%q) = %v; want %v", tt.input, result, tt.expected)
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("UnmarshalJSON(%q) error = %v; want %v", tt.input, err, tt.err)
			}
			if err == nil && tt.err != nil {
				t.Errorf("UnmarshalJSON(%q) expected error %v, got nil", tt.input, tt.err)
			}
		})
	}
}
func TestNodeTypeMarshalJSON(t *testing.T) {
	tests := []struct {
		nodeType NodeType
		expected string
		err      error
	}{
		{NodeTypeDoc, `"doc"`, nil},
		{NodeTypeParagraph, `"paragraph"`, nil},
		{NodeTypeHeading, `"heading"`, nil},
		{NodeTypeBulletList, `"bulletList"`, nil},
		{NodeTypeOrderedList, `"orderedList"`, nil},
		{NodeTypeListItem, `"listItem"`, nil},
		{NodeTypeBlockquote, `"blockquote"`, nil},
		{NodeTypeCodeBlock, `"codeBlock"`, nil},
		{NodeTypeHorizontalRule, `"horizontalRule"`, nil},
		{NodeTypeText, `"text"`, nil},
		{NodeTypeHardBreak, `"hardBreak"`, nil},
		{NodeTypeImage, `"image"`, nil},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result, err := tt.nodeType.MarshalJSON()

			if string(result) != tt.expected {
				t.Errorf("MarshalJSON(%v) = %q; want %q", tt.nodeType, result, tt.expected)
			}
			if err != tt.err {
				t.Errorf("MarshalJSON(%v) error = %v; want %v", tt.nodeType, err, tt.err)
			}
		})
	}
	t.Run("undefined", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for undefined NodeType, but none occurred")
			}
		}()
		NodeType(-1).MarshalJSON()
	})
}
func TestNodeTypeIsBlock(t *testing.T) {
	tests := []struct {
		nodeType NodeType
		expected bool
	}{
		{NodeTypeDoc, true},
		{NodeTypeParagraph, true},
		{NodeTypeHeading, true},
		{NodeTypeBulletList, true},
		{NodeTypeOrderedList, true},
		{NodeTypeListItem, true},
		{NodeTypeBlockquote, true},
		{NodeTypeCodeBlock, true},
		{NodeTypeHorizontalRule, true},
		{NodeTypeText, false},
		{NodeTypeHardBreak, false},
		{NodeTypeImage, false},
	}
	for _, tt := range tests {
		t.Run(tt.nodeType.String(), func(t *testing.T) {
			result := tt.nodeType.IsBlock()

			if result != tt.expected {
				t.Errorf("NodeType(%v).IsBlock() = %v; want %v", tt.nodeType, result, tt.expected)
			}
		})
	}
}

func TestNodeTypeIsInline(t *testing.T) {
	tests := []struct {
		nodeType NodeType
		expected bool
	}{
		{NodeTypeDoc, false},
		{NodeTypeParagraph, false},
		{NodeTypeHeading, false},
		{NodeTypeBulletList, false},
		{NodeTypeOrderedList, false},
		{NodeTypeListItem, false},
		{NodeTypeBlockquote, false},
		{NodeTypeCodeBlock, false},
		{NodeTypeHorizontalRule, false},
		{NodeTypeText, true},
		{NodeTypeHardBreak, true},
		{NodeTypeImage, true},
	}
	for _, tt := range tests {
		t.Run(tt.nodeType.String(), func(t *testing.T) {
			result := tt.nodeType.IsInline()

			if result != tt.expected {
				t.Errorf("NodeType(%v).IsInline() = %v; want %v", tt.nodeType, result, tt.expected)
			}
		})
	}
}

func TestNodeTypeGroup(t *testing.T) {
	tests := []struct {
		nodeType NodeType
		expected NodeGroup
	}{
		{NodeTypeDoc, NodeGroupBlock},
		{NodeTypeParagraph, NodeGroupBlock},
		{NodeTypeHeading, NodeGroupBlock},
		{NodeTypeBulletList, NodeGroupBlock},
		{NodeTypeOrderedList, NodeGroupBlock},
		{NodeTypeListItem, NodeGroupBlock},
		{NodeTypeBlockquote, NodeGroupBlock},
		{NodeTypeCodeBlock, NodeGroupBlock},
		{NodeTypeHorizontalRule, NodeGroupBlock},
		{NodeTypeText, NodeGroupInline},
		{NodeTypeHardBreak, NodeGroupInline},
		{NodeTypeImage, NodeGroupInline},
	}
	for _, tt := range tests {
		t.Run(tt.nodeType.String(), func(t *testing.T) {
			result := tt.nodeType.Group()

			if result != tt.expected {
				t.Errorf("NodeType(%v).Group() = %v; want %v", tt.nodeType, result, tt.expected)
			}
		})
	}
	t.Run("undefined", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for undefined NodeType, but none occurred")
			}
		}()
		NodeType(-1).Group()
	})
}

func TestNodeTypeHTMLTagFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fn       func(string) bool
		expected bool
	}{
		{"isParagraph", "p", isParagraph, true},
		{"isParagraph", "div", isParagraph, false},
		{"isHeading", "h1", isHeading, true},
		{"isHeading", "h2", isHeading, true},
		{"isHeading", "h3", isHeading, true},
		{"isHeading", "h4", isHeading, true},
		{"isHeading", "h5", isHeading, true},
		{"isHeading", "h6", isHeading, true},
		{"isHeading", "p", isHeading, false},
		{"isBulletList", "ul", isBulletList, true},
		{"isBulletList", "ol", isBulletList, false},
		{"isOrderedList", "ol", isOrderedList, true},
		{"isOrderedList", "ul", isOrderedList, false},
		{"isListItem", "li", isListItem, true},
		{"isListItem", "div", isListItem, false},
		{"isBlockquote", "blockquote", isBlockquote, true},
		{"isBlockquote", "p", isBlockquote, false},
		{"isCodeBlock", "pre", isCodeBlock, true},
		{"isCodeBlock", "code", isCodeBlock, false},
		{"isHorizontalRule", "hr", isHorizontalRule, true},
		{"isHorizontalRule", "br", isHorizontalRule, false},
		{"isHardBreak", "br", isHardBreak, true},
		{"isHardBreak", "hr", isHardBreak, false},
		{"isImage", "img", isImage, true},
		{"isImage", "a", isImage, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.input)

			if result != tt.expected {
				t.Errorf("%s(%q) = %v; want %v", tt.name, tt.input, result, tt.expected)
			}
		})
	}
}
