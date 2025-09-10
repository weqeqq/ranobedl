package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Mark struct {
	Type  MarkType       `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}
type Node struct {
	Type    NodeType       `json:"type"`
	Text    string         `json:"text,omitempty"`
	Marks   []Mark         `json:"marks,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []Node         `json:"content,omitempty"`
}

func (n *Node) ImageSrc() (string, error) {
	if n.Type != NodeTypeImage {
		return "", fmt.Errorf("Node is not image")
	}
	if n.Attrs == nil {
		return "", fmt.Errorf("Node has no attributes")
	}
	if src, found := n.Attrs["src"].(string); found {
		return src, nil
	} else {
		return "", fmt.Errorf("Node has no src attribute")
	}
}
func (m *Mark) LinkHref() (string, error) {
	if m.Type != MarkTypeLink {
		return "", fmt.Errorf("Mark is not link")
	}
	if m.Attrs == nil {
		return "", fmt.Errorf("Mark has no attributes")
	}
	if href, found := m.Attrs["href"].(string); found {
		return href, nil
	} else {
		return "", fmt.Errorf("Mark has no href attribute")
	}
}
func FromString(jsonStr string) (Node, error) {
	var node Node

	if err := json.Unmarshal([]byte(jsonStr), &node); err != nil {
		return Node{}, err
	}
	return node, nil
}
func FromStream(stream io.Reader) (Node, error) {
	var node Node

	if err := json.NewDecoder(stream).Decode(&node); err != nil {
		return Node{}, err
	}
	return node, nil
}
func FromFile(filename string) (Node, error) {
	if file, err := os.Open(filename); err != nil {
		return Node{}, err
	} else {
		defer file.Close()
		return FromStream(file)
	}
}
func (node *Node) ToString() (string, error) {
	bytes, err := json.MarshalIndent(node, "", "  ")
	return string(bytes), err
}
func (node *Node) ToStream(stream io.Writer) error {
	encoder := json.NewEncoder(stream)
	encoder.SetIndent("", "  ")

	return encoder.Encode(node)
}
func (node *Node) ToFile(filename string) error {
	if file, err := os.Create(filename); err != nil {
		return err
	} else {
		defer file.Close()
		return node.ToStream(file)
	}
}
