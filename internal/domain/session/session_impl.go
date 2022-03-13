package session

import (
	"crypto/sha256"
	"encoding/base32"
	"math/big"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServerInterfaceImpl struct{}

func (*ServerInterfaceImpl) CreateSession(ctx echo.Context) error {
	var requestBody CreateSessionJSONRequestBody
	if err := ctx.Bind(&requestBody); err != nil {
		logrus.Infof("bind body failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "bad body, expected format: Session.json")
	}
	var newEntry = Session{}
	hash := sha256.New().Sum(big.NewInt(time.Now().UnixNano()).Bytes())
	joinCode := base32.StdEncoding.EncodeToString(hash)[:8]
	newEntry.JoinCode = &joinCode
	newEntry.Description = requestBody.Description
	newEntry.CardSelectionList = requestBody.CardSelectionList
	newEntry.OwnerClientId = requestBody.OwnerClientId
	newEntry.GameStatus = requestBody.GameStatus
	if err := CreateSession(&newEntry); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create")
	}
	return ctx.JSON(http.StatusOK, newEntry)
}

func (*ServerInterfaceImpl) ReadSession(ctx echo.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad id, expected format: uuid")
	}
	var response Session
	if err := ReadSession(uuid, &response); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read")
	}
	return ctx.JSON(http.StatusOK, response)
}

func (*ServerInterfaceImpl) ReadSessionJoinCode(ctx echo.Context, joinCode string) error {
	var response Session
	if err := ReadSessionJoinCode(joinCode, &response); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read")
	}
	return ctx.JSON(http.StatusOK, response)
}

func (*ServerInterfaceImpl) UpdateSession(ctx echo.Context, id string, params UpdateSessionParams) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad id, expected format: uuid")
	}
	var session Session
	if err := ReadSession(uuid, &session); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read")
	}
	if session.OwnerClientId != nil && *session.OwnerClientId != params.ClientId {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var requestBody UpdateSessionJSONRequestBody
	if err := ctx.Bind(&requestBody); err != nil {
		logrus.Infof("bind body failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "bad body, expected format: Session.json")
	}
	session.Description = requestBody.Description
	session.CardSelectionList = requestBody.CardSelectionList
	session.OwnerClientId = requestBody.OwnerClientId
	session.GameStatus = requestBody.GameStatus
	if err := UpdateSession(uuid, &session); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update")
	}
	return ctx.NoContent(http.StatusNoContent)
}
