package main

import (
	"context"
	"fmt"
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

	if err := domain.MigrateDB(db); err != nil {
		panic(err)
	}

	serverCtx := types.ServerContext{
		Config: config,
		Log:    logger,
		DB:     db,
	}

	internal.Run(ctx, serverCtx)

	cr := cron.New()
	job := job.NewJob(cr, db)

	go job.Run()
	defer cr.Stop()

	logger.Info("Start...")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	logger.Info("Shutdown...")

	cancel()
	os.Exit(0)
}
