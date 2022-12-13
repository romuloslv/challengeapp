package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romuloslv/challengeapp/internal/database"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Service struct {
	queries *database.Queries
}

func NewService(queries *database.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/", s.Home)
	router.GET("/health", s.CheckHealth)
	router.GET("/version", s.Version)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	router.GET("/accounts", s.List)
	router.GET("/accounts/:person_id", s.Get)
	router.POST("/accounts", s.Create)
	router.PUT("/accounts/:person_id", s.FullUpdate)
	router.PATCH("/accounts/:person_id", s.PartialUpdate)
	router.DELETE("/accounts/:person_id", s.Delete)

}

type apiAccount struct {
	PersonID   string `json:"person_id,omitempty" binding:"omitempty,max=11"`
	FirstName  string `json:"first_name,omitempty" binding:"required,max=30"`
	LastName   string `json:"last_name,omitempty" binding:"required,max=20"`
	WebAddress string `json:"web_address,omitempty" binding:"required"`
	DateBirth  string `json:"date_birth,omitempty" binding:"required"`
}

type apiAccountPartialUpdate struct {
	PersonID   *string `json:"person_id,omitempty" binding:"omitempty,max=11"`
	FirstName  *string `json:"first_name,omitempty" binding:"omitempty,max=30"`
	LastName   *string `json:"last_name,omitempty" binding:"omitempty,max=20"`
	WebAddress *string `json:"web_address,omitempty" binding:"omitempty"`
	DateBirth  *string `json:"date_birth,omitempty" binding:"omitempty"`
}

func fromDB(account database.Account) *apiAccount {
	return &apiAccount{
		PersonID:   account.PersonID,
		FirstName:  account.FirstName,
		LastName:   account.LastName,
		WebAddress: account.WebAddress,
		DateBirth:  account.DateBirth,
	}
}

type pathParameters struct {
	PersonID string `uri:"person_id" binding:"required"`
}

// @Summary Home
// @Description Home
// @Tags home
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router / [get]
func (s *Service) Home(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the Challenge App!")
}

// @Summary Check health
// @Description Check health
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /health [get]
func (s *Service) CheckHealth(c *gin.Context) {
	c.String(http.StatusOK, "{ \"status\" : \"UP\" }")
}

// @Summary Version
// @Description Version
// @Tags version
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /version [get]
func (s *Service) Version(c *gin.Context) {
	c.String(http.StatusOK, "{ \"version\" : \"0.1.0\" }")
}

// @Summary Create account
// @Description Create account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body apiAccount true "Account"
// @Success 201 {object} apiAccount
// @Failure 400 {object} string
// @Failure 503 {object} string
// @Router /accounts [post]
func (s *Service) Create(c *gin.Context) {
	// Parse request
	var request apiAccount
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	params := database.CreateAccountParams{
		PersonID:   request.PersonID,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		WebAddress: request.WebAddress,
		DateBirth:  request.DateBirth,
	}

	account, err := s.queries.CreateAccount(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"Error": err.Error()})
		return
	}

	// Build response
	response := fromDB(account)
	c.IndentedJSON(http.StatusCreated, response)
}

// @Summary Get account
// @Description Get account
// @Tags accounts
// @Accept json
// @Produce json
// @Param person_id path int true "Account PersonID"
// @Success 200 {object} apiAccount
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 503 {object} string
// @Router /accounts/{person_id} [get]
func (s *Service) Get(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Get account
	account, err := s.queries.GetAccount(context.Background(), pathParams.PersonID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"Error": err.Error()})
		return
	}

	// Build response
	response := fromDB(account)
	c.IndentedJSON(http.StatusOK, response)
}

// @Summary Full update account
// @Description Full update account
// @Tags accounts
// @Accept json
// @Produce json
// @Param person_id path int true "Account PersonID"
// @Param account body apiAccount true "Account"
// @Success 200 {object} apiAccount
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 503 {object} string
// @Router /accounts/{person_id} [put]
func (s *Service) FullUpdate(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	var request apiAccount
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Update account
	params := database.UpdateAccountParams{
		PersonID:   pathParams.PersonID,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		WebAddress: request.WebAddress,
		DateBirth:  request.DateBirth,
	}

	fmt.Println(params)
	account, err := s.queries.UpdateAccount(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"Error": err.Error()})
		return
	}

	// Build response
	response := fromDB(account)
	c.IndentedJSON(http.StatusOK, response)
}

// @Summary Partial update account
// @Description Partial update account
// @Tags accounts
// @Accept json
// @Produce json
// @Param person_id path int true "Account PersonID"
// @Param account body apiAccountPartialUpdate true "Account"
// @Success 200 {object} apiAccount
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 503 {object} string
// @Router /accounts/{person_id} [patch]
func (s *Service) PartialUpdate(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	var request apiAccountPartialUpdate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Update account
	params := database.PartialUpdateAccountParams{PersonID: pathParams.PersonID}

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
		params.WebAddress = *request.WebAddress
	}
	if request.DateBirth != nil {
		params.UpdateDateBirth = true
		params.DateBirth = *request.DateBirth
	}

	account, err := s.queries.PartialUpdateAccount(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"Error": err.Error()})
		return
	}

	// Build response
	response := fromDB(account)
	c.IndentedJSON(http.StatusOK, response)
}

// @Summary Delete account
// @Description Delete account
// @Tags accounts
// @Accept json
// @Produce json
// @Param person_id path int true "Account PersonID"
// @Success 204
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 503 {object} string
// @Router /accounts/{person_id} [delete]
func (s *Service) Delete(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Delete account
	if err := s.queries.DeleteAccount(context.Background(), pathParams.PersonID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"Error": err.Error()})
		return
	}

	// Build response
	c.Status(http.StatusOK)
}

// @Summary List accounts
// @Description List accounts
// @Tags accounts
// @Accept json
// @Produce json
// @Param person_id query int false "Person PersonID"
// @Param first_name query string false "First name"
// @Param last_name query string false "Last name"
// @Param web_address query string false "Web address"
// @Param date_birth query string false "Date birth"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} apiAccount
// @Failure 400 {object} string
// @Failure 503 {object} string
// @Router /accounts [get]
func (s *Service) List(c *gin.Context) {
	// List accounts
	accounts, err := s.queries.ListAccounts(context.Background())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"Error": err.Error()})
		return
	}

	if len(accounts) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"Message": "No accounts found!"})
		return
	}

	// Build response
	var response []*apiAccount
	for _, account := range accounts {
		response = append(response, fromDB(account))
	}
	c.IndentedJSON(http.StatusOK, accounts)
}
