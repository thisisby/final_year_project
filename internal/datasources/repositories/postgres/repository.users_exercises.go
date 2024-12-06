package postgres

import "database/sql"

type postgresUsersExercisesRepository struct {
	db *sql.DB
}
