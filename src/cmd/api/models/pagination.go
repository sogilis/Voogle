package models

type PaginationAttribute int

const (
	TITLE PaginationAttribute = iota
	UPLOADEDAT
	CREATEDAT
	UPDATEDAT
)

type Pagination struct {
	Page      uint
	Limit     uint
	Ascending bool
	Attribute PaginationAttribute
}
