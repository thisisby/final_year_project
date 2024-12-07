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
	query, args, err := squirrel.
		Select("e.*").
		From("exercises e").
		LeftJoin("user_exercises ue ON e.id = ue.exercise_id").
		Where(squirrel.Expr("ue.exercise_id IS NULL")).
		OrderBy("id ASC").
		ToSql()
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

func (r *postgresExercisesRepository) CreateCustomExercise(exercise records.Exercises, userID int) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - CreateCustomExercise - db.Beginx: %w", err))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	query, args, err := squirrel.
		Insert("exercises").
		Columns("name").
		Values(exercise.Name).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - CreateCustomExercise - squirrel.Insert: %w", err))
	}

	var exerciseID int
	if err := tx.Get(&exerciseID, query, args...); err != nil {
		tx.Rollback()
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - CreateCustomExercise - tx.Get: %w", err))
	}

	query, args, err = squirrel.
		Insert("user_exercises").
		Columns("user_id", "exercise_id").
		Values(userID, exerciseID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		tx.Rollback()
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - CreateCustomExercise - squirrel.Insert: %w", err))
	}

	if _, err := tx.Exec(query, args...); err != nil {
		tx.Rollback()
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - CreateCustomExercise - tx.Exec: %w", err))
	}

	if err := tx.Commit(); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - CreateCustomExercise - tx.Commit: %w", err))
	}

	return exerciseID, nil
}

func (r *postgresExercisesRepository) FindAllUserExercises(userID int) ([]records.Exercises, error) {
	query, args, err := squirrel.
		Select("e.*").
		From("exercises e").
		Join("user_exercises ue ON e.id = ue.exercise_id").
		Where(squirrel.Eq{"ue.user_id": userID}).
		OrderBy("e.id ASC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindAllUserExercises - squirrel.Select: %w", err))
	}

	var exercises []records.Exercises
	if err := r.db.Select(&exercises, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindAllUserExercises - db.Select: %w", err))
	}

	return exercises, nil
}

func (r *postgresExercisesRepository) FindAllWithWorkoutCheck(workoutID int) ([]records.ExercisesWithWorkoutCheck, error) {
	query, args, err := squirrel.
		Select(`
			e.*,
			CASE 
				WHEN we.workout_id IS NOT NULL THEN true
				ELSE false
			END AS "is_in_workout"
		`).
		From("exercises e").
		LeftJoin("workout_exercises we ON e.id = we.exercise_id AND we.workout_id = ?", workoutID).
		OrderBy("e.id ASC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindAllWithWorkoutCheck - squirrel.Select: %w", err))
	}

	var exercises []records.ExercisesWithWorkoutCheck
	if err := r.db.Select(&exercises, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExercisesRepository - FindAllWithWorkoutCheck - db.Select: %w", err))
	}

	return exercises, nil
}
