package model

type Sorting struct {
	SortBy   SortBy
	SortType SortType
}

type SortBy int

const (
	SortByID SortBy = iota
	SortByTitle
	SortByDescription
	SortByDueDate
)

type SortType int

const (
	SortAscending SortType = iota
	SortDescending
)

type Pagination struct {
	Page  int
	Limit int
}
