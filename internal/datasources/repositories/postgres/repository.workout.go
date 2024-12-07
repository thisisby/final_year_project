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
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - db.Beginx: %w", err))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	query, args, err := squirrel.
		Select("*").
		From("workouts").
		ToSql()
	if err != nil {
		tx.Rollback()
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - squirrel.Select: %w", err))
	}

	var workouts []records.Workouts
	if err := tx.Select(&workouts, query, args...); err != nil {
		tx.Rollback()
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - tx.Select: %w", err))
	}

	for i, workout := range workouts {
		query, args, err = squirrel.
			Select("*").
			From("workout_exercises").
			PlaceholderFormat(squirrel.Dollar).
			Where(squirrel.Eq{"workout_id": workout.ID}).
			OrderBy("id ASC").
			ToSql()
		if err != nil {
			tx.Rollback()
			return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - squirrel.Select: %w", err))
		}

		var exercises []records.WorkoutExercises
		if err := tx.Select(&exercises, query, args...); err != nil {
			tx.Rollback()
			return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - tx.Select: %w", err))
		}

		workouts[i].Exercises = exercises
	}

	return workouts, nil
}

func (r *postgresWorkoutsRepository) FindByID(id int) (records.Workouts, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return records.Workouts{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - db.Beginx: %w", err))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	query, args, err := squirrel.
		Select("*").
		From("workouts").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		tx.Rollback()
		return records.Workouts{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindByID - squirrel.Select: %w", err))
	}

	var workout records.Workouts
	if err := tx.Get(&workout, query, args...); err != nil {
		tx.Rollback()
		return records.Workouts{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindByID - tx.Get: %w", err))
	}

	query, args, err = squirrel.
		Select("*").
		From("workout_exercises").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"workout_id": workout.ID}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		tx.Rollback()
		return records.Workouts{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindByID - squirrel.Select: %w", err))
	}

	var exercises []records.WorkoutExercises
	if err := tx.Select(&exercises, query, args...); err != nil {
		tx.Rollback()
		return records.Workouts{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindByID - tx.Select: %w", err))
	}

	workout.Exercises = exercises

	return workout, nil
}

func (r *postgresWorkoutsRepository) FindAllByOwnerID(ownerID int) ([]records.Workouts, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - db.Beginx: %w", err))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	query, args, err := squirrel.
		Select("*").
		From("workouts").
		Where(squirrel.Eq{"owner_id": ownerID}).
		PlaceholderFormat(squirrel.Dollar).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		tx.Rollback()
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - squirrel.Select: %w", err))
	}

	var workouts []records.Workouts
	if err := tx.Select(&workouts, query, args...); err != nil {
		tx.Rollback()
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - tx.Select: %w", err))
	}

	for i, workout := range workouts {
		query, args, err = squirrel.
			Select("*").
			From("workout_exercises").
			PlaceholderFormat(squirrel.Dollar).
			Where(squirrel.Eq{"workout_id": workout.ID}).
			OrderBy("id ASC").
			ToSql()
		if err != nil {
			tx.Rollback()
			return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - squirrel.Select: %w", err))
		}

		var exercises []records.WorkoutExercises
		if err := tx.Select(&exercises, query, args...); err != nil {
			tx.Rollback()
			return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - FindAll - tx.Select: %w", err))
		}

		workouts[i].Exercises = exercises
	}

	return workouts, nil
}

func (r *postgresWorkoutsRepository) Save(workout records.Workouts, exerciseNames []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - db.Beginx: %w", err))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var workoutID int
	query, args, err := squirrel.
		Insert("workouts").
		Columns("title", "description", "owner_id").
		Values(workout.Title, workout.Description, workout.OwnerID).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		tx.Rollback()
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - squirrel.Insert: %w", err))
	}

	err = tx.Get(&workoutID, query, args...)
	if err != nil {
		tx.Rollback()
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - tx.Get: %w", err))
	}

	for _, exerciseName := range exerciseNames {
		query, args, err = squirrel.
			Insert("workout_exercises").
			Columns("name", "workout_id").
			Values(exerciseName, workoutID).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			tx.Rollback()
			return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - squirrel.Insert: %w", err))
		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - tx.Exec: %w", err))
		}
	}

	return tx.Commit()
}

func (r *postgresWorkoutsRepository) AddExercise(workoutID int, exerciseNames []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - db.Beginx: %w", err))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	for _, exerciseName := range exerciseNames {
		query, args, err := squirrel.
			Insert("workout_exercises").
			Columns("name", "workout_id").
			Values(exerciseName, workoutID).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			tx.Rollback()
			return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - squirrel.Insert: %w", err))
		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - tx.Exec: %w", err))
		}
	}

	return tx.Commit()
}

func (r *postgresWorkoutsRepository) RemoveExercise(workoutID int, exerciseIDs []int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - db.Beginx: %w", err))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	for _, exerciseID := range exerciseIDs {
		query, args, err := squirrel.
			Delete("workout_exercises").
			Where(squirrel.Eq{"id": exerciseID}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			tx.Rollback()
			return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - squirrel.Delete: %w", err))
		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			return helpers.PostgresErrorTransform(fmt.Errorf("postgresWorkoutsRepository - Save - tx.Exec: %w", err))
		}
	}

	return tx.Commit()
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
