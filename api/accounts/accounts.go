package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romuloslv/go-rest-api/local/database"
)

type Service struct {
	queries *database.Queries
}

func NewService(queries *database.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.POST("/accounts", s.Create)
	router.GET("/accounts/:id", s.Get)
	router.PUT("/accounts/:id", s.FullUpdate)
	router.PATCH("/accounts/:id", s.PartialUpdate)
	router.DELETE("/accounts/:id", s.Delete)
	router.GET("/accounts", s.List)
}

type apiAccount struct {
	ID         int64
	PersonID   string         `json:"person_id,omitempty" binding:"omitempty,max=11"`
	FirstName  string         `json:"first_name,omitempty" binding:"required,max=30"`
	LastName   string         `json:"last_name,omitempty" binding:"required,max=20"`
	WebAddress sql.NullString `json:"web_address,omitempty" binding:"required,email"`
	DateBirth  sql.NullTime   `json:"date_birth,omitempty" binding:"required,datetime"`
}

type apiAccountPartialUpdate struct {
	PersonID   *string         `json:"person_id,omitempty" binding:"omitempty,max=11"`
	FirstName  *string         `json:"first_name,omitempty" binding:"omitempty,max=30"`
	LastName   *string         `json:"last_name,omitempty" binding:"omitempty,max=20"`
	WebAddress *sql.NullString `json:"web_address,omitempty" binding:"omitempty,email"`
	DateBirth  *sql.NullTime   `json:"date_birth,omitempty" binding:"omitempty,datetime"`
}

func fromDB(account database.Account) *apiAccount {
	return &apiAccount{
		ID:         account.ID,
		PersonID:   account.PersonID,
		FirstName:  account.FirstName,
		LastName:   account.LastName,
		WebAddress: account.WebAddress,
		DateBirth:  account.DateBirth,
	}
}

type pathParameters struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Service) Create(c *gin.Context) {
	// Parse request
	var request apiAccount
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create account
	params := database.CreateAccountParams{
		PersonID:   request.PersonID,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		WebAddress: request.WebAddress,
		DateBirth:  request.DateBirth,
	}

	account, err := s.queries.CreateAccount(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(account)
	c.IndentedJSON(http.StatusCreated, response)
}

func (s *Service) Get(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get account
	account, err := s.queries.GetAccount(context.Background(), pathParams.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(account)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) FullUpdate(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request apiAccount
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update account
	params := database.UpdateAccountParams{
		ID:         pathParams.ID,
		PersonID:   request.PersonID,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		WebAddress: request.WebAddress,
		DateBirth:  request.DateBirth,
	}

	fmt.Println(params)
	account, err := s.queries.UpdateAccount(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(account)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) PartialUpdate(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request apiAccountPartialUpdate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update account
	params := database.PartialUpdateAccountParams{ID: pathParams.ID}

	if request.PersonID != nil {
		params.UpdatePersonID = true
		params.PersonID = *request.PersonID
	}
	if request.FirstName != nil {
		params.UpdateFirstName = true
		params.FirstName = *request.FirstName
	}
	if request.LastName != nil {
		params.UpdateLastName = true
		params.LastName = *request.LastName
	}
	if request.WebAddress != nil {
		params.UpdateWebAddress = true
		params.WebAddress = (*request.WebAddress).String
	}
	if request.DateBirth != nil {
		params.UpdateDateBirth = true
		params.DateBirth = (*request.DateBirth).Time
	}

	account, err := s.queries.PartialUpdateAccount(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(account)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) Delete(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete account
	if err := s.queries.DeleteAccount(context.Background(), pathParams.ID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	c.Status(http.StatusOK)
}

func (s *Service) List(c *gin.Context) {
	// List accounts
	accounts, err := s.queries.ListAccounts(context.Background())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	if len(accounts) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Build response
	var response []*apiAccount
	for _, account := range accounts {
		response = append(response, fromDB(account))
	}
	c.IndentedJSON(http.StatusOK, accounts)
}
