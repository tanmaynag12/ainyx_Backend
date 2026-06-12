package repository

import (
	"context"
	"database/sql"
	"time"

	db "github.com/tanmaynag12/ainyx_Backend/db/sqlc"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{
		queries: db.New(database),
	}
}

func (r *UserRepository) Create(ctx context.Context, name string, dob time.Time) (db.User, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

func (r *UserRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *UserRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (db.User, error) {
	return r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
}

func (r *UserRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *UserRepository) List(ctx context.Context) ([]db.User, error) {
	return r.queries.ListUsers(ctx)
}