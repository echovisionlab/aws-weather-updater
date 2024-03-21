package mock

import (
	"context"
	"database/sql"
)

type Database struct {
}

func (d Database) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (d Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (d Database) Rebind(query string) string {
	//TODO implement me
	panic("implement me")
}

func (d Database) Close() error {
	//TODO implement me
	panic("implement me")
}

func (d Database) Query(query string, args ...any) (*sql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (d Database) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	//TODO implement me
	panic("implement me")
}
