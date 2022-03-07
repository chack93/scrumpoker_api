package client

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServerInterfaceImpl struct{}

func (*ServerInterfaceImpl) CreateClient(ctx echo.Context) error {
	var requestBody CreateClientJSONRequestBody
	if err := ctx.Bind(&requestBody); err != nil {
		logrus.Infof("bind body failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "bad body, expected format: Client.json")
	}
	var newEntry = Client{}
	newEntry.Connected = requestBody.Connected
	newEntry.Estimation = requestBody.Estimation
	newEntry.Name = requestBody.Name
	newEntry.SessionId = requestBody.SessionId
	newEntry.Viewer = requestBody.Viewer
	if err := CreateClient(&newEntry); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create")
	}
	return ctx.JSON(http.StatusOK, newEntry)
}

func (*ServerInterfaceImpl) ReadClient(ctx echo.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad id, expected format: uuid")
	}
	var response Client
	if err := ReadClient(uuid, &response); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read")
	}
	return ctx.JSON(http.StatusOK, response)
}

func (*ServerInterfaceImpl) UpdateClient(ctx echo.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad id, expected format: uuid")
	}
	var client Client
	if err := ReadClient(uuid, &client); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read")
	}

	var requestBody UpdateClientJSONRequestBody
	if err := ctx.Bind(&requestBody); err != nil {
		logrus.Infof("bind body failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "bad body, expected format: Client.json")
	}
	client.Connected = requestBody.Connected
	client.Estimation = requestBody.Estimation
	client.Name = requestBody.Name
	client.SessionId = requestBody.SessionId
	client.Viewer = requestBody.Viewer
	if err := UpdateClient(uuid, &client); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update")
	}
	return ctx.NoContent(http.StatusNoContent)
}
