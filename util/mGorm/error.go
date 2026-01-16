package mGorm

import (
	"errors"

	pq "github.com/jackc/pgconn"
)

var (
	// ErrorDuplicateValues duplicate value
	ErrorDuplicateValues = errors.New("Duplicate Values")
	// ErrorForeignKeyConstraint foreign key constraint
	ErrorForeignKeyConstraint = errors.New("Foreignkey Constraint")
	// ErrorCheckConstraint check constraint
	ErrorCheckConstraint = errors.New("Check Constraint")
)

var (
	// https://www.postgresql.org/docs/10/errcodes-appendix.html
	codeToError = map[string]error{
		"23505": ErrorDuplicateValues,
		"23503": ErrorForeignKeyConstraint,
		"23514": ErrorCheckConstraint,
	}
)

// ParseDBError parse db error
func ParseDBError(err error) error {
	if err, ok := err.(*pq.PgError); ok {
		if cErr, ok := codeToError[err.Code]; ok {
			return cErr
		}
	}
	return err
}
