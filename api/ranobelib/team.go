package ranobelib

type team struct {
	Id      int    `json:"id"`
	Slug    string `json:"slug"`
	SlugUrl string `json:"slug_url"`
	Model   string `json:"model"`
	Name    string `json:"name"`
	Cover   cover  `json:"cover"`
}
