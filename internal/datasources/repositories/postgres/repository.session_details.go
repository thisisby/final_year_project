package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/services"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type postgresSessionDetailsRepository struct {
	db *sqlx.DB
}

func NewPostgresSessionDetailsRepository(db *sqlx.DB) services.SessionDetailsRepository {
	return &postgresSessionDetailsRepository{db: db}
}

func (r *postgresSessionDetailsRepository) FindAll() ([]records.SessionDetails, error) {
	query, args, err := squirrel.
		Select("*").
		From("session_details").
		Where("deleted_at IS NULL").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - FindAll - squirrel.Select: %w", err))
	}

	var details []records.SessionDetails
	if err := r.db.Select(&details, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - FindAll - db.Select: %w", err))
	}
	return details, nil
}

func (r *postgresSessionDetailsRepository) Save(sessionDetails records.SessionDetails) (int, error) {
	query, args, err := squirrel.
		Insert("session_details").
		Columns("session_id", "name", "value").
		Values(sessionDetails.SessionID, sessionDetails.Name, sessionDetails.Value).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - Create - squirrel.Insert: %w", err))
	}

	var id int
	if err := r.db.Get(&id, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - Save - db.Get: %w", err))
	}
	return id, nil
}

func (r *postgresSessionDetailsRepository) FindByID(id int) (records.SessionDetails, error) {
	query, args, err := squirrel.
		Select("*").
		From("session_details").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return records.SessionDetails{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - FindByID - squirrel.Select: %w", err))
	}

	var detail records.SessionDetails
	if err := r.db.Get(&detail, query, args...); err != nil {
		return records.SessionDetails{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - FindByID - db.Get: %w", err))
	}
	return detail, nil
}

func (r *postgresSessionDetailsRepository) Update(id int, sessionDetailUpdates map[string]interface{}) error {
	updateQuery := squirrel.Update("session_details").PlaceholderFormat(squirrel.Dollar)

	for key, value := range sessionDetailUpdates {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - Update - squirrel.Update: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresSessionDetailsRepository) Delete(id int) error {
	query, args, err := squirrel.
		Delete("session_details").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - Delete - squirrel.Delete: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - Delete - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresSessionDetailsRepository) FindAllBySessionID(sessionID int) ([]records.SessionDetails, error) {
	query, args, err := squirrel.
		Select("*").
		From("session_details").
		Where(squirrel.Eq{"session_id": sessionID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - FindAllBySessionID - squirrel.Select: %w", err))
	}

	var details []records.SessionDetails
	if err := r.db.Select(&details, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionDetailsRepository - FindAllBySessionID - db.Select: %w", err))
	}
	return details, nil
}
