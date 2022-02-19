package main

import (
	"sync"

	"github.com/chack93/go_base/internal/domain"
	"github.com/chack93/go_base/internal/service/config"
	"github.com/chack93/go_base/internal/service/database"
	"github.com/chack93/go_base/internal/service/logger"
	"github.com/chack93/go_base/internal/service/server"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.Init(); err != nil {
		logrus.Fatalf("config init failed, err: %v", err)
	}
	if err := logger.Init(); err != nil {
		logrus.Fatalf("log init failed, err: %v", err)
	}
	if err := database.New().Init(); err != nil {
		logrus.Fatalf("database init failed, err: %v", err)
	}
	if err := domain.DbMigrate(); err != nil {
		logrus.Fatalf("domain init failed, err: %v", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	if err := server.New().Init(wg); err != nil {
		logrus.Fatalf("server init failed, err: %v", err)
	}

	wg.Wait()
}
