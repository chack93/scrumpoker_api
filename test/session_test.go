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

func TestSessionCRUD(t *testing.T) {
	var ctx echo.Context
	var rec *httptest.ResponseRecorder
	var baseURL = "/api/scrumpoker_api/session/"
	var impl = session.ServerInterfaceImpl{}

	// CREATE
	var respCreate session.Session
	var descCreate = "new session"
	var cardSelectionListCreate = "1,2,3"
	var ownerClientIdCreate = "1234"
	var gameStatusCreate = "init"
	ctx, rec = Request("POST", baseURL, session.CreateSessionJSONRequestBody{
		Description:       &descCreate,
		CardSelectionList: &cardSelectionListCreate,
		OwnerClientId:     &ownerClientIdCreate,
		GameStatus:        &gameStatusCreate,
	})
	assert.NoError(t, impl.CreateSession(ctx))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respCreate))
	assert.True(t, respCreate.JoinCode != nil)
	assert.Equal(t, 8, len(*respCreate.JoinCode))
	assert.Equal(t, descCreate, *respCreate.Description)
	assert.Equal(t, cardSelectionListCreate, *respCreate.CardSelectionList)
	assert.Equal(t, ownerClientIdCreate, *respCreate.OwnerClientId)
	assert.Equal(t, gameStatusCreate, *respCreate.GameStatus)

	// READ
	ctx, rec = Request("GET", baseURL+":id", nil)
	assert.NoError(t, impl.ReadSession(ctx, respCreate.ID.String()))
	assert.Equal(t, http.StatusOK, rec.Code)
	var respRead session.Session
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respRead))
	assert.Equal(t, *respCreate.JoinCode, *respRead.JoinCode)
	assert.Equal(t, *respCreate.Description, *respRead.Description)
	assert.Equal(t, *respCreate.CardSelectionList, *respRead.CardSelectionList)
	assert.Equal(t, *respCreate.OwnerClientId, *respRead.OwnerClientId)
	assert.Equal(t, *respCreate.GameStatus, *respRead.GameStatus)

	// READ JOINCODE
	ctx, rec = Request("GET", baseURL+"/join/:joinCode", nil)
	assert.NoError(t, impl.ReadSessionJoinCode(ctx, *respCreate.JoinCode))
	assert.Equal(t, http.StatusOK, rec.Code)
	var respReadJoinCode session.Session
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respReadJoinCode))
	assert.Equal(t, respCreate.ID.String(), respReadJoinCode.ID.String())
	assert.Equal(t, *respCreate.JoinCode, *respReadJoinCode.JoinCode)
	assert.Equal(t, *respCreate.Description, *respReadJoinCode.Description)
	assert.Equal(t, *respCreate.CardSelectionList, *respReadJoinCode.CardSelectionList)
	assert.Equal(t, *respCreate.OwnerClientId, *respReadJoinCode.OwnerClientId)
	assert.Equal(t, *respCreate.GameStatus, *respReadJoinCode.GameStatus)

	// UPDATE
	var descUpdate = "updated description"
	var cardSelectionListUpdate = "1,2,3,5,6"
	var ownerClientIdUpdate = "4321"
	var gameStatusUpdate = "reveal"
	ctx, rec = Request("PUT", baseURL+":id", session.UpdateSessionJSONRequestBody{
		Description:       &descUpdate,
		CardSelectionList: &cardSelectionListUpdate,
		OwnerClientId:     &ownerClientIdUpdate,
		GameStatus:        &gameStatusUpdate,
	})
	assert.NoError(t, impl.UpdateSession(
		ctx,
		respCreate.ID.String(),
		session.UpdateSessionParams{
			ClientId: ownerClientIdCreate,
		},
	))
	assert.Equal(t, http.StatusNoContent, rec.Code)
	// UPDATE-READ
	ctx, rec = Request("GET", baseURL+":id", nil)
	assert.NoError(t, impl.ReadSession(ctx, respCreate.ID.String()))
	assert.Equal(t, http.StatusOK, rec.Code)
	var respUpdate session.Session
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respUpdate))
	assert.Equal(t, *respCreate.JoinCode, *respUpdate.JoinCode)
	assert.Equal(t, descUpdate, *respUpdate.Description)
	assert.Equal(t, cardSelectionListUpdate, *respUpdate.CardSelectionList)
	assert.Equal(t, ownerClientIdUpdate, *respUpdate.OwnerClientId)
	assert.Equal(t, gameStatusUpdate, *respUpdate.GameStatus)

	// UPDATE BAD CLIENT
	ctx, rec = Request("PUT", baseURL+":id", session.UpdateSessionJSONRequestBody{
		Description: &descUpdate,
	})
	errRead := impl.UpdateSession(
		ctx,
		respCreate.ID.String(),
		session.UpdateSessionParams{
			ClientId: "bad",
		},
	)
	assert.Error(t, errRead)
	respError := errRead.(*echo.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, respError.Code)
}
