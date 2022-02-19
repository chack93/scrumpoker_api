package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chack93/scrumpoker_api/internal/domain/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSessionCLRUD(t *testing.T) {
	var ctx echo.Context
	var rec *httptest.ResponseRecorder
	var baseURL = "/api/scrumpoker_api/session/"
	var impl = session.ServerInterfaceImpl{}

	// CREATE
	var respCreate session.Session
	descCreate := "new session"
	ctx, rec = Request("POST", baseURL, session.CreateSessionJSONRequestBody{
		Description: &descCreate,
	})
	assert.NoError(t, impl.CreateSession(ctx))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respCreate))
	assert.True(t, respCreate.JoinCode != nil)
	assert.Equal(t, 8, len(*respCreate.JoinCode))
	assert.Equal(t, descCreate, *respCreate.Description)

	// LIST
	ctx, rec = Request("GET", baseURL, nil)
	assert.NoError(t, impl.ListSession(ctx, session.ListSessionParams{}))
	var respList []session.Session
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respList))
	assert.Equal(t, 1, len(respList))
	assert.Equal(t, *respCreate.JoinCode, *respList[0].JoinCode)
	assert.Equal(t, *respCreate.Description, *respList[0].Description)

	// READ
	ctx, rec = Request("GET", baseURL+":id", nil)
	assert.NoError(t, impl.ReadSession(ctx, respCreate.ID.String()))
	assert.Equal(t, http.StatusOK, rec.Code)
	var respRead session.Session
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respRead))
	assert.Equal(t, *respCreate.JoinCode, *respRead.JoinCode)
	assert.Equal(t, *respCreate.Description, *respRead.Description)

	// UPDATE
	var descUpdate = "updated description"
	ctx, rec = Request("PUT", baseURL+":id", session.UpdateSessionJSONRequestBody{
		Description: &descUpdate,
	})
	assert.NoError(t, impl.UpdateSession(ctx, respCreate.ID.String()))
	assert.Equal(t, http.StatusNoContent, rec.Code)
	// UPDATE-READ
	ctx, rec = Request("GET", baseURL+":id", nil)
	assert.NoError(t, impl.ReadSession(ctx, respCreate.ID.String()))
	assert.Equal(t, http.StatusOK, rec.Code)
	var respUpdate session.Session
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respUpdate))
	assert.Equal(t, *respCreate.JoinCode, *respUpdate.JoinCode)
	assert.Equal(t, descUpdate, *respUpdate.Description)

	// DELETE
	ctx, rec = Request("DELETE", baseURL+":id", nil)
	assert.NoError(t, impl.DeleteSession(ctx, respCreate.ID.String()))
	assert.Equal(t, http.StatusNoContent, rec.Code)
	// DELETE-READ
	ctx, rec = Request("GET", baseURL+":id", nil)
	errRead := impl.ReadSession(ctx, respCreate.ID.String())
	assert.Error(t, errRead)
	respError := errRead.(*echo.HTTPError)
	assert.Equal(t, http.StatusNotFound, respError.Code)
}
