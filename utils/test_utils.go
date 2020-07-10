package utils

import (
	"database/sql/driver"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AnyTime is a mock for time.Time.
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface.
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// NewDBMock creates sqlmock database connection and a mock to manage expectations.
func NewDBMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gdb, err := gorm.Open("mysql", db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	gdb.LogMode(true)

	return gdb, mock
}

// SetUpRouter returns an instance of Engine that was set test mode.
func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

// SetUpAuthentication attached middleware to an instance of Engine.
// add a session token to request header.
// and add an expectation of authenticate user.
func SetUpAuthentication(r *gin.Engine, req *http.Request, mock sqlmock.Sqlmock, middleware ...gin.HandlerFunc) {
	token := "sampletoken"
	req.Header.Add("Authorization", token)

	for _, m := range middleware {
		r.Use(m)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(token).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
}
