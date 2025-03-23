package postgres

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/services"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type postgresNutritionsRepository struct {
	db *sqlx.DB
}

func NewPostgresNutritionsRepository(db *sqlx.DB) services.NutritionsRepository {
	return &postgresNutritionsRepository{db: db}
}

func (p *postgresNutritionsRepository) FindAllByOwnerID(ownerID int) ([]records.Nutritions, error) {
	query, args, err := squirrel.
		Select("*").
		From("nutritions").
		Where("owner_id = ?", ownerID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - FindAllByOwnerID - squirrel.Select: %w", err))
	}

	var nutritions []records.Nutritions
	if err := p.db.Select(&nutritions, query, args...); err != nil {
		return nil, helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - FindAllByOwnerID - db.Select: %w", err))
	}

	return nutritions, nil
}

func (p *postgresNutritionsRepository) FindByID(id int) (records.Nutritions, error) {
	query, args, err := squirrel.
		Select("*").
		From("nutritions").
		Where("id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return records.Nutritions{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - FindByID - squirrel.Select: %w", err))
	}

	var nutrition records.Nutritions
	if err := p.db.Get(&nutrition, query, args...); err != nil {
		return records.Nutritions{}, helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - FindByID - db.Get: %w", err))
	}

	return nutrition, nil
}

func (p *postgresNutritionsRepository) Save(nutrition records.Nutritions) (int, error) {
	query, args, err := squirrel.
		Insert("nutritions").
		Columns("owner_id", "name", "value").
		Values(nutrition.OwnerID, nutrition.Name, nutrition.Value).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - Save - squirrel.Insert: %w", err))
	}

	var id int
	if err := p.db.Get(&id, query, args...); err != nil {
		return 0, helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - Save - db.Get: %w", err))
	}

	return id, nil
}

func (p *postgresNutritionsRepository) Update(id int, nutrition map[string]interface{}) error {
	updateQuery := squirrel.Update("nutritions").PlaceholderFormat(squirrel.Dollar)
	for key, value := range nutrition {
		updateQuery = updateQuery.Set(key, value)
	}

	query, args, err := updateQuery.Where("id = ?", id).ToSql()
	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - Update - squirrel.Update: %w", err))
	}

	if _, err := p.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - Update - db.Exec: %w", err))
	}

	return nil
}

func (p *postgresNutritionsRepository) Delete(id int) error {
	query, args, err := squirrel.
		Delete("nutritions").
		Where("id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - Delete - squirrel.Delete: %w", err))
	}

	if _, err := p.db.Exec(query, args...); err != nil {
		return helpers.PostgresErrorTransform(fmt.Errorf("postgresNutritionsRepository - Delete - db.Exec: %w", err))
	}

	return nil
}
