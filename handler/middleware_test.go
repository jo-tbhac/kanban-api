package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"local.packages/repository"
	"local.packages/utils"
	"local.packages/validator"
)

func TestShouldSetUserIDToContextUponSuccessfullyAuthenticate(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	userID := uint(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	uh := NewUserHandler(repository.NewUserRepository(db))

	r.Use(uh.Authenticate())
	r.GET("/test", func(c *gin.Context) {
		assert.Equal(t, c.Keys["uid"], userID)
		c.Status(200)
	})

	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, w.Code, 200)
}

func TestShouldReturnsStatusUnAuthorizationWhenRequestTokenIsInvalid(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WillReturnError(gorm.ErrRecordNotFound)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	uh := NewUserHandler(repository.NewUserRepository(db))

	r.Use(uh.Authenticate())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, w.Code, 401)
}

func TestShouldSetIDToContextWhenRequestParamKeyContainsSuffixID(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	boardID := uint(1)
	cardID := uint(2)

	r.Use(MapIDParamsToContext())
	r.GET("/test/:boardID/:cardID", func(c *gin.Context) {
		assert.Equal(t, c.Keys["boardID"], boardID)
		assert.Equal(t, c.Keys["cardID"], cardID)
		c.Status(200)
	})

	url := fmt.Sprintf("/test/%d/%d", boardID, cardID)

	c.Request, _ = http.NewRequest(http.MethodGet, url, nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, w.Code, 200)
}

func TestShouldReturnsStatusBadRequestWhenRequestParamKeyDoesNotContainsSuffixID(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	r.Use(MapIDParamsToContext())
	r.GET("/test/:boardID", func(c *gin.Context) {
		c.Status(200)
	})

	c.Request, _ = http.NewRequest(http.MethodGet, "/test/sfreibdsd", nil)
	r.ServeHTTP(w, c.Request)

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Contains(t, res["errors"][0].Text, "must be an integer")
}
