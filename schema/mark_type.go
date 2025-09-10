package schema

import (
	"encoding/json"
	"fmt"
)

type MarkType int

const (
	MarkTypeBold MarkType = iota
	MarkTypeItalic
	MarkTypeUnderline
	MarkTypeStrike
	MarkTypeCode
	MarkTypeLink
)

func MarkTypeFromString(str string) (MarkType, error) {
	switch str {
	case "bold":
		return MarkTypeBold, nil
	case "italic":
		return MarkTypeItalic, nil
	case "underline":
		return MarkTypeUnderline, nil
	case "strike":
		return MarkTypeStrike, nil
	case "code":
		return MarkTypeCode, nil
	case "link":
		return MarkTypeLink, nil
	default:
		return -1, fmt.Errorf("Undefined MarkType: %s", str)
	}
}
func (mt *MarkType) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)

	if markType, err := MarkTypeFromString(str); err != nil {
		*mt = -1
		return err
	} else {
		*mt = markType
		return nil
	}
}
func (mt MarkType) String() string {
	switch mt {
	case MarkTypeBold:
		return "bold"
	case MarkTypeItalic:
		return "italic"
	case MarkTypeUnderline:
		return "underline"
	case MarkTypeStrike:
		return "strike"
	case MarkTypeCode:
		return "code"
	case MarkTypeLink:
		return "link"
	default:
		panic(fmt.Sprintf("Undefined MarkType: %d", mt))
	}
}
func (mt MarkType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}
func isBold(html string) bool {
	return html == "b"
}
func isItalic(html string) bool {
	return html == "i"
}
func isUnderline(html string) bool {
	return html == "u"
}
func isStrike(html string) bool {
	return html == "s"
}
func isCode(html string) bool {
	return html == "code"
}
func isLink(html string) bool {
	return html == "a"
}
