package history

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ServerInterfaceImpl struct{}

func (*ServerInterfaceImpl) ListHistory(ctx echo.Context, sessionID string) error {
	uuid, err := uuid.Parse(sessionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad id, expected format: uuid")
	}
	var response []History
	if err := ListHistoryBySessionID(uuid, &response); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read")
	}
	return ctx.JSON(http.StatusOK, HistoryList(response))
}
