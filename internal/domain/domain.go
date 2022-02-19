package domain

import (
	"github.com/chack93/scrumpoker_api/internal/domain/session"
	"github.com/chack93/scrumpoker_api/internal/service/database"
	"github.com/labstack/echo/v4"
)

func DbMigrate() error {
	database.Get().AutoMigrate(
		&session.Session{},
	)
	return nil
}

func RegisterHandlers(e *echo.Echo, baseURL string) {
	session.RegisterHandlersWithBaseURL(e, &session.ServerInterfaceImpl{}, baseURL)
}
