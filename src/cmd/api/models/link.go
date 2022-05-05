package models

type Link struct {
	Href   string
	Method string
}

func CreateLink(href string, method string) *Link {
	return &Link{Href: href, Method: method}
}
