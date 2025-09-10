package ranobelib

import (
	"net/url"
	"strings"
)

const apiUrl = "https://api.cdnlibs.org/api"

func GetUniqueName(urlStr string) (string, error) {
	const nameIndex = 3

	if url, err := url.Parse(urlStr); err != nil {
		return "", err
	} else {
		return strings.Split(url.Path, "/")[nameIndex], nil
	}
}
