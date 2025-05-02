package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Cladkoewka/effective-mobile-api/internal/dto"
	"github.com/Cladkoewka/effective-mobile-api/internal/logger"
	"github.com/Cladkoewka/effective-mobile-api/internal/service"
)

type PersonHandler struct {
	personService *service.PersonService
}

func NewPersonHandler(personService *service.PersonService) *PersonHandler {
	return &PersonHandler{personService: personService}
}

func (h *PersonHandler) RegisterRoutes(rg *gin.RouterGroup) {
	person := rg.Group("/persons")
	person.POST("/", h.createPerson)
	person.GET("/", h.getAllPersons)
	person.GET("/:id", h.getPersonByID)
	person.PUT("/:id", h.updatePerson)
	person.DELETE("/:id", h.deletePerson)
}

// @Summary Get a person by ID
// @Tags Person
// @Description Retrieve a person by their ID
// @Accept json
// @Produce json
// @Param id path uint64 true "Person ID"
// @Success 200 {object} model.Person "Successfully retrieved person"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Person not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /persons/{id} [get]
func (h *PersonHandler) getPersonByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error("Invalid ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	person, err := h.personService.GetByID(id)
	if err != nil {
		logger.Log.Error("Person not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("Successfully retrieved person", "id", id)
	c.JSON(http.StatusOK, person)
}

// @Summary Get all persons with filters
// @Tags Person
// @Description Retrieve all persons with optional filters, sorting, and pagination
// @Accept json
// @Produce json
// @Param request query dto.GetPersonRequest true "Filter and pagination parameters"
// @Success 200 {array} model.Person "List of persons"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /persons [get]
func (h *PersonHandler) getAllPersons(c *gin.Context) {
	var req dto.GetPersonRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Log.Error("Failed to bind query parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	persons, err := h.personService.GetAll(&req)
	if err != nil {
		logger.Log.Error("Error retrieving persons", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("Successfully retrieved all persons")
	c.JSON(http.StatusOK, persons)
}

// @Summary Create a new person
// @Tags Person
// @Description Create a new person with enriched data
// @Accept json
// @Produce json
// @Param request body dto.CreatePersonRequest true "Create person request"
// @Success 201 {object} map[string]interface{} "Successfully created person"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /persons [post]
func (h *PersonHandler) createPerson(c *gin.Context) {
	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Failed to bind request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.personService.Create(c.Request.Context(), &req)
	if err != nil {
		logger.Log.Error("Error creating person", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("Successfully created person", "id", id)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary Update a person by ID
// @Tags Person
// @Description Update an existing person by their ID
// @Accept json
// @Produce json
// @Param id path uint64 true "Person ID"
// @Param request body dto.UpdatePersonRequest true "Update person request"
// @Success 200 {object} map[string]interface{} "Successfully updated person"
// @Failure 400 {object} map[string]interface{} "Invalid ID or input"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /persons/{id} [put]
func (h *PersonHandler) updatePerson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error("Invalid ID parameter for update", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Failed to bind request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.personService.Update(id, &req); err != nil {
		logger.Log.Error("Error updating person", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("Successfully updated person", "id", id)
	c.Status(http.StatusOK)
}

// @Summary Delete a person by ID
// @Tags Person
// @Description Delete a person by their ID
// @Accept json
// @Produce json
// @Param id path uint64 true "Person ID"
// @Success 204 {object} map[string]interface{} "Successfully deleted person"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /persons/{id} [delete]
func (h *PersonHandler) deletePerson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error("Invalid ID parameter for deletion", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.personService.Delete(id); err != nil {
		logger.Log.Error("Error deleting person", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("Successfully deleted person", "id", id)
	c.Status(http.StatusNoContent)
}
