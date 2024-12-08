package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/services"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type postgresWorkoutsRepository struct {
	db *sqlx.DB
}

func NewPostgresWorkoutsRepository(db *sqlx.DB) services.WorkoutsRepository {
	return &postgresWorkoutsRepository{db}
}

func (r *postgresWorkoutsRepository) FindAll() ([]records.Workouts, error) {
	query, args, err := squirrel.
		Select("*").
		From("workouts").
		Where(squirrel.Eq{"is_private": false}).
		PlaceholderFormat(squirrel.Dollar).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - squirrel.Select: %w", err))
	}

	var workouts []records.Workouts
	if err := r.db.Select(&workouts, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - tx.Select: %w", err))
	}

	return workouts, nil
}

func (r *postgresWorkoutsRepository) FindByID(id int) (records.Workouts, error) {
	query, args, err := squirrel.
		Select("*").
		From("workouts").
		Where(squirrel.Eq{"id": id}, squirrel.Eq{"is_private": false}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return records.Workouts{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindByID - squirrel.Select: %w", err))
	}

	var workout records.Workouts
	if err := r.db.Get(&workout, query, args...); err != nil {
		return records.Workouts{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindByID - r.db.Get: %w", err))
	}

	return workout, nil
}

func (r *postgresWorkoutsRepository) FindAllByOwnerID(ownerID int) ([]records.Workouts, error) {
	query, args, err := squirrel.
		Select("*").
		From("workouts").
		Where(squirrel.Eq{"owner_id": ownerID}, squirrel.Eq{"is_private": false}).
		PlaceholderFormat(squirrel.Dollar).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - squirrel.Select: %w", err))
	}

	var workouts []records.Workouts
	if err := r.db.Select(&workouts, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - r.db.Select: %w", err))
	}

	return workouts, nil
}

func (r *postgresWorkoutsRepository) FindAllByCurrentUserID(ownerID int) ([]records.Workouts, error) {
	query, args, err := squirrel.
		Select("*").
		From("workouts").
		Where(squirrel.Eq{"owner_id": ownerID}).
		PlaceholderFormat(squirrel.Dollar).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAllByCurrentUserID - squirrel.Select: %w", err))
	}

	var workouts []records.Workouts
	if err := r.db.Select(&workouts, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAllByCurrentUserID - r.db.Select: %w", err))
	}

	return workouts, nil
}

func (r *postgresWorkoutsRepository) Save(workout records.Workouts) (int, error) {
	var workoutID int
	query, args, err := squirrel.
		Insert("workouts").
		Columns("title", "description", "owner_id", "is_private", "price").
		Values(workout.Title, workout.Description, workout.OwnerID, workout.IsPrivate, workout.Price).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - squirrel.Insert: %w", err))
	}

	err = r.db.Get(&workoutID, query, args...)
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - db.Get: %w", err))
	}

	return workoutID, nil
}

func (r *postgresWorkoutsRepository) Update(id int, workout map[string]interface{}) error {
	updateQuery := squirrel.
		Update("workouts").
		PlaceholderFormat(squirrel.Dollar)

	for key, value := range workout {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Update - squirrel.Update: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresWorkoutsRepository) Delete(id int) error {
	query, args, err := squirrel.
		Delete("workouts").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Delete - squirrel.Delete: %w", err))
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Delete - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresWorkoutsRepository) Copy(id int, userID int) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - r.db.Beginx: %w", err))
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	workoutQuery, args, err := squirrel.
		Select("*").
		From("workouts").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - squirrel.Select: %w", err))
	}

	var workout records.Workouts
	if err := tx.Get(&workout, workoutQuery, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - tx.Get: %w", err))
	}

	workout.OwnerID = userID

	queryInsert, args, err := squirrel.
		Insert("workouts").
		Columns("title", "description", "owner_id").
		Values(workout.Title, workout.Description, workout.OwnerID).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - squirrel.Insert: %w", err))
	}

	var workoutID int
	if err := tx.Get(&workoutID, queryInsert, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - tx.Get: %w", err))
	}

	workoutExercisesQuery, args, err := squirrel.
		Select("*").
		From("workout_exercises").
		Where(squirrel.Eq{"workout_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - squirrel.Select: %w", err))
	}

	var workoutExercises []records.WorkoutExercises
	if err := tx.Select(&workoutExercises, workoutExercisesQuery, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - tx.Select: %w", err))
	}

	for _, workoutExercise := range workoutExercises {
		workoutExercise.WorkoutID = workoutID

		queryInsert, args, err := squirrel.
			Insert("workout_exercises").
			Columns("workout_id", "exercise_id", "main_note", "secondary_note", "owner_id").
			Values(workoutExercise.WorkoutID, workoutExercise.ExerciseID, workoutExercise.MainNote, workoutExercise.SecondaryNote, workoutExercise.OwnerID).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - squirrel.Insert: %w", err))
		}

		if _, err := tx.Exec(queryInsert, args...); err != nil {
			return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - tx.Exec: %w", err))
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Copy - tx.Commit: %w", err))
	}

	return workoutID, nil
}
