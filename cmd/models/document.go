package models

type Document struct {
	Name string
}

func (document Document) String() string {
	return document.Name
}
