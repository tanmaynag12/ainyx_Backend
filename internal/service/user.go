package service

import (
	"context"
	"fmt"
	"time"

	db "github.com/tanmaynag12/ainyx_Backend/db/sqlc"
	"github.com/tanmaynag12/ainyx_Backend/internal/models"
	"github.com/tanmaynag12/ainyx_Backend/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("invalid dob format, use YYYY-MM-DD")
	}

	user, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toResponse(user, false), nil
}

func (s *UserService) GetByID(ctx context.Context, id int32) (models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toResponse(user, true), nil
}

func (s *UserService) Update(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("invalid dob format, use YYYY-MM-DD")
	}

	user, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toResponse(user, false), nil
}

func (s *UserService) Delete(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.UserResponse
	for _, u := range users {
		result = append(result, toResponse(u, true))
	}
	return result, nil
}

func calculateAge(dob time.Time) int {
	today := time.Now()
	age := today.Year() - dob.Year()
	if today.Month() < dob.Month() || (today.Month() == dob.Month() && today.Day() < dob.Day()) {
		age--
	}
	return age
}

func daysUntilNextBirthday(dob time.Time) int {
	today := time.Now()
	next := time.Date(today.Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, time.Local)
	if !next.After(today) {
		next = next.AddDate(1, 0, 0)
	}
	return int(next.Sub(today).Hours() / 24)
}

func toResponse(user db.User, withAge bool) models.UserResponse {
	resp := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format("2006-01-02"),
	}

	if withAge {
		age := calculateAge(user.Dob)
		days := daysUntilNextBirthday(user.Dob)
		resp.Age = &age
		resp.NextBirthday = &days
	}

	return resp
}