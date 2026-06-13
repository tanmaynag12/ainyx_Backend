package handler

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/tanmaynag12/ainyx_Backend/internal/logger"
	"github.com/tanmaynag12/ainyx_Backend/internal/models"
	"github.com/tanmaynag12/ainyx_Backend/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc      *service.UserService
	validate *validator.Validate
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:      svc,
		validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_BODY",
			"message": "request body is not valid json",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "VALIDATION_FAILED",
			"message": err.Error(),
		})
	}

	user, err := h.svc.Create(c.Context(), req)
	if err != nil {
		logger.Log.Error("failed to create user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "CREATE_FAILED",
			"message": err.Error(),
		})
	}

	logger.Log.Info("user created", zap.Int32("user_id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_ID",
			"message": "id must be a number",
		})
	}

	user, err := h.svc.GetByID(c.Context(), int32(id))
	if err != nil {
		logger.Log.Error("failed to get user", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    "USER_NOT_FOUND",
			"message": "no user with id " + strconv.Itoa(id),
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_ID",
			"message": "id must be a number",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_BODY",
			"message": "request body is not valid json",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "VALIDATION_FAILED",
			"message": err.Error(),
		})
	}

	user, err := h.svc.Update(c.Context(), int32(id), req)
	if err != nil {
		logger.Log.Error("failed to update user", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "UPDATE_FAILED",
			"message": err.Error(),
		})
	}

	logger.Log.Info("user updated", zap.Int("user_id", id))
	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_ID",
			"message": "id must be a number",
		})
	}

	if err := h.svc.Delete(c.Context(), int32(id)); err != nil {
		logger.Log.Error("failed to delete user", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "DELETE_FAILED",
			"message": err.Error(),
		})
	}

	logger.Log.Info("user deleted", zap.Int("user_id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.svc.List(c.Context())
	if err != nil {
		logger.Log.Error("failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "LIST_FAILED",
			"message": err.Error(),
		})
	}

	return c.JSON(users)
}