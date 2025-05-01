package model

type Person struct {
	ID          int64   `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Surname     string  `db:"surname" json:"surname"`
	Patronymic  *string `db:"patronymic,omitempty" json:"patronymic,omitempty"`
	Gender      string  `db:"gender" json:"gender"`
	Age         int     `db:"age" json:"age"`
	Nationality string  `db:"nationality" json:"nationality"`
}
