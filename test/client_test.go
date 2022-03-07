package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chack93/scrumpoker_api/internal/domain/client"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestClientCRUD(t *testing.T) {
	var ctx echo.Context
	var rec *httptest.ResponseRecorder
	var baseURL = "/api/scrumpoker_api/client/"
	var impl = client.ServerInterfaceImpl{}

	// CREATE
	var respCreate client.Client
	var connected = false
	var estimation = "3"
	var name = "John Doe"
	var sessionId = "1234"
	var viewer = false
	var createRequest = client.CreateClientJSONRequestBody{
		Connected:  &connected,
		Estimation: &estimation,
		Name:       &name,
		SessionId:  &sessionId,
		Viewer:     &viewer,
	}
	ctx, rec = Request("POST", baseURL, createRequest)
	assert.NoError(t, impl.CreateClient(ctx))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respCreate))
	assert.True(t, respCreate.ID.String() != "")
	assert.Equal(t, *createRequest.Connected, *respCreate.Connected)
	assert.Equal(t, *createRequest.Estimation, *respCreate.Estimation)
	assert.Equal(t, *createRequest.Name, *respCreate.Name)
	assert.Equal(t, *createRequest.SessionId, *respCreate.SessionId)
	assert.Equal(t, *createRequest.Viewer, *respCreate.Viewer)

	// READ
	ctx, rec = Request("GET", baseURL+":id", nil)
	assert.NoError(t, impl.ReadClient(ctx, respCreate.ID.String()))
	assert.Equal(t, http.StatusOK, rec.Code)
	var respRead client.Client
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respRead))
	assert.Equal(t, *createRequest.Connected, *respRead.Connected)
	assert.Equal(t, *createRequest.Estimation, *respRead.Estimation)
	assert.Equal(t, *createRequest.Name, *respRead.Name)
	assert.Equal(t, *createRequest.SessionId, *respRead.SessionId)
	assert.Equal(t, *createRequest.Viewer, *respRead.Viewer)

	// UPDATE
	connected = true
	estimation = "12"
	name = "Jane Viewer"
	sessionId = "4321"
	viewer = true
	var updateRequest = client.UpdateClientJSONRequestBody{
		Connected:  &connected,
		Estimation: &estimation,
		Name:       &name,
		SessionId:  &sessionId,
		Viewer:     &viewer,
	}
	ctx, rec = Request("PUT", baseURL+":id", updateRequest)
	assert.NoError(t, impl.UpdateClient(
		ctx,
		respCreate.ID.String(),
	))
	assert.Equal(t, http.StatusNoContent, rec.Code)
	// UPDATE-READ
	ctx, rec = Request("GET", baseURL+":id", nil)
	assert.NoError(t, impl.ReadClient(ctx, respCreate.ID.String()))
	assert.Equal(t, http.StatusOK, rec.Code)
	var respUpdate client.Client
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respUpdate))
	assert.Equal(t, *createRequest.Connected, *respUpdate.Connected)
	assert.Equal(t, *createRequest.Estimation, *respUpdate.Estimation)
	assert.Equal(t, *createRequest.Name, *respUpdate.Name)
	assert.Equal(t, *createRequest.SessionId, *respUpdate.SessionId)
	assert.Equal(t, *createRequest.Viewer, *respUpdate.Viewer)

	// READ NOT FOUND
	ctx, rec = Request("GET", baseURL+":id", nil)
	errRead := impl.ReadClient(ctx, uuid.New().String())
	assert.Error(t, errRead)
	respError := errRead.(*echo.HTTPError)
	assert.Equal(t, http.StatusNotFound, respError.Code)
}
