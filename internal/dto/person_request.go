package dto

type CreatePersonRequest struct {
	Name       string  `json:"name" binding:"required"`
	Surname    string  `json:"surname" binding:"required"`
	Patronymic *string `json:"patronymic,omitempty"`
}

type UpdatePersonRequest struct {
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
}

type GetPersonRequest struct {
	Name        *string `form:"name"`
	Surname     *string `form:"surname"`
	Patronymic  *string `form:"patronymic"`
	Gender      *string `form:"gender"`
	Nationality *string `form:"nationality"`
	AgeMin      *int    `form:"age_min"`
	AgeMax      *int    `form:"age_max"`
	Limit       int     `form:"limit,default=10"`
	Offset      int     `form:"offset,default=0"`
	SortBy      string  `form:"sort_by,default=id"`
	Order       string  `form:"order,default=asc"`
}
