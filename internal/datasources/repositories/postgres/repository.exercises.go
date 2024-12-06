package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/services"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type postgresExercisesRepository struct {
	db *sqlx.DB
}

func NewPostgresExercisesRepository(db *sqlx.DB) services.ExercisesRepository {
	return &postgresExercisesRepository{db}
}

func (r *postgresExercisesRepository) FindAll() ([]records.Exercises, error) {
	query, args, err := squirrel.Select("*").From("exercises").ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindAll - squirrel.Select: %w", err))
	}

	var exercises []records.Exercises
	if err := r.db.Select(&exercises, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindAll - db.Select: %w", err))
	}

	return exercises, nil
}

func (r *postgresExercisesRepository) FindByID(id int) (records.Exercises, error) {
	query, args, err := squirrel.Select("*").
		From("exercises").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return records.Exercises{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindByID - squirrel.Select: %w", err))
	}

	var exercise records.Exercises
	if err := r.db.Get(&exercise, query, args...); err != nil {
		return records.Exercises{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindByID - db.Get: %w", err))
	}

	return exercise, nil
}

func (r *postgresExercisesRepository) FindByName(name string) (records.Exercises, error) {
	query, args, err := squirrel.Select("*").
		From("exercises").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"name": name}).ToSql()
	if err != nil {
		return records.Exercises{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindByName - squirrel.Select: %w", err))
	}

	var exercise records.Exercises
	if err := r.db.Get(&exercise, query, args...); err != nil {
		return records.Exercises{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindByName - db.Get: %w", err))
	}

	return exercise, nil
}

func (r *postgresExercisesRepository) Save(exercise records.Exercises) error {
	query, args, err := squirrel.Insert("exercises").
		Columns("name").
		Values(exercise.Name).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - Save - squirrel.Insert: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - Save - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresExercisesRepository) Update(id int, exercise map[string]interface{}) error {
	updateQuery := squirrel.Update("exercises").PlaceholderFormat(squirrel.Dollar)

	for key, value := range exercise {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - Update - squirrel.Update: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresExercisesRepository) Delete(id int) error {
	query, args, err := squirrel.Delete("exercises").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - Delete - squirrel.Delete: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - Delete - db.Exec: %w", err))
	}

	return nil
}
