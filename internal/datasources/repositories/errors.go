package repositories

import "errors"

var (
	ErrorRowExists           = errors.New("row already exists")
	ErrorRowNotFound         = errors.New("row not found")
	ErrorForeignKeyViolation = errors.New("foreign key violation")
)
