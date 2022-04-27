package main

import (
	"sync"

	"github.com/chack93/scrumpoker_api/internal/domain"
	"github.com/chack93/scrumpoker_api/internal/service/cleanup"
	"github.com/chack93/scrumpoker_api/internal/service/config"
	"github.com/chack93/scrumpoker_api/internal/service/database"
	"github.com/chack93/scrumpoker_api/internal/service/datasync"
	"github.com/chack93/scrumpoker_api/internal/service/logger"
	"github.com/chack93/scrumpoker_api/internal/service/msgsystem"
	"github.com/chack93/scrumpoker_api/internal/service/server"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.Init(); err != nil {
		logrus.Fatalf("config init failed, err: %v", err)
	}
	if err := logger.Init(); err != nil {
		logrus.Fatalf("log init failed, err: %v", err)
	}
	if err := msgsystem.New().Init(); err != nil {
		logrus.Fatalf("msgsystem init failed, err: %v", err)
	}
	if err := database.New().Init(); err != nil {
		logrus.Fatalf("database init failed, err: %v", err)
	}
	if err := domain.Init(); err != nil {
		logrus.Fatalf("domain init failed, err: %v", err)
	}
	if err := datasync.Init(); err != nil {
		logrus.Fatalf("datasync init failed, err: %v", err)
	}
	if err := cleanup.Init(); err != nil {
		logrus.Fatalf("cleanup init failed, err: %v", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	if err := server.New().Init(wg); err != nil {
		logrus.Fatalf("server init failed, err: %v", err)
	}

	wg.Wait()
}
