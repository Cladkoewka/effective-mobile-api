package repository

import (
	"fmt"
	"log"

	"github.com/Cladkoewka/effective-mobile-api/internal/model"
	"github.com/jmoiron/sqlx"
)

type PersonRepository struct {
	db *sqlx.DB
}

func NewPersonRepository(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) GetAll(f *model.PersonFilter) ([]model.Person, error) {
	var persons []model.Person
	query := `SELECT * FROM persons WHERE 1=1`

	// filtration
	if f.Name != nil {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", *f.Name)
	}
	if f.Surname != nil {
		query += fmt.Sprintf(" AND surname ILIKE '%%%s%%'", *f.Surname)
	}
	if f.Patronymic != nil {
		query += fmt.Sprintf(" AND patronymic ILIKE '%%%s%%'", *f.Patronymic)
	}
	if f.Gender != nil {
		query += fmt.Sprintf(" AND gender ILIKE '%%%s%%'", *f.Gender)
	}
	if f.Nationality != nil {
		query += fmt.Sprintf(" AND nationality ILIKE '%%%s%%'", *f.Nationality)
	}
	if f.AgeMin != nil {
		query += fmt.Sprintf(" AND age >= %d", *f.AgeMin)
	}
	if f.AgeMax != nil {
		query += fmt.Sprintf(" AND age <= %d", *f.AgeMax)
	}

	// sorting
	query += " ORDER BY " + f.SortBy + " " + f.Order

	// pagination
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", f.Limit, f.Offset)

	log.Println("Query: ", query)

	err := r.db.Select(&persons, query)
	return persons, err
}


func (r *PersonRepository) GetByID(id uint64) (*model.Person, error) {
	var person model.Person
	query := `SELECT * FROM persons WHERE id = $1`
	err := r.db.Get(&person, query, id)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *PersonRepository) Create(p *model.Person) (uint64, error) {
	var id uint64
	query := `INSERT INTO persons (name, surname, patronymic, gender, age, nationality)
						VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRowx(query, p.Name, p.Surname, p.Patronymic, p.Gender, p.Age, p.Nationality).Scan(&id)
	return id, err
}

func (r *PersonRepository) Update(p *model.Person) error {
	query := `UPDATE persons 
			  SET name = $1, surname = $2, patronymic = $3, gender = $4, age = $5, nationality = $6
			  WHERE id = $7`
	_, err := r.db.Exec(query, p.Name, p.Surname, p.Patronymic, p.Gender, p.Age, p.Nationality, p.ID)
	return err
}

func (r *PersonRepository) Delete(id uint64) error {
	query := `DELETE FROM persons WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
