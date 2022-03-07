package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/chack93/scrumpoker_api/internal/domain"
	"github.com/chack93/scrumpoker_api/internal/domain/common"
	"github.com/chack93/scrumpoker_api/internal/service/config"
	"github.com/chack93/scrumpoker_api/internal/service/database"
	"github.com/chack93/scrumpoker_api/internal/service/logger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	if err := config.Init(); err != nil {
		logrus.Fatalf("config init failed, err: %v", err)
	}
	if err := logger.Init(); err != nil {
		logrus.Fatalf("log init failed, err: %v", err)
	}
	if err := database.New().Init(); err != nil {
		logrus.Fatalf("database init failed, err: %v", err)
	}
	if err := domain.Init(); err != nil {
		logrus.Fatalf("domain init failed, err: %v", err)
	}

	os.Exit(m.Run())
}

func Request(
	method string,
	path string,
	body interface{},
) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	bodyJson, _ := json.Marshal(body)
	req := httptest.NewRequest(method, "/", bytes.NewReader(bodyJson))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath(path)
	return ctx, rec
}

func CleanModelTS(a *common.BaseModel) {
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	//for gorm soft delete feature
	//a.DeletedAt.Time = now
	//a.DeletedAt.Valid = false
}
