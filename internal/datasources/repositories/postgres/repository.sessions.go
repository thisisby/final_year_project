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

type postgresSessionsRepository struct {
	db *sqlx.DB
}

func NewPostgresSessionsRepository(db *sqlx.DB) services.SessionsRepository {
	return &postgresSessionsRepository{db: db}
}

func (r *postgresSessionsRepository) FindAll() ([]records.Sessions, error) {
	query, args, err := squirrel.
		Select(`
			sessions.*,
			activities.id AS "activity.id",
			activities.created_at AS "activity.created_at",
			activities.updated_at AS "activity.updated_at",
			activities.deleted_at AS "activity.deleted_at",
			activities.name AS "activity.name",
			activities.activity_group_id AS "activity.activity_group_id"
		`).
		From("sessions").
		LeftJoin("activities ON sessions.activity_id = activities.id").
		PlaceholderFormat(squirrel.Dollar).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - FindAll - squirrel: %w", err))
	}

	var sessions []records.Sessions
	if err := r.db.Select(&sessions, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository -FindAll - db.Select: %w", err))
	}

	return sessions, nil
}

func (r *postgresSessionsRepository) FindByID(id int) (records.Sessions, error) {
	query, args, err := squirrel.
		Select(`
			sessions.*,
			activities.id AS "activity.id",
			activities.created_at AS "activity.created_at",
			activities.updated_at AS "activity.updated_at",
			activities.deleted_at AS "activity.deleted_at",
			activities.name AS "activity.name",
			activities.activity_group_id AS "activity.activity_group_id"
		`).
		From("sessions").
		LeftJoin("activities ON sessions.activity_id = activities.id").
		Where("sessions.id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return records.Sessions{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - FindByID - squirrel: %w", err))
	}

	var session records.Sessions
	if err := r.db.Get(&session, query, args...); err != nil {
		return records.Sessions{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository -FindByID - db.Get: %w", err))
	}

	return session, nil
}

func (r *postgresSessionsRepository) Save(session records.Sessions) (int, error) {
	query, args, err := squirrel.
		Insert("sessions").
		Columns("notes", "start_time", "end_time", "activity_id", "owner_id").
		Values(session.Notes, session.StartTime, session.EndTime, session.ActivityID, session.OwnerID).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - Save - squirrel: %w", err))
	}

	var id int
	if err := r.db.Get(&id, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - Save - db.Get: %w", err))
	}

	return id, nil
}

func (r *postgresSessionsRepository) Update(id int, sessionMap map[string]interface{}) error {
	updateQuery := squirrel.Update("sessions").PlaceholderFormat(squirrel.Dollar)
	for key, value := range sessionMap {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - Update - squirrel: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresSessionsRepository) Delete(id int) error {
	query, args, err := squirrel.
		Delete("sessions").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - Delete - squirrel: %w", err))
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - Delete - db.Exec: %w", err))
	}

	return nil
}

func (r *postgresSessionsRepository) FindAllByStartTime(ownerID int, createdAt time.Time) ([]records.Sessions, error) {
	query, args, err := squirrel.
		Select(`
			sessions.*,
			activities.id AS "activity.id",
			activities.created_at AS "activity.created_at",
			activities.updated_at AS "activity.updated_at",
			activities.deleted_at AS "activity.deleted_at",
			activities.name AS "activity.name",
			activities.activity_group_id AS "activity.activity_group_id"
		`).
		From("sessions").
		LeftJoin("activities ON sessions.activity_id = activities.id").
		Where(squirrel.Eq{"sessions.owner_id": ownerID}).
		Where(squirrel.Expr("DATE(sessions.start_time) = ?", createdAt.Format("2006-01-02"))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return []records.Sessions{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - FindAllByStartTime - squirrel: %w", err))
	}

	var sessions []records.Sessions
	if err := r.db.Select(&sessions, query, args...); err != nil {
		return []records.Sessions{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - FindAllByStartTime - db.Select: %w", err))
	}

	fmt.Printf("sessions: %v\n", sessions)
	return sessions, nil
}

func (r *postgresSessionsRepository) FindAllInDateRange(ownerID int, startDate time.Time, endDate time.Time) ([]records.Sessions, error) {
	query, args, err := squirrel.
		Select(`
			sessions.*,
			activities.id AS "activity.id",
			activities.created_at AS "activity.created_at",
			activities.updated_at AS "activity.updated_at",
			activities.deleted_at AS "activity.deleted_at",
			activities.name AS "activity.name",
			activities.activity_group_id AS "activity.activity_group_id"
		`).
		From("sessions").
		LeftJoin("activities ON sessions.activity_id = activities.id").
		Where(squirrel.Eq{"sessions.owner_id": ownerID}).
		Where(squirrel.Expr("DATE(sessions.start_time) >= ?", startDate.Format("2006-01-02"))).
		Where(squirrel.Expr("DATE(sessions.start_time) <= ?", endDate.Format("2006-01-02"))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return []records.Sessions{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - FindAllInDateRange - squirrel: %w", err))
	}

	var sessions []records.Sessions
	if err := r.db.Select(&sessions, query, args...); err != nil {
		return []records.Sessions{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresSessionsRepository - FindAllInDateRange - db.Select: %w", err))
	}

	return sessions, nil
}
