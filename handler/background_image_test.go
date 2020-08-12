package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"local.packages/entity"
	"local.packages/repository"
	"local.packages/utils"
)

func TestShouldReturnsBackgroundImageInstancesAndStatusOKAsHTTPResponse(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBackgroundImageHandler(repository.NewBackgroundImageRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/background_images", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	backgroundImageID := uint(1)
	backgroundImageURL := "http://localhost/image"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `background_images`")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url"}).AddRow(backgroundImageID, backgroundImageURL))

	r.GET("/background_images", bh.IndexBackgroundImage)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]entity.BackgroundImage{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Len(t, res["background_images"], 1)
	assert.Equal(t, res["background_images"][0].ID, backgroundImageID)
	assert.Equal(t, res["background_images"][0].URL, backgroundImageURL)
}
