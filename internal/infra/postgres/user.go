package postgres

import (
	"context"
	"database/sql"
	"errors"
	"gravitum/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

type UserData struct {
	ID        uuid.UUID
	Login     string
	FirstName string
	LastName  string
}

func (r *UserRepository) Create(ctx context.Context, u user.User) error {
	_, err := r.pool.Exec(ctx,
		"insert into users (id, login, first_name, last_name) values ($1, $2, $3, $4) ",
		u.ID(), u.Login(), u.FirsName(), u.LastName(),
	)
	return err
}

func (r *UserRepository) Update(ctx context.Context, u user.User) error {
	_, err := r.pool.Exec(ctx,
		"update users SET login = $1, first_name = $2, last_name = $3 where ID = $4",
		u.Login(), u.FirsName(), u.LastName(), u.ID(),
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, u user.User) error {
	_, err := r.pool.Exec(ctx,
		"delete from users where id = $1",
		u.ID(),
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	var data UserData

	err := r.pool.QueryRow(ctx, "select id, login, first_name, last_name from users where id=$1", id).Scan(
		&data.ID,
		&data.Login,
		&data.FirstName,
		&data.LastName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, user.ErrNotFound
	}

	if err != nil {
		// var pgErr *pgconn.PgError
		// if errors.As(err, &pgErr) && pgErr.Code == "20000" {
		// 	return user.User{}, user.ErrNotFound
		// }

		return user.User{}, err
	}

	return user.New(data.ID, data.Login, data.FirstName, data.LastName)
}
