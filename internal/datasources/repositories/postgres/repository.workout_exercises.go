package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/services"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type postgresWorkoutExercisesRepository struct {
	db *sqlx.DB
}

func NewPostgresWorkoutExercisesRepository(db *sqlx.DB) services.WorkoutExercisesRepository {
	return &postgresWorkoutExercisesRepository{
		db: db,
	}
}

func (r *postgresWorkoutExercisesRepository) FindAll() ([]records.WorkoutExercises, error) {
	query, args, err := squirrel.
		Select(`
			workout_exercises.*,
			exercises.id AS "exercise.id",
			exercises.created_at AS "exercise.created_at",
			exercises.updated_at AS "exercise.updated_at",
			exercises.deleted_at AS "exercise.deleted_at",
			exercises.name AS "exercise.name"
		`).
		From("workout_exercises").
		Join("exercises ON workout_exercises.exercise_id = exercises.id").
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - FindAll - squirrel.Select: %w", err))
	}

	var workoutExercises []records.WorkoutExercises
	if err := r.db.Select(&workoutExercises, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - FindAll - db.Select: %w", err))
	}

	return workoutExercises, nil
}

func (r *postgresWorkoutExercisesRepository) FindByID(id int) (records.WorkoutExercises, error) {
	query, args, err := squirrel.
		Select(`
			workout_exercises.*,
			exercises.id AS "exercise.id",
			exercises.created_at AS "exercise.created_at",
			exercises.updated_at AS "exercise.updated_at",
			exercises.deleted_at AS "exercise.deleted_at",
			exercises.name AS "exercise.name"
		`).
		From("workout_exercises").
		Join("exercises ON workout_exercises.exercise_id = exercises.id").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"workout_exercises.id": id}).
		ToSql()
	if err != nil {
		return records.WorkoutExercises{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - FindByID - squirrel.Select: %w", err))
	}

	var workoutExercise records.WorkoutExercises
	if err := r.db.Get(&workoutExercise, query, args...); err != nil {
		return records.WorkoutExercises{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - FindByID - db.Get: %w", err))
	}

	return workoutExercise, nil
}

func (r *postgresWorkoutExercisesRepository) FindAllByWorkoutID(workoutID int) ([]records.WorkoutExercises, error) {
	query, args, err := squirrel.
		Select(`
			workout_exercises.*,
			exercises.id AS "exercise.id",
			exercises.created_at AS "exercise.created_at",
			exercises.updated_at AS "exercise.updated_at",
			exercises.deleted_at AS "exercise.deleted_at",
			exercises.name AS "exercise.name"
		`).
		From("workout_exercises").
		Join("exercises ON workout_exercises.exercise_id = exercises.id").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"workout_id": workoutID}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - FindAllByWorkoutID - squirrel.Select: %w", err))
	}

	var workoutExercises []records.WorkoutExercises
	if err := r.db.Select(&workoutExercises, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - FindAllByWorkoutID - db.Select: %w", err))
	}

	return workoutExercises, nil
}

func (r *postgresWorkoutExercisesRepository) Save(workoutExercise records.WorkoutExercises) (int, error) {
	query, args, err := squirrel.
		Insert("workout_exercises").
		Columns("main_note", "secondary_note", "workout_id", "owner_id", "exercise_id").
		Values(workoutExercise.MainNote, workoutExercise.SecondaryNote, workoutExercise.WorkoutID, workoutExercise.OwnerID, workoutExercise.ExerciseID).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - Save - squirrel.Insert: %w", err))
	}

	var id int
	if err := r.db.Get(&id, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - Save - db.Get: %w", err))
	}

	return id, nil
}

func (r *postgresWorkoutExercisesRepository) Update(id int, workoutExercise map[string]interface{}) error {
	updateQuery := squirrel.Update("workout_exercises").PlaceholderFormat(squirrel.Dollar)

	for key, value := range workoutExercise {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - Update - squirrel.Update: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresWorkoutExercisesRepository) Delete(id int) error {
	query, args, err := squirrel.
		Delete("workout_exercises").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - Delete - squirrel.Delete: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutExercisesRepository - Delete - db.Exec: %w", err))
	}

	return nil
}
