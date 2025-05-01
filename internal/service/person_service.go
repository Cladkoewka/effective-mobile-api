package service

import (
	"context"
	"fmt"

	"github.com/Cladkoewka/effective-mobile-api/internal/api/enrichment"
	"github.com/Cladkoewka/effective-mobile-api/internal/dto"
	"github.com/Cladkoewka/effective-mobile-api/internal/model"
	"github.com/Cladkoewka/effective-mobile-api/internal/repository"
)

type PersonService struct {
	repo     *repository.PersonRepository
	enricher enrichment.PersonEnricher
}

func NewPersonService(repo *repository.PersonRepository, enricher enrichment.PersonEnricher) *PersonService {
	return &PersonService{repo: repo, enricher: enricher}
}

func (s *PersonService) GetByID(id uint64) (*model.Person, error) {
	return s.repo.GetByID(id)
}

func (s *PersonService) GetAll(resp *dto.GetPersonRequest) ([]model.Person, error) {
	filter := &model.PersonFilter{
		Name:        resp.Name,
		Surname:     resp.Surname,
		Patronymic:  resp.Patronymic,
		Gender:      resp.Gender,
		Nationality: resp.Nationality,
		AgeMin:      resp.AgeMin,
		AgeMax:      resp.AgeMax,
		Limit:       resp.Limit,
		Offset:      resp.Offset,
		SortBy:      resp.SortBy,
		Order:       resp.Order,
	}

	return s.repo.GetAll(filter)
}

func (s *PersonService) Create(ctx context.Context, req *dto.CreatePersonRequest) (uint64, error) {
	if req.Name == "" || req.Surname == "" {
		return 0, fmt.Errorf("name and surname are required")
	}

	enriched, err := s.enricher.Enrich(ctx, req.Name)
	if err != nil {
		return 0, fmt.Errorf("failed to enrich person: %w", err)
	}

	person := &model.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Gender:      enriched.Gender,
		Age:         enriched.Age,
		Nationality: enriched.Nationality,
	}

	return s.repo.Create(person)
}

func (s *PersonService) Update(id uint64, req *dto.UpdatePersonRequest) error {
	person, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("not found: %w", err)
	}

	if req.Name != nil {
		person.Name = *req.Name
	}
	if req.Surname != nil {
		person.Surname = *req.Surname
	}
	if req.Patronymic != nil {
		person.Patronymic = req.Patronymic
	}

	return s.repo.Update(person)
}

func (s *PersonService) Delete(id uint64) error {
	return s.repo.Delete(id)
}
