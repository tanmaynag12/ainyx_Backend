package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tanmaynag12/ainyx_Backend/internal/handler"
	"github.com/tanmaynag12/ainyx_Backend/internal/middleware"
)

func Setup(app *fiber.App, h *handler.UserHandler) {
	app.Use(middleware.RequestLogger())
	app.Use(middleware.RequestID())

	api := app.Group("/api/v1")

	users := api.Group("/users")
	users.Post("/", h.CreateUser)
	users.Get("/", h.ListUsers)
	users.Get("/:id", h.GetUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
}