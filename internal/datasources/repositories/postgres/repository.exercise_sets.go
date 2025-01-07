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

func (r *postgresExerciseSetsRepository) FindTotalSetsByDate(ownerID int, date time.Time) (int, error) {
	query, args, err := squirrel.
		Select("COUNT(*)").
		From("exercise_sets").
		Join("workout_exercises ON exercise_sets.workout_exercise_id = workout_exercises.id").
		Join("workouts ON workout_exercises.workout_id = workouts.id").
		Where(squirrel.Eq{"exercise_sets.owner_id": ownerID}).
		Where(squirrel.Expr("DATE(exercise_sets.created_at) = ?", date.Format("2006-01-02"))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindTotalSetsByDate - squirrel.Select: %w", err))
	}

	var totalSets int
	if err := r.db.Get(&totalSets, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindTotalSetsByDate - db.Get: %w", err))
	}

	return totalSets, nil
}

func (r *postgresExerciseSetsRepository) FindTotalRepsByDate(ownerID int, date time.Time) (int, error) {
	query, args, err := squirrel.
		Select("SUM(reps)").
		From("exercise_sets").
		Join("workout_exercises ON exercise_sets.workout_exercise_id = workout_exercises.id").
		Join("workouts ON workout_exercises.workout_id = workouts.id").
		Where(squirrel.Eq{"exercise_sets.owner_id": ownerID}).
		Where(squirrel.Expr("DATE(exercise_sets.created_at) = ?", date.Format("2006-01-02"))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindTotalRepsByDate - squirrel.Select: %w", err))
	}

	var totalReps int
	if err := r.db.Get(&totalReps, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindTotalRepsByDate - db.Get: %w", err))
	}

	return totalReps, nil
}

func (r *postgresExerciseSetsRepository) FindUniqueWorkoutExercisesByDate(ownerID int, date time.Time) (int, error) {
	query, args, err := squirrel.
		Select("COUNT(DISTINCT workout_exercises.id)").
		From("exercise_sets").
		Join("workout_exercises ON exercise_sets.workout_exercise_id = workout_exercises.id").
		Join("workouts ON workout_exercises.workout_id = workouts.id").
		Where(squirrel.Eq{"exercise_sets.owner_id": ownerID}).
		Where(squirrel.Expr("DATE(exercise_sets.created_at) = ?", date.Format("2006-01-02"))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindUniqueWorkoutsByDate - squirrel.Select: %w", err))
	}

	var uniqueWorkouts int
	if err := r.db.Get(&uniqueWorkouts, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindUniqueWorkoutsByDate - db.Get: %w", err))
	}

	return uniqueWorkouts, nil
}

func (r *postgresExerciseSetsRepository) FindExercisesDetailsByDate(ownerID int, date time.Time) ([]records.ExerciseDetails, error) {
	query, args, err := squirrel.
		Select(
			"workout_exercises.exercise_name",
			"exercise_sets.reps",
			"exercise_sets.weight",
		).
		From("exercise_sets").
		Join("workout_exercises ON exercise_sets.workout_exercise_id = workout_exercises.id").
		Join("workouts ON workout_exercises.workout_id = workouts.id").
		Where(squirrel.Eq{"exercise_sets.owner_id": ownerID}).
		Where(squirrel.Expr("DATE(exercise_sets.created_at) = ?", date.Format("2006-01-02"))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindExercisesDetailsByDate - squirrel.Select: %w", err))
	}

	var exercises []records.ExerciseDetails
	if err := r.db.Select(&exercises, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindExercisesDetailsByDate - db.Select: %w", err))
	}

	return exercises, nil
}

func (r *postgresExerciseSetsRepository) FindAllByCreatedAt(ownerID int, createdAt time.Time) ([]records.ExerciseSets, error) {
	query, args, err := squirrel.
		Select(`
			exercise_sets.*,
			workout_exercises.id AS "workout_exercise.id",
			workout_exercises.created_at AS "workout_exercise.created_at",
			workout_exercises.updated_at AS "workout_exercise.updated_at",
		    workout_exercises.deleted_at AS "workout_exercise.deleted_at",
			workout_exercises.main_note AS "workout_exercise.main_note",
			workout_exercises.secondary_note AS "workout_exercise.secondary_note",
			workout_exercises.workout_id AS "workout_exercise.workout_id",
			workout_exercises.owner_id AS "workout_exercise.owner_id",
			workout_exercises.exercise_id AS "workout_exercise.exercise_id",
			exercises.id AS "workout_exercise.exercise.id",
			exercises.created_at AS "workout_exercise.exercise.created_at",
			exercises.updated_at AS "workout_exercise.exercise.updated_at",
			exercises.deleted_at AS "workout_exercise.exercise.deleted_at",
			exercises.name AS "workout_exercise.exercise.name"	
		`).
		From("exercise_sets").
		Join("workout_exercises ON exercise_sets.workout_exercise_id = workout_exercises.id").
		Join("exercises ON workout_exercises.exercise_id = exercises.id").
		Where(squirrel.Eq{"exercise_sets.owner_id": ownerID}).
		Where(squirrel.Expr("DATE(exercise_sets.created_at) = ?", createdAt.Format("2006-01-02"))).
		OrderBy("exercise_sets.id ASC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindAllByStartTime - squirrel.Select: %w", err))
	}

	var exerciseSets []records.ExerciseSets
	if err := r.db.Select(&exerciseSets, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindAllByStartTime - db.Select: %w", err))
	}

	return exerciseSets, nil
}

func (r *postgresExerciseSetsRepository) FindAllInDateRange(ownerID int, startDate time.Time, endDate time.Time) ([]records.ExerciseSets, error) {
	query, args, err := squirrel.
		Select(`
			exercise_sets.*,
			workout_exercises.id AS "workout_exercise.id",
			workout_exercises.created_at AS "workout_exercise.created_at",
			workout_exercises.updated_at AS "workout_exercise.updated_at",
		    workout_exercises.deleted_at AS "workout_exercise.deleted_at",
			workout_exercises.main_note AS "workout_exercise.main_note",
			workout_exercises.secondary_note AS "workout_exercise.secondary_note",
			workout_exercises.workout_id AS "workout_exercise.workout_id",
			workout_exercises.owner_id AS "workout_exercise.owner_id",
			workout_exercises.exercise_id AS "workout_exercise.exercise_id",
			exercises.id AS "workout_exercise.exercise.id",
			exercises.created_at AS "workout_exercise.exercise.created_at",
			exercises.updated_at AS "workout_exercise.exercise.updated_at",
			exercises.deleted_at AS "workout_exercise.exercise.deleted_at",
			exercises.name AS "workout_exercise.exercise.name"	
		`).
		From("exercise_sets").
		Join("workout_exercises ON exercise_sets.workout_exercise_id = workout_exercises.id").
		Join("exercises ON workout_exercises.exercise_id = exercises.id").
		Where(squirrel.Eq{"exercise_sets.owner_id": ownerID}).
		Where(squirrel.Expr("DATE(exercise_sets.created_at) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))).
		OrderBy("exercise_sets.id ASC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindAllInDateRange - squirrel.Select: %w", err))
	}

	var exerciseSets []records.ExerciseSets
	if err := r.db.Select(&exerciseSets, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresExerciseSetsRepository - FindAllInDateRange - db.Select: %w", err))
	}

	return exerciseSets, nil
}
