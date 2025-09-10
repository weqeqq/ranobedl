package schema

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarkTypeFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected MarkType
		err      error
	}{
		{"bold", MarkTypeBold, nil},
		{"italic", MarkTypeItalic, nil},
		{"underline", MarkTypeUnderline, nil},
		{"strike", MarkTypeStrike, nil},
		{"code", MarkTypeCode, nil},
		{"link", MarkTypeLink, nil},
		{"invalid", -1, fmt.Errorf("Undefined MarkType: invalid")},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := MarkTypeFromString(tt.input)

			if result != tt.expected {
				t.Errorf("MarkTypeFromString(%q) = %v; want %v", tt.input, result, tt.expected)
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("MarkTypeFromString(%q) error = %v; want %v", tt.input, err, tt.err)
			}
			if err == nil && tt.err != nil {
				t.Errorf("MarkTypeFromString(%q) error = nil; want %v", tt.input, tt.err)
			}
		})
	}
}
func TestMarkTypeString(t *testing.T) {
	tests := []struct {
		input    MarkType
		expected string
	}{
		{MarkTypeBold, "bold"},
		{MarkTypeItalic, "italic"},
		{MarkTypeUnderline, "underline"},
		{MarkTypeStrike, "strike"},
		{MarkTypeCode, "code"},
		{MarkTypeLink, "link"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("MarkType(%v).String() = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}

	t.Run("undefined", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MarkType(%v).String() did not panic", MarkType(-1))
			}
		}()
		_ = MarkType(-1).String()
	})
}
func TestMarkTypeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected MarkType
		err      error
	}{
		{`"bold"`, MarkTypeBold, nil},
		{`"italic"`, MarkTypeItalic, nil},
		{`"underline"`, MarkTypeUnderline, nil},
		{`"strike"`, MarkTypeStrike, nil},
		{`"code"`, MarkTypeCode, nil},
		{`"link"`, MarkTypeLink, nil},
		{`"invalid"`, -1, fmt.Errorf("Undefined MarkType: invalid")},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var result MarkType
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
func TestMarkTypeMarshalJSON(t *testing.T) {
	tests := []struct {
		markType MarkType
		expected string
		err      error
	}{
		{MarkTypeBold, `"bold"`, nil},
		{MarkTypeItalic, `"italic"`, nil},
		{MarkTypeUnderline, `"underline"`, nil},
		{MarkTypeStrike, `"strike"`, nil},
		{MarkTypeCode, `"code"`, nil},
		{MarkTypeLink, `"link"`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result, err := tt.markType.MarshalJSON()

			if string(result) != tt.expected {
				t.Errorf("MarshalJSON(%v) = %q; want %q", tt.markType, result, tt.expected)
			}
			if err != tt.err {
				t.Errorf("MarshalJSON(%v) error = %v; want %v", tt.markType, err, tt.err)
			}
		})
	}
	t.Run("undefined", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for undefined MarkType, but none occurred")
			}
		}()
		MarkType(-1).MarshalJSON()
	})
}

func TestMarkTypeHTMLTagFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fn       func(string) bool
		expected bool
	}{
		{"isBold", "b", isBold, true},
		{"isBold", "i", isBold, false},
		{"isItalic", "i", isItalic, true},
		{"isItalic", "b", isItalic, false},
		{"isUnderline", "u", isUnderline, true},
		{"isUnderline", "s", isUnderline, false},
		{"isStrike", "s", isStrike, true},
		{"isStrike", "u", isStrike, false},
		{"isCode", "code", isCode, true},
		{"isCode", "a", isCode, false},
		{"isLink", "a", isLink, true},
		{"isLink", "code", isLink, false},
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
