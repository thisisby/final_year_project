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
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

func (r *postgresWorkoutsRepository) RemoveExercise(workoutID int, exerciseIDs []int) error {
	//TODO implement me
	panic("implement me")
}

func (r *postgresWorkoutsRepository) Update(id int, workout map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (r *postgresWorkoutsRepository) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
