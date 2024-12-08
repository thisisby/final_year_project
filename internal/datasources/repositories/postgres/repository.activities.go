package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/services"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type postgresActivitiesRepository struct {
	db *sqlx.DB
}

func NewPostgresActivitiesRepository(db *sqlx.DB) services.ActivitiesRepository {
	return &postgresActivitiesRepository{
		db: db,
	}
}

func (r *postgresActivitiesRepository) FindAll() ([]records.Activities, error) {
	query, args, err := squirrel.
		Select("*").
		From("activities").
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - FindAll - squirrel.Select: %w", err))
	}

	var activities []records.Activities
	if err := r.db.Select(&activities, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - FindAll - db.Select: %w", err))
	}

	return activities, nil
}

func (r *postgresActivitiesRepository) FindByID(id int) (records.Activities, error) {
	query, args, err := squirrel.
		Select("*").
		From("activities").
		Where("id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return records.Activities{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - FindByID - squirrel.Select: %w", err))
	}

	var activity records.Activities
	if err := r.db.Get(&activity, query, args...); err != nil {
		return records.Activities{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - FindByID - db.Get: %w", err))
	}

	return activity, nil
}

func (r *postgresActivitiesRepository) Save(activity records.Activities) (int, error) {
	query, args, err := squirrel.
		Insert("activities").
		Columns("name", "activity_group_id").
		Values(activity.Name, activity.ActivityGroupID).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - Save - squirrel.Insert: %w", err))
	}

	var id int
	if err := r.db.Get(&id, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - Save - db.Get: %w", err))
	}

	return id, nil
}

func (r *postgresActivitiesRepository) Update(id int, activityMap map[string]interface{}) error {
	updateQuery := squirrel.Update("activities").PlaceholderFormat(squirrel.Dollar)
	for key, value := range activityMap {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - Update - squirrel.Update: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresActivitiesRepository) Delete(id int) error {
	query, args, err := squirrel.
		Delete("activities").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - Delete - squirrel.Delete: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresActivitiesRepository - Delete - db.Exec: %w", err))
	}

	return nil
}
