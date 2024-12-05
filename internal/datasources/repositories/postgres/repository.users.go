package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/services"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"time"
)

type postgresUsersRepository struct {
	db *sqlx.DB
}

func NewPostgresUsersRepository(db *sqlx.DB) services.UsersRepository {
	return &postgresUsersRepository{db}
}

func (r *postgresUsersRepository) FindAll() ([]records.Users, error) {
	query, args, err := squirrel.Select("*").From("users").ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindAll - squirrel.Select: %w", err))
	}

	var users []records.Users
	if err := r.db.Select(&users, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindAll - db.Select: %w", err))
	}

	return users, nil
}

func (r *postgresUsersRepository) FindByID(id int) (records.Users, error) {
	query, args, err := squirrel.Select("*").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return records.Users{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindByID - squirrel.Select: %w", err))
	}

	var user records.Users
	if err := r.db.Get(&user, query, args...); err != nil {
		return records.Users{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindByID - db.Get: %w", err))
	}

	return user, nil
}

func (r *postgresUsersRepository) FindByEmail(email string) (records.Users, error) {
	query, args, err := squirrel.Select("*").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"email": email}).ToSql()
	if err != nil {
		return records.Users{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindByEmail - squirrel.Select: %w", err))
	}

	var user records.Users
	if err := r.db.Get(&user, query, args...); err != nil {
		return records.Users{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindByEmail - db.Get: %w", err))
	}

	return user, nil
}

func (r *postgresUsersRepository) Save(user records.Users) error {
	query, args, err := squirrel.Insert("users").
		Columns("email", "password", "created_at", "updated_at", "deleted_at").
		Values(user.Email, user.Password, time.Now(), nil, nil).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - Save - squirrel.Insert: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - Save - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresUsersRepository) Update(id int, users map[string]interface{}) error {
	updateQuery := squirrel.Update("users").PlaceholderFormat(squirrel.Dollar)

	for key, value := range users {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - Update - squirrel.Update: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresUsersRepository) Delete(id int) error {
	query, args, err := squirrel.Delete("users").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - Delete - squirrel.Delete: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - Delete - db.Exec: %w", err))
	}

	return nil
}
