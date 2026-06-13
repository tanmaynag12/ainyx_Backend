package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/tanmaynag12/ainyx_Backend/config"
	"github.com/tanmaynag12/ainyx_Backend/internal/handler"
	"github.com/tanmaynag12/ainyx_Backend/internal/logger"
	"github.com/tanmaynag12/ainyx_Backend/internal/repository"
	"github.com/tanmaynag12/ainyx_Backend/internal/routes"
	"github.com/tanmaynag12/ainyx_Backend/internal/service"
)

func main() {
	logger.Init()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			})
		},
	})

	routes.Setup(app, userHandler)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}