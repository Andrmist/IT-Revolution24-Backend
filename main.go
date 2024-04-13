package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"itrevolution-backend/internal"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"os"
	"os/signal"
	"syscall"
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

	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Fish{})

	serverCtx := types.ServerContext{
		Config: config,
		Log:    logger,
		DB:     db,
	}

	internal.Run(ctx, serverCtx)

	cr := cron.New()
	//go jobs.ServerPingJob(serverCtx)
	cr.AddFunc("0 0-59/5 * * * *", func() {
	})
	cr.Start()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit
	logger.Info("Shutdown...")
	cancel()
	cr.Stop()
	os.Exit(0)
}
