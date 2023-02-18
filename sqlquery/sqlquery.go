package sqlquery

import (
	"context"
	"database/sql"

	"github.com/anton7r/orava/dbquery"
)

// Querier is something that sqlscan can query and get the *sql.Rows from.
// For example, it can be: *sql.DB, *sql.Conn or *sql.Tx.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

var (
	_ Querier = &sql.DB{}
	_ Querier = &sql.Conn{}
	_ Querier = &sql.Tx{}
)

type API struct {
	dbqueryAPI *dbquery.API
}

// NewAPI creates new API instance from dbquery.API instance.
func NewAPI(dbqueryAPI *dbquery.API) (*API, error) {
	api := &API{
		dbqueryAPI: dbqueryAPI,
	}
	return api, nil
}
