package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chack93/scrumpoker_api/internal/domain/history"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHistoryCRUD(t *testing.T) {
	var ctx echo.Context
	var rec *httptest.ResponseRecorder
	var baseURL = "/api/scrumpoker_api/history/"
	var impl = history.ServerInterfaceImpl{}

	// CREATE
	var clientID = "1234"
	var clientName = "John Doe"
	var estimation = "3"
	var sessionID = uuid.New().String()
	var gameID = uuid.New().String()

	var createRequest = history.History{
		HistoryNew: history.HistoryNew{
			ClientId:   &clientID,
			ClientName: &clientName,
			Estimation: &estimation,
			SessionId:  &sessionID,
		},
		GameId: &gameID,
	}

	assert.NoError(t, history.CreateHistory(&createRequest))

	// READ
	ctx, rec = Request("GET", baseURL, nil)
	assert.NoError(t, impl.ListHistory(ctx, sessionID))
	assert.Equal(t, http.StatusOK, rec.Code)
	var respRead history.HistoryList
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respRead))
	assert.Equal(t, 1, len(respRead))
}
