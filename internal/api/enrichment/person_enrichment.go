package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Cladkoewka/effective-mobile-api/internal/model"
	"net/http"
	"time"
)

const (
	clientTimeoutInSeconds = 5
	agifyURL               = "https://api.agify.io/?name=%s"
	genderizeURL           = "https://api.genderize.io/?name=%s"
	nationalizeURL         = "https://api.nationalize.io/?name=%s"
)

type PersonEnricher interface {
	Enrich(ctx context.Context, name string) (*model.EnrichmentResult, error)
}

type HTTPEnrichmentClient struct {
	client *http.Client
}

func NewHTTPEnrichmentClient() *HTTPEnrichmentClient {
	return &HTTPEnrichmentClient{
		client: &http.Client{Timeout: clientTimeoutInSeconds * time.Second},
	}
}

func (c *HTTPEnrichmentClient) Enrich(ctx context.Context, name string) (*model.EnrichmentResult, error) {
	result := &model.EnrichmentResult{}

	if err := c.enrichAge(ctx, name, result); err != nil {
		return nil, fmt.Errorf("age enrichment failed: %w", err)
	}
	if err := c.enrichGender(ctx, name, result); err != nil {
		return nil, fmt.Errorf("gender enrichment failed: %w", err)
	}
	if err := c.enrichNationality(ctx, name, result); err != nil {
		return nil, fmt.Errorf("nationality enrichment failed: %w", err)
	}

	return result, nil
}

func (c *HTTPEnrichmentClient) enrichAge(ctx context.Context, name string, result *model.EnrichmentResult) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(agifyURL, name), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data struct {
		Age int `json:"age"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	result.Age = data.Age
	return nil
}

func (c *HTTPEnrichmentClient) enrichGender(ctx context.Context, name string, result *model.EnrichmentResult) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(genderizeURL, name), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data struct {
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	result.Gender = data.Gender
	return nil
}

func (c *HTTPEnrichmentClient) enrichNationality(ctx context.Context, name string, result *model.EnrichmentResult) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(nationalizeURL, name), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if len(data.Country) > 0 {
		result.Nationality = data.Country[0].CountryID
	}
	return nil
}
