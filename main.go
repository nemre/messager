// @title Messager API
// @version 1.0
// @description Fast & Strong Messaging Tool

// @host localhost:2025
// @BasePath /

package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	messageservice "messager/application/service/message"
	"messager/infrastructure/client"
	"messager/infrastructure/config"
	"messager/infrastructure/database/postgresql"
	"messager/infrastructure/database/redis"
	"messager/infrastructure/logger"
	messagepersistence "messager/infrastructure/persistence/message"
	messageconsumer "messager/presentation/consumer/message"
	messagehandler "messager/presentation/handler/message"
	messagejob "messager/presentation/job/message"

	"messager/infrastructure/server"
)

func main() {
	time.Local = time.UTC
	logger := logger.New()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("failed to initialize config", err)
	}

	srv := server.New(server.Config{
		Host:     cfg.GetServer().Host,
		Port:     cfg.GetServer().Port,
		IDHeader: cfg.GetServer().IDHeader,
		OnListen: func(address string) {
			logger.Debug("server started", "address", address)
		},
		OnRequestStart: func(ctx server.RequestContext) {
			logger.Debug("request received", "id", ctx.GetID(), "method", ctx.GetMethod(), "url", ctx.GetURL())
		},
		OnRequestEnd: func(ctx server.RequestContext, status uint16) {
			logger.Debug("request succeeded", "id", ctx.GetID(), "status", status)
		},
		OnRequestError: func(ctx server.RequestContext, status uint16, message string, err error) {
			logger.Error("request failed", err, "id", ctx.GetID(), "status", status, "cause", message)
		},
		OnRequestPanic: func(ctx server.RequestContext, status uint16, message string, err error, stackTrace string) {
			logger.FatalWithoutExit("request panicked", err, "id", ctx.GetID(), "status", status, "cause", message, "stacktrace", stackTrace)
		},
	})

	router := srv.NewRouter()

	router.AddRoute("GET /health", func(ctx server.RequestContext) (any, error) {
		return map[string]string{
			"status": "green",
		}, nil
	})

	postgreSQL, err := postgresql.New(postgresql.Config{
		Host:     cfg.GetPostgreSQL().Host,
		Port:     cfg.GetPostgreSQL().Port,
		User:     cfg.GetPostgreSQL().User,
		Password: cfg.GetPostgreSQL().Password,
		Name:     cfg.GetPostgreSQL().Name,
	})
	if err != nil {
		logger.Fatal("failed to initialize postgresql", err)
	}

	redis, err := redis.New(redis.Config{
		Host:     cfg.GetRedis().Host,
		Port:     cfg.GetRedis().Port,
		User:     cfg.GetRedis().User,
		Password: cfg.GetRedis().Password,
		DB:       cfg.GetRedis().DB,
	})
	if err != nil {
		logger.Fatal("failed to initialize redis", err)
	}

	messageRepository, err := messagepersistence.New(postgreSQL, redis)
	if err != nil {
		logger.Fatal("failed to initialize message repository", err)
	}

	client := client.New(client.Config{
		URL:     cfg.GetClient().URL,
		Token:   cfg.GetClient().Token,
		Timeout: cfg.GetClient().Timeout,
	})

	messageService := messageservice.New(messageRepository, client)

	messageJob := messagejob.New(messageService, cfg.GetJob().Interval, func(err error) {
		logger.FatalWithoutExit("message job failed", err)
	})

	messageConsumer, err := messageconsumer.New(
		messageService,
		cfg.GetKafka().Brokers,
		cfg.GetKafka().GroupID,
		cfg.GetKafka().Topic,
		func(err error) {
			logger.FatalWithoutExit("message consume failed", err)
		},
	)
	if err != nil {
		logger.Fatal("failed to initialize message consumer", err)
	}

	_ = messagehandler.New(router, messageService, messageJob)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			logger.FatalWithoutExit("failed to start server", err)
			stop <- syscall.SIGINT
		}
	}()

	go func() {
		messageConsumer.Start()
	}()

	<-stop

	if err := srv.Stop(); err != nil {
		logger.FatalWithoutExit("failed to stop server", err)
	}

	if err := messageConsumer.Stop(); err != nil {
		logger.FatalWithoutExit("failed to stop message consumer", err)
	}

	messageJob.Stop()
	postgreSQL.Close()

	if err := redis.Close(); err != nil {
		logger.FatalWithoutExit("failed to stop redis", err)
	}

	logger.Debug("application gracefully stopped")
}
