package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/services"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type postgresActivityGroupsRepository struct {
	db *sqlx.DB
}

func NewPostgresActivityGroupsRepository(db *sqlx.DB) services.ActivityGroupsRepository {
	return &postgresActivityGroupsRepository{
		db: db,
	}
}

func (r *postgresActivityGroupsRepository) FindAll() ([]records.ActivityGroups, error) {
	query, args, err := squirrel.
		Select("*").
		From("activity_groups").
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - FindAll - squirrel.Select: %w", err))
	}

	var activityGroups []records.ActivityGroups
	if err := r.db.Select(&activityGroups, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - FindAll - db.Select: %w", err))
	}

	return activityGroups, nil
}

func (r *postgresActivityGroupsRepository) FindByID(id int) (records.ActivityGroups, error) {
	query, args, err := squirrel.
		Select("*").
		From("activity_groups").
		Where("id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return records.ActivityGroups{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - FindByID - squirrel.Select: %w", err))
	}

	var activityGroup records.ActivityGroups
	if err := r.db.Get(&activityGroup, query, args...); err != nil {
		return records.ActivityGroups{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - FindByID - db.Get: %w", err))
	}

	return activityGroup, nil
}

func (r *postgresActivityGroupsRepository) Save(activityGroup records.ActivityGroups) (int, error) {
	query, args, err := squirrel.
		Insert("activity_groups").
		Columns("name", "description").
		Values(activityGroup.Name, activityGroup.Description).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - Save - squirrel.Insert: %w", err))
	}

	var id int
	if err := r.db.Get(&id, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - Save - db.Get: %w", err))
	}

	return id, nil
}

func (r *postgresActivityGroupsRepository) Update(id int, activityGroupMap map[string]interface{}) error {
	updateQuery := squirrel.Update("activity_groups").PlaceholderFormat(squirrel.Dollar)
	for key, value := range activityGroupMap {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - Update - squirrel.Update: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresActivityGroupsRepository) Delete(id int) error {
	query, args, err := squirrel.
		Delete("activity_groups").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - Delete - squirrel.Delete: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresActivityGroupsRepository - Delete - db.Exec: %w", err))
	}

	return nil
}
