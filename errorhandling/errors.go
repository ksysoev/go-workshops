package errorhandling

import "github.com/jackc/pgx/v5/pgconn"

func GetUsers() error {
	return &pgconn.PgError{
		Severity: "ERROR",
		Code:     "42P01",
		Message:  "relation \"users\" does not exist",
		Detail:   "",
		Hint:     "",
	}
}
