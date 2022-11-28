package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/romuloslv/go-rest-api/cmd/app/config"
	"github.com/romuloslv/go-rest-api/internal/database"
	"github.com/stretchr/testify/suite"
)

type apiError struct {
	Error string
}

type ServiceTestSuite struct {
	suite.Suite
	router  *gin.Engine
	queries *database.Queries
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupSuite() {
	cfg, err := config.Read()
	suite.Require().NoError(err)

	postgres, err := database.NewPostgres(cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password)
	suite.Require().NoError(err)

	suite.queries = database.New(postgres.DB)
	service := NewService(suite.queries)

	suite.router = gin.Default()
	service.RegisterHandlers(suite.router)
}

func (suite *ServiceTestSuite) TestCreate() {
	request := apiAccount{
		PersonID:   "11111111111",
		FirstName:  "John",
		LastName:   "Doe",
		WebAddress: "john.doe@test.local",
		DateBirth:  "1990-01-01",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(request))

	req, err := http.NewRequest("POST", "/accounts", &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusCreated, rec.Result().StatusCode)
	var account apiAccount
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&account))
	suite.Require().Equal(request.PersonID, account.PersonID)
	suite.Require().Equal(request.FirstName, account.FirstName)
	suite.Require().Equal(request.LastName, account.LastName)
	suite.Require().Equal(request.WebAddress, account.WebAddress)
	suite.Require().Equal(request.DateBirth, account.DateBirth)
}

func (suite *ServiceTestSuite) TestCreateBadRequest() {
	request := apiAccount{
		PersonID:   "111111111111",
		FirstName:  "John",
		LastName:   "Doe",
		WebAddress: "john.doe@test.local",
		DateBirth:  "1990-01-01",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(request))

	req, err := http.NewRequest("POST", "/accounts", &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)
	var apiErr apiError
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&apiErr))
	suite.Require().Contains(apiErr.Error, "max")
}

func (suite *ServiceTestSuite) TestGet() {
	account, err := suite.queries.CreateAccount(context.Background(), database.CreateAccountParams{
		PersonID:   "55555555555",
		FirstName:  "John",
		LastName:   "Doe",
		WebAddress: "john.doe@test.local",
		DateBirth:  "1990-01-01",
	})
	suite.Require().NoError(err)

	req, err := http.NewRequest("GET", fmt.Sprintf("/accounts/%d", account.ID), nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	var got apiAccount
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&got))
	suite.Require().Equal(account.ID, got.ID)
	suite.Require().Equal(account.PersonID, got.PersonID)
	suite.Require().Equal(account.FirstName, got.FirstName)
	suite.Require().Equal(account.LastName, got.LastName)
	suite.Require().Equal(account.WebAddress, got.WebAddress)
	suite.Require().Equal(account.DateBirth, got.DateBirth)
}

func (suite *ServiceTestSuite) TestGetNotFound() {
	req, err := http.NewRequest("GET", "/accounts/123", nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusNotFound, rec.Result().StatusCode)
}

func (suite *ServiceTestSuite) TestGetBadRequest() {
	req, err := http.NewRequest("GET", "/accounts/bad-request", nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)
}

func (suite *ServiceTestSuite) TestFullUpdateBadRequest() {
	account, err := suite.queries.CreateAccount(context.Background(), database.CreateAccountParams{
		PersonID:   "44444444444",
		FirstName:  "John",
		LastName:   "Doe",
		WebAddress: "john.doe@test.local",
		DateBirth:  "1990-01-01",
	})
	suite.Require().NoError(err)

	request := apiAccount{
		PersonID:   "444444444444",
		FirstName:  "John",
		LastName:   "Doe",
		WebAddress: "john.doe@test.local",
		DateBirth:  "1990-01-01",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(request))

	req, err := http.NewRequest("PUT", fmt.Sprintf("/accounts/%d", account.ID), &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)
}

func (suite *ServiceTestSuite) TestFullUpdate() {
	account, err := suite.queries.CreateAccount(context.Background(), database.CreateAccountParams{
		PersonID:   "33333333333",
		FirstName:  "John",
		LastName:   "Doe",
		WebAddress: "john.doe@test.local",
		DateBirth:  "1990-01-01",
	})
	suite.Require().NoError(err)

	update := apiAccount{
		PersonID:   "33333333333",
		FirstName:  "Jane",
		LastName:   "Carson",
		WebAddress: "jane.carson@test.local",
		DateBirth:  "2022-01-01",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(update))

	req, err := http.NewRequest("PUT", fmt.Sprintf("/accounts/%d", account.ID), &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	var got apiAccount
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&got))
	suite.Require().Equal(account.ID, got.ID)
	suite.Require().Equal(update.PersonID, got.PersonID)
	suite.Require().Equal(update.FirstName, got.FirstName)
	suite.Require().Equal(update.LastName, got.LastName)
	suite.Require().Equal(update.WebAddress, got.WebAddress)
	suite.Require().Equal(update.DateBirth, got.DateBirth)

}

func (suite *ServiceTestSuite) TestPartialUpdateBadRequest() {
	account, err := suite.queries.CreateAccount(context.Background(), database.CreateAccountParams{
		PersonID:   "7777777777",
		FirstName:  "John",
		LastName:   "Doe",
		WebAddress: "john.doe@test.local",
		DateBirth:  "1990-01-01",
	})
	suite.Require().NoError(err)

	request := apiAccount{
		PersonID: "888888888888",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(request))

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/accounts/%d", account.ID), &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)
}

func (suite *ServiceTestSuite) TestPartialUpdate() {
	account, err := suite.queries.CreateAccount(context.Background(), database.CreateAccountParams{
		PersonID:   "66666666666",
		FirstName:  "John",
		LastName:   "Doe",
		WebAddress: "john.doe@test.local",
		DateBirth:  "1990-01-01",
	})
	suite.Require().NoError(err)

	update := apiAccount{
		FirstName: "Jane",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(update))

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/accounts/%d", account.ID), &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	var got apiAccount
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&got))
	suite.Require().Equal(account.ID, got.ID)
	suite.Require().Equal(account.PersonID, got.PersonID)
	suite.Require().Equal(update.FirstName, got.FirstName)
	suite.Require().Equal(account.LastName, got.LastName)
	suite.Require().Equal(account.WebAddress, got.WebAddress)
	suite.Require().Equal(account.DateBirth, got.DateBirth)
}

func (suite *ServiceTestSuite) TestDelete() {
	account, err := suite.queries.CreateAccount(context.Background(), database.CreateAccountParams{
		PersonID:   "22222222222",
		FirstName:  "John",
		LastName:   "Doe",
		WebAddress: "john.doe@test.local",
		DateBirth:  "1990-01-01",
	})
	suite.Require().NoError(err)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/accounts/%d", account.ID), nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	_, err = suite.queries.GetAccount(context.Background(), account.ID)
	suite.Require().Error(err)
}

func (suite *ServiceTestSuite) TestList() {
	req, err := http.NewRequest("GET", "/accounts/1", nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusNotFound, rec.Result().StatusCode)
}
