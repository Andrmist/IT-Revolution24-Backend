package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"itrevolution-backend/internal"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/job"
	"itrevolution-backend/internal/types"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger := logrus.New()

	config, err := types.InitConfig()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		config.PostgresHost, config.PostgresUser, config.PostgresPassword, config.PostgresDB, config.PostgresPort)), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&domain.User{}, &domain.Pet{}, &domain.Message{})

	serverCtx := types.ServerContext{
		Config:  config,
		Log:     logger,
		DB:      db,
		WsConns: make(map[uint][]*websocket.Conn),
	}

	internal.Run(ctx, serverCtx)

	cr := cron.New()
	j := job.NewJob(cr, db, serverCtx.WsConns)

	go j.Run()
	defer cr.Stop()

	logger.Info("Start...")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	logger.Info("Shutdown...")

	cancel()
	os.Exit(0)
}
