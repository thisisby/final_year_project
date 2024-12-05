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

type postgresTokensRepository struct {
	db *sqlx.DB
}

func NewPostgresTokensRepository(db *sqlx.DB) services.TokensRepository {
	return &postgresTokensRepository{db}
}

func (r *postgresTokensRepository) Save(token records.Tokens) error {
	query, args, err := squirrel.Insert("tokens").
		Columns("token", "user_id", "created_at", "updated_at", "deleted_at").
		Values(token.Token, token.UserID, time.Now(), nil, nil).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresTokensRepository - Save - squirrel.Insert: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresTokensRepository - Save - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresTokensRepository) FindByToken(token string) (records.Tokens, error) {
	query, args, err := squirrel.Select("*").
		From("tokens").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"token": token}).ToSql()
	if err != nil {
		return records.Tokens{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresTokensRepository - FindByToken - squirrel.Select: %w", err))
	}

	var tokenRecord records.Tokens
	if err := r.db.Get(&tokenRecord, query, args...); err != nil {
		return records.Tokens{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresTokensRepository - FindByToken - db.Get: %w", err))
	}

	return tokenRecord, nil
}

func (r *postgresTokensRepository) Delete(token string) error {
	query, args, err := squirrel.Update("tokens").
		Set("deleted_at", time.Now()).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"token": token}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresTokensRepository - Delete - squirrel.Update: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresTokensRepository - Delete - db.Exec: %w", err))
	}

	return nil
}
