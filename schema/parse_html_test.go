package schema

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func nodesEqual(a, b Node) bool {
	if a.Type != b.Type || a.Text != b.Text || len(a.Content) != len(b.Content) || len(a.Marks) != len(b.Marks) {
		return false
	}
	if !reflect.DeepEqual(a.Attrs, b.Attrs) {
		return false
	}
	for i := range a.Content {
		if !nodesEqual(a.Content[i], b.Content[i]) {
			return false
		}
	}
	for i := range a.Marks {
		if a.Marks[i].Type != b.Marks[i].Type || !reflect.DeepEqual(a.Marks[i].Attrs, b.Marks[i].Attrs) {
			return false
		}
	}
	return true
}
func nodeSlicesEqual(a, b []Node) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !nodesEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestParseBlockChildren(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Node
		wantErr  bool
	}{
		{
			name:  "Simple paragraph with text",
			input: `<p>Hello, world!</p>`,
			expected: []Node{
				{
					Type: NodeTypeParagraph,
					Content: []Node{
						{Type: NodeTypeText, Text: "Hello, world!", Marks: []Mark{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "Multiple paragraphs",
			input: `<p>First</p><p>Second</p>`,
			expected: []Node{
				{
					Type: NodeTypeParagraph,
					Content: []Node{
						{Type: NodeTypeText, Text: "First", Marks: []Mark{}},
					},
				},
				{
					Type: NodeTypeParagraph,
					Content: []Node{
						{Type: NodeTypeText, Text: "Second", Marks: []Mark{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "Mixed inline and block elements",
			input: `<p>Text <b>Bold</b></p><h1>Heading</h1>`,
			expected: []Node{
				{
					Type: NodeTypeParagraph,
					Content: []Node{
						{Type: NodeTypeText, Text: "Text ", Marks: []Mark{}},
						{Type: NodeTypeText, Text: "Bold", Marks: []Mark{{Type: MarkTypeBold}}},
					},
				},
				{
					Type:  NodeTypeHeading,
					Attrs: map[string]any{"level": 1},
					Content: []Node{
						{Type: NodeTypeText, Text: "Heading", Marks: []Mark{}},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid block tag",
			input:   `<invalid>Text</invalid>`,
			wantErr: true,
		},
		{
			name:  "Empty paragraph",
			input: `<p></p>`,
			expected: []Node{
				{
					Type:    NodeTypeParagraph,
					Content: []Node{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := html.Parse(strings.NewReader(tt.input))
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Failed to parse HTML: %v", err)
				}
				return
			}
			parser := newHtmlParser(node)
			body := parser.findBody(node)
			if body == nil && !tt.wantErr {
				t.Fatal("No body found in HTML")
			}

			result, err := parser.parseBlockChildren(body)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBlockChildren() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !nodeSlicesEqual(result, tt.expected) {
					t.Errorf("parseBlockChildren() got %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestHandleHeading(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
		wantErr  bool
	}{
		{
			name:  "H1 heading",
			input: `<h1>Hello</h1>`,
			expected: Node{
				Type:  NodeTypeHeading,
				Attrs: map[string]any{"level": 1},
				Content: []Node{
					{Type: NodeTypeText, Text: "Hello", Marks: []Mark{}},
				},
			},
			wantErr: false,
		},
		{
			name:  "H2 with bold text",
			input: `<h2>Bold <b>text</b></h2>`,
			expected: Node{
				Type:  NodeTypeHeading,
				Attrs: map[string]any{"level": 2},
				Content: []Node{
					{Type: NodeTypeText, Text: "Bold ", Marks: []Mark{}},
					{Type: NodeTypeText, Text: "text", Marks: []Mark{{Type: MarkTypeBold}}},
				},
			},
			wantErr: false,
		},
		{
			name:  "H6 heading (max level)",
			input: `<h6>Max level</h6>`,
			expected: Node{
				Type:  NodeTypeHeading,
				Attrs: map[string]any{"level": 6},
				Content: []Node{
					{Type: NodeTypeText, Text: "Max level", Marks: []Mark{}},
				},
			},
			wantErr: false,
		},
		// {
		// 	name:    "Invalid heading level",
		// 	input:   `<h7>Invalid</h7>`,
		// 	wantErr: true,
		// },
		{
			name:  "Empty heading",
			input: `<h1></h1>`,
			expected: Node{
				Type:    NodeTypeHeading,
				Attrs:   map[string]any{"level": 1},
				Content: []Node{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := html.Parse(strings.NewReader(tt.input))
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Failed to parse HTML: %v", err)
				}
				return
			}
			parser := newHtmlParser(node)
			body := parser.findBody(node)
			if body == nil && !tt.wantErr {
				t.Fatal("No body found in HTML")
			}

			result, err := parser.handleHeading(body.FirstChild)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleHeading() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !nodesEqual(result, tt.expected) {
					t.Errorf("handleHeading() got %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestHandleBulletList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
		wantErr  bool
	}{
		{
			name:  "Simple bullet list",
			input: `<ul><li>Item 1</li><li>Item 2</li></ul>`,
			expected: Node{
				Type: NodeTypeBulletList,
				Content: []Node{
					{
						Type: NodeTypeListItem,
						Content: []Node{
							{
								Type: NodeTypeParagraph,
								Content: []Node{
									{Type: NodeTypeText, Text: "Item 1", Marks: []Mark{}},
								},
							},
						},
					},
					{
						Type: NodeTypeListItem,
						Content: []Node{
							{
								Type: NodeTypeParagraph,
								Content: []Node{
									{Type: NodeTypeText, Text: "Item 2", Marks: []Mark{}},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "Nested bullet list",
			input: `<ul><li>Item 1<ul><li>Nested</li></ul></li></ul>`,
			expected: Node{
				Type: NodeTypeBulletList,
				Content: []Node{
					{
						Type: NodeTypeListItem,
						Content: []Node{
							{
								Type: NodeTypeParagraph,
								Content: []Node{
									{Type: NodeTypeText, Text: "Item 1", Marks: []Mark{}},
								},
							},
							{
								Type: NodeTypeBulletList,
								Content: []Node{
									{
										Type: NodeTypeListItem,
										Content: []Node{
											{
												Type: NodeTypeParagraph,
												Content: []Node{
													{Type: NodeTypeText, Text: "Nested", Marks: []Mark{}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "Empty list",
			input: `<ul></ul>`,
			expected: Node{
				Type:    NodeTypeBulletList,
				Content: []Node{},
			},
			wantErr: false,
		},
		{
			name:    "Invalid list item",
			input:   `<ul><invalid>Item</invalid></ul>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := html.Parse(strings.NewReader(tt.input))
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Failed to parse HTML: %v", err)
				}
				return
			}
			parser := newHtmlParser(node)
			body := parser.findBody(node)
			if body == nil && !tt.wantErr {
				t.Fatal("No body found in HTML")
			}

			result, err := parser.handleBulletList(body.FirstChild)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleBulletList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !nodesEqual(result, tt.expected) {
					t.Errorf("handleBulletList() \ngot_ %v, \nwant %v", result, tt.expected)
				}
			}
		})
	}
}

func TestParseInline(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Node
		wantErr  bool
	}{
		{
			name:  "Bold and italic text",
			input: `<b>Bold <i>italic</i></b>`,
			expected: []Node{
				{
					Type:  NodeTypeText,
					Text:  "Bold ",
					Marks: []Mark{{Type: MarkTypeBold}},
				},
				{
					Type:  NodeTypeText,
					Text:  "italic",
					Marks: []Mark{{Type: MarkTypeBold}, {Type: MarkTypeItalic}},
				},
			},
			wantErr: false,
		},
		{
			name:  "Link with text",
			input: `<a href="example.com">Link</a>`,
			expected: []Node{
				{
					Type:  NodeTypeText,
					Text:  "Link",
					Marks: []Mark{{Type: MarkTypeLink, Attrs: map[string]any{"href": "example.com"}}},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid inline tag",
			input:   `<invalid>Text</invalid>`,
			wantErr: true,
		},
		{
			name:  "Multiple formatting",
			input: `<b>Bold <i>italic <u>underline</u></i></b>`,
			expected: []Node{
				{
					Type:  NodeTypeText,
					Text:  "Bold ",
					Marks: []Mark{{Type: MarkTypeBold}},
				},
				{
					Type:  NodeTypeText,
					Text:  "italic ",
					Marks: []Mark{{Type: MarkTypeBold}, {Type: MarkTypeItalic}},
				},
				{
					Type:  NodeTypeText,
					Text:  "underline",
					Marks: []Mark{{Type: MarkTypeBold}, {Type: MarkTypeItalic}, {Type: MarkTypeUnderline}},
				},
			},
			wantErr: false,
		},
		{
			name:     "Empty inline element",
			input:    `<b></b>`,
			expected: []Node{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := html.Parse(strings.NewReader(tt.input))
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Failed to parse HTML: %v", err)
				}
				return
			}
			parser := newHtmlParser(node)
			body := parser.findBody(node)
			if body == nil && !tt.wantErr {
				t.Fatal("No body found in HTML")
			}

			result, err := parser.parseInline(body.FirstChild, []Mark{})
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !nodeSlicesEqual(result, tt.expected) {
					t.Errorf("parseInline() got %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestFromHtmlString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
		wantErr  bool
	}{
		{
			name:  "Complete document",
			input: `<html><body><p>Hello</p></body></html>`,
			expected: Node{
				Type: NodeTypeDoc,
				Content: []Node{
					{
						Type: NodeTypeParagraph,
						Content: []Node{
							{Type: NodeTypeText, Text: "Hello", Marks: []Mark{}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "Fragment without html/body",
			input: `<p>Fragment</p>`,
			expected: Node{
				Type: NodeTypeDoc,
				Content: []Node{
					{
						Type: NodeTypeParagraph,
						Content: []Node{
							{Type: NodeTypeText, Text: "Fragment", Marks: []Mark{}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "Complex document",
			input: `<html><body><h1>Title</h1><p>Text <b>bold</b></p><ul><li>Item</li></ul></body></html>`,
			expected: Node{
				Type: NodeTypeDoc,
				Content: []Node{
					{
						Type:  NodeTypeHeading,
						Attrs: map[string]any{"level": 1},
						Content: []Node{
							{Type: NodeTypeText, Text: "Title", Marks: []Mark{}},
						},
					},
					{
						Type: NodeTypeParagraph,
						Content: []Node{
							{Type: NodeTypeText, Text: "Text ", Marks: []Mark{}},
							{Type: NodeTypeText, Text: "bold", Marks: []Mark{{Type: MarkTypeBold}}},
						},
					},
					{
						Type: NodeTypeBulletList,
						Content: []Node{
							{
								Type: NodeTypeListItem,
								Content: []Node{
									{
										Type: NodeTypeParagraph,
										Content: []Node{
											{Type: NodeTypeText, Text: "Item", Marks: []Mark{}},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromHtmlString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromHtmlString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !nodesEqual(result, tt.expected) {
					t.Errorf("FromHtmlString() got %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
