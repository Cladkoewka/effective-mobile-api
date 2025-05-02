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

func (s *PersonService) GetAll(req *dto.GetPersonRequest) ([]model.Person, error) {
	if err := validateGetPersonRequest(req); err != nil {
		return nil, err
	}

	filter := mapGetPersonRequestToFilter(req)
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

func validateGetPersonRequest(req *dto.GetPersonRequest) error {
	allowedSortBy := map[string]bool{
		"id": true, "name": true, "surname": true, "patronymic": true,
		"gender": true, "nationality": true, "age": true,
	}
	if !allowedSortBy[req.SortBy] {
		return fmt.Errorf("invalid sort_by field: %s", req.SortBy)
	}

	if req.Order != "asc" && req.Order != "desc" {
		return fmt.Errorf("invalid order value: %s (must be 'asc' or 'desc')", req.Order)
	}

	if req.Limit < 1 || req.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}
	if req.Offset < 0 {
		return fmt.Errorf("offset must be 0 or greater")
	}

	if req.AgeMin != nil && *req.AgeMin < 0 {
		return fmt.Errorf("age_min must be non-negative")
	}
	if req.AgeMax != nil && *req.AgeMax < 0 {
		return fmt.Errorf("age_max must be non-negative")
	}
	if req.AgeMin != nil && req.AgeMax != nil && *req.AgeMin > *req.AgeMax {
		return fmt.Errorf("age_min cannot be greater than age_max")
	}

	return nil
}

func mapGetPersonRequestToFilter(req *dto.GetPersonRequest) *model.PersonFilter {
	return &model.PersonFilter{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Gender:      req.Gender,
		Nationality: req.Nationality,
		AgeMin:      req.AgeMin,
		AgeMax:      req.AgeMax,
		Limit:       req.Limit,
		Offset:      req.Offset,
		SortBy:      req.SortBy,
		Order:       req.Order,
	}
}
