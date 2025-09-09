package container

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/frkntplglu/insider/internal/config"
	"github.com/frkntplglu/insider/internal/handler"
	"github.com/frkntplglu/insider/internal/repositories"
	"github.com/frkntplglu/insider/internal/services"
	"github.com/frkntplglu/insider/pkg/database"
	"github.com/frkntplglu/insider/pkg/logger"
	"github.com/frkntplglu/insider/pkg/redis"
	smsclient "github.com/frkntplglu/insider/pkg/sms_client"
	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Container struct {
	Config         *config.Config
	Database       *database.Database
	RedisClient    *redis.RedisClient
	MessageRepo    *repositories.MessageRepository
	MessageService *services.MessageService
	MessageHandler *handler.MessageHandler
	App            *fiber.App
}

func NewContainer() *Container {
	logger.Init(slog.LevelInfo)
	logger.Info("Initializing application container")

	cfg := config.LoadConfig()

	db := database.NewDatabase(database.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
	})

	redisClient := redis.NewRedisClient(redis.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Database: cfg.Redis.Database,
	})

	messageRepo := repositories.NewMessageRepository(db)

	smsClient := smsclient.NewSmsClient(cfg.SMS.Host)

	messageService := services.NewMessageService(messageRepo, redisClient, smsClient, cfg.Redis.Key, cfg.Ticker.Period)

	messageHandler := handler.NewMessageHandler(messageService)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(recover.New())

	messageHandler.SetRoutes(app)

	sw := swagno.New(swagno.Config{Title: cfg.App.Name, Version: cfg.App.Version, Host: fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)})

	sw.AddEndpoints(handler.MessageSwaggerEndpoints)

	swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

	if err := messageService.StartAutoSending(); err != nil {
		logger.Warn("Could not start auto sending", "error", err)
	}

	return &Container{
		Config:         cfg,
		Database:       db,
		RedisClient:    redisClient,
		MessageRepo:    messageRepo,
		MessageService: messageService,
		MessageHandler: messageHandler,
		App:            app,
	}
}

func (c *Container) Start() error {
	addr := c.Config.Server.Host + ":" + c.Config.Server.Port
	logger.Info("Starting server", "address", addr)
	return c.App.Listen(addr)
}

func (c *Container) Stop(ctx context.Context) error {
	logger.Info("Shutting down application")
	c.Database.Close()
	c.RedisClient.Close()
	return c.App.ShutdownWithContext(ctx)
}
