package handler

import (
	"encoding/json"
	"fmt"
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

func TestIndexFilesHandlerShouldReturnsStatusOKWithFileData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	fh := NewFileHandler(repository.NewFileRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	boardID := uint(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/board/%d/files", boardID), nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mockFile := entity.File{
		ID:          uint(3),
		DisplayName: "file1",
		URL:         "http://test",
		ContentType: "image/png",
		CardID:      uint(4),
	}

	query := utils.ReplaceQuotationForQuery(`
		SELECT files.id, files.display_name, files.url, files.content_type, files.card_id
		FROM 'files'`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "display_name", "url", "content_type", "card_id"}).
				AddRow(mockFile.ID, mockFile.DisplayName, mockFile.URL, mockFile.ContentType, mockFile.CardID))

	r.GET("/board/:boardID/files", fh.IndexFiles)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]entity.File{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["files"][0].ID, mockFile.ID)
	assert.Equal(t, res["files"][0].DisplayName, mockFile.DisplayName)
	assert.Equal(t, res["files"][0].URL, mockFile.URL)
	assert.Equal(t, res["files"][0].ContentType, mockFile.ContentType)
	assert.Equal(t, res["files"][0].CardID, mockFile.CardID)
}
