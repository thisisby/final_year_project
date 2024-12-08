package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/services"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type postgresExerciseSetsRepository struct {
	db *sqlx.DB
}

func NewPostgresExerciseSetsRepository(db *sqlx.DB) services.ExerciseSetsRepository {
	return &postgresExerciseSetsRepository{
		db: db,
	}
}

func (r *postgresExerciseSetsRepository) Save(exerciseSet records.ExerciseSets) (int, error) {
	query, args, err := squirrel.
		Insert("exercise_sets").
		Columns("reps", "weight", "workout_exercise_id", "owner_id").
		Values(exerciseSet.Reps, exerciseSet.Weight, exerciseSet.WorkoutExerciseID, exerciseSet.OwnerID).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - Save - squirrel.Insert: %w", err))
	}

	var id int
	if err := r.db.Get(&id, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - Save - db.Get: %w", err))
	}

	return id, nil
}

func (r *postgresExerciseSetsRepository) FindAllByWorkoutExerciseID(workoutExerciseID int) ([]records.ExerciseSets, error) {
	query, args, err := squirrel.
		Select("*").
		From("exercise_sets").
		Where(squirrel.Eq{"workout_exercise_id": workoutExerciseID}).
		OrderBy("id ASC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindAllByWorkoutExerciseID - squirrel.Select: %w", err))
	}

	var exerciseSets []records.ExerciseSets
	if err := r.db.Select(&exerciseSets, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindAllByWorkoutExerciseID - db.Select: %w", err))
	}

	return exerciseSets, nil
}

func (r *postgresExerciseSetsRepository) FindByID(id int) (records.ExerciseSets, error) {
	query, args, err := squirrel.
		Select("id", "reps", "weight", "workout_exercise_id", "owner_id").
		From("exercise_sets").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return records.ExerciseSets{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindByID - squirrel.Select: %w", err))
	}

	var exerciseSet records.ExerciseSets
	if err := r.db.Get(&exerciseSet, query, args...); err != nil {
		return records.ExerciseSets{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindByID - db.Get: %w", err))
	}

	return exerciseSet, nil
}

func (r *postgresExerciseSetsRepository) Update(id int, exerciseSetMap map[string]interface{}) error {
	updateQuery := squirrel.Update("exercise_sets").PlaceholderFormat(squirrel.Dollar)
	for key, value := range exerciseSetMap {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - Update - squirrel.Update: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresExerciseSetsRepository) Delete(id int) error {
	query, args, err := squirrel.
		Delete("exercise_sets").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - Delete - squirrel.Delete: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - Delete - db.Exec: %w", err))
	}

	return nil
}
