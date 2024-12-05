package records

import (
	"database/sql"
	"time"
)

type Record struct {
	ID        int          `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
