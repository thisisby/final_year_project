package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
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
	query, args, err := squirrel.
		Select("*").
		From("users").
		OrderBy("id ASC").
		ToSql()
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

func (r *postgresUsersRepository) FindByUsername(username string) (records.Users, error) {
	query, args, err := squirrel.Select("*").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"username": username}).ToSql()
	if err != nil {
		return records.Users{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindByUsername - squirrel.Select: %w", err))
	}

	var user records.Users
	if err := r.db.Get(&user, query, args...); err != nil {
		return records.Users{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindByUsername - db.Get: %w", err))
	}

	return user, nil
}

func (r *postgresUsersRepository) Save(user records.Users) error {
	query, args, err := squirrel.Insert("users").
		Columns("username", "bio", "avatar", "email", "password", "created_at", "updated_at", "deleted_at").
		Values(user.Username, user.Bio, user.Avatar, user.Email, user.Password, time.Now(), nil, nil).
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

func (r *postgresUsersRepository) ChangeAvatar(id int, avatar string) error {
	query, args, err := squirrel.Update("users").
		PlaceholderFormat(squirrel.Dollar).
		Set("avatar", avatar).
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - ChangeAvatar - squirrel.Update: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - ChangeAvatar - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresUsersRepository) UsernameExists(username string) (bool, error) {
	query, args, err := squirrel.Select("count(*)").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"username": username}).ToSql()
	if err != nil {
		return false, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - UsernameExists - squirrel.Select: %w", err))
	}

	var count int
	if err := r.db.Get(&count, query, args...); err != nil {
		return false, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - UsernameExists - db.Get: %w", err))
	}

	return count > 0, nil
}

func (r *postgresUsersRepository) FindAllWithFilters(params repositories.QueryParams) ([]records.Users, int, error) {
	querySelectUsers := squirrel.
		Select("*").
		From("users").
		PlaceholderFormat(squirrel.Dollar)

	queryCountUsers := squirrel.
		Select("COUNT(*)").
		From("users").
		PlaceholderFormat(squirrel.Dollar)

	querySelectUsers = repositories.ApplyFilters(querySelectUsers, params.Filters)
	queryCountUsers = repositories.ApplyFilters(queryCountUsers, params.Filters)

	querySelectUsers = querySelectUsers.
		OrderBy("id ASC")

	queryCountUsers = queryCountUsers.Where(squirrel.Eq{"deleted_at": nil})

	var total int
	query, args, err := queryCountUsers.ToSql()
	if err != nil {
		return nil, 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindAllWithFilters - squirrel.Select: %w", err))
	}

	if err := r.db.Get(&total, query, args...); err != nil {
		return nil, 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindAllWithFilters - r.db.Get: %w", err))
	}

	if params.Pagination.Limit > 0 {
		querySelectUsers = repositories.ApplyPagination(querySelectUsers, params.Pagination.Page, params.Pagination.Limit)
	}

	query, args, err = querySelectUsers.ToSql()
	if err != nil {
		return nil, 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindAllWithFilters - squirrel.Select: %w", err))
	}

	var users []records.Users
	if err := r.db.Select(&users, query, args...); err != nil {
		return nil, 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresUsersRepository - FindAllWithFilters - r.db.Select: %w", err))
	}

	return users, total, nil
}
