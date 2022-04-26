package server

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/chack93/scrumpoker_api/internal/domain"
	"github.com/chack93/scrumpoker_api/internal/service/websocket"
	"github.com/gobwas/ws"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	echo *echo.Echo
}

var server *Server

func Get() *Server {
	return server
}

func New() *Server {
	server = &Server{}
	return server
}

//go:embed swagger/swagger_gen.yaml
var swaggerYaml []byte

//go:embed swagger/index.html
var swaggerHtml []byte

func (srv *Server) Init(wg *sync.WaitGroup) error {
	srv.echo = echo.New()
	srv.echo.HideBanner = true
	srv.echo.HidePort = true
	srv.echo.Use(middleware.Logger())
	srv.echo.Use(middleware.Recover())

	baseURL := "/api/scrumpoker"
	apiGroup := srv.echo.Group(baseURL)
	apiGroup.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Status string
		}{"ok"})
	})
	apiGroup.GET("/doc", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "doc/")
	})
	apiGroup.GET("/doc/", func(c echo.Context) error {
		return c.HTMLBlob(http.StatusOK, swaggerHtml)
	})
	apiGroup.GET("/doc/swagger.yaml", func(c echo.Context) error {
		return c.HTMLBlob(http.StatusOK, swaggerYaml)
	})
	apiGroup.GET("/ws/:clientId/:groupId", func(c echo.Context) error {
		clientID := c.Param("clientId")
		groupID := c.Param("groupId")
		conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response().Writer)
		if err != nil {
			logrus.Errorf("upgrade http failed: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to upgrade connection to websocket")
		}
		websocket.CreateHandler(conn, clientID, groupID)
		return nil
	})
	domain.RegisterHandlers(srv.echo, baseURL)

	address := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
	go func() {
		if err := srv.echo.Start(address); err != nil && err != http.ErrServerClosed {
			logrus.Warnf("server start failed, err: %v", err)
			wg.Done()
		}
	}()

	defer func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.echo.Shutdown(ctx); err != nil {
			logrus.Errorf("server shutdown failed, err: %v", err)
		}
		logrus.Info("server shutdown")
		wg.Done()
	}()

	for _, el := range srv.echo.Routes() {
		lastSlash := strings.LastIndex(el.Name, "/")
		domainHandler := el.Name[lastSlash:]
		logrus.Infof("%6s %s -> %s", el.Method, el.Path, domainHandler)
	}
	logrus.Infof("http server started on %s", address)
	return nil
}
