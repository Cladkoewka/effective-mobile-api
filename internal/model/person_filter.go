package model

type PersonFilter struct {
	Name        *string
	Surname     *string
	Patronymic  *string
	Gender      *string
	Nationality *string
	AgeMin      *int
	AgeMax      *int
	Limit       int
	Offset      int
	SortBy      string
	Order       string
}
