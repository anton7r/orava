package pgxquery

import (
	"context"

	"github.com/anton7r/orava/dbquery"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// Querier is something that pgxscan can query and get the pgx.Rows from.
// For example, it can be: *pgxpool.Pool, *pgx.Conn or pgx.Tx.
type Querier interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

var (
	_ Querier = &pgxpool.Pool{}
	_ Querier = &pgx.Conn{}
	_ Querier = pgx.Tx(nil)
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

// Select is a high-level function that queries rows from Querier and calls the ScanAll function.
// See ScanAll for details.
func (api *API) Select(ctx context.Context, db Querier, dst interface{}, query string, args ...interface{}) error {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "orava: query multiple result rows")
	}
	err = api.ScanAll(dst, rows)
	return errors.WithStack(err)
}

// SelectNamed is a high-level function that queries rows from Querier and calls the ScanAll function.
// See ScanAll for details.
func (api *API) SelectNamed(ctx context.Context, db Querier, dst interface{}, query string, arg interface{}) error {
	compiledQuery, args, err := api.dbqueryAPI.NamedQueryParams(query, arg)
	if err != nil {
		return err
	}

	return api.Select(ctx, db, dst, compiledQuery, args)
}

// Get is a high-level function that queries rows from Querier and calls the ScanOne function.
// See ScanOne for details.
func (api *API) Get(ctx context.Context, db Querier, dst interface{}, query string, args ...interface{}) error {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "orava: query one result row")
	}
	err = api.ScanOne(dst, rows)
	return errors.WithStack(err)
}

// GetNamed is a high-level function that queries rows from Querier and calls the ScanOne function.
// See ScanOne for details.
func (api *API) GetNamed(ctx context.Context, db Querier, dst interface{}, query string, arg interface{}) error {
	compiledQuery, args, err := api.dbqueryAPI.NamedQueryParams(query, arg)
	if err != nil {
		return err
	}

	return api.Get(ctx, db, dst, compiledQuery, args)
}

// Exec is a high-level function that sends an executable action to the database
func (api *API) Exec(ctx context.Context, db Querier, query string, args ...interface{}) (pgconn.CommandTag, error) {
	tag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return pgconn.CommandTag{}, errors.Wrap(err, "orava: exec")
	}

	return tag, nil
}

// ExecNamed is a high-level function that sends an executable action to the database with named parameters
func (api *API) ExecNamed(ctx context.Context, db Querier, query string, arg interface{}) (pgconn.CommandTag, error) {
	compiledQuery, args, err := api.dbqueryAPI.NamedQueryParams(query, arg)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return api.Exec(ctx, db, compiledQuery, args)
}

// QueryNamed is a high-level function that is used to retrieve pgx.Rows from the database with named parameters
func (api *API) QueryNamed(ctx context.Context, db Querier, query string, arg interface{}) (pgx.Rows, error) {
	compiledQuery, args, err := api.dbqueryAPI.NamedQueryParams(query, arg)
	if err != nil {
		return nil, err
	}

	return api.Query(ctx, db, compiledQuery, args)
}

// Query is a wrapper around pgx's own query method
func (api *API) Query(ctx context.Context, db Querier, query string, args ...interface{}) (pgx.Rows, error) {
	return db.Query(ctx, query, args)
}

type PreparedQuery struct {
	api  *API
	prep *dbquery.PreparedQuery
}

func (api *API) PrepareNamed(query string, assertableStruct ...interface{}) (*PreparedQuery, error) {
	dbPrep, err := api.dbqueryAPI.PrepareNamed(query, assertableStruct...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &PreparedQuery{api, dbPrep}, nil
}

// SelectNamed is a high-level function that queries rows from Querier and calls the ScanAll function.
// See ScanAll for details.
func (pq *PreparedQuery) SelectNamed(ctx context.Context, db Querier, dst interface{}, arg interface{}) error {
	query, args, err := pq.prep.GetQuery(arg)
	if err != nil {
		return err
	}

	return pq.api.Select(ctx, db, dst, query, args)
}

// GetNamed is a high-level function that queries rows from Querier and calls the ScanOne function.
// See ScanOne for details.
func (pq *PreparedQuery) GetNamed(ctx context.Context, db Querier, dst interface{}, arg interface{}) error {
	query, args, err := pq.prep.GetQuery(arg)
	if err != nil {
		return err
	}

	return pq.api.Get(ctx, db, dst, query, args)
}

// ExecNamed is a high-level function that sends an executable action to the database with named parameters
func (pq *PreparedQuery) ExecNamed(ctx context.Context, db Querier, arg interface{}) (pgconn.CommandTag, error) {
	query, args, err := pq.prep.GetQuery(arg)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return pq.api.Exec(ctx, db, query, args)
}

// QueryNamed is a high-level function that is used to retrieve *sql.Rows from the database with named parameters
func (pq *PreparedQuery) QueryNamed(ctx context.Context, db Querier, arg interface{}) (pgx.Rows, error) {
	query, args, err := pq.prep.GetQuery(arg)
	if err != nil {
		return nil, err
	}

	return db.Query(ctx, query, args)
}

// NotFound is a helper function to check if an error
// is `pgx.ErrNoRows`.
func NotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// ScanAll is a wrapper around the dbscan.ScanAll function.
// See dbscan.ScanAll for details.
func (api *API) ScanAll(dst interface{}, rows pgx.Rows) error {
	err := api.dbqueryAPI.ScanAll(dst, NewRowsAdapter(rows))
	return errors.WithStack(err)
}

// ScanOne is a wrapper around the dbscan.ScanOne function.
// See dbscan.ScanOne for details. If no rows are found it
// returns a pgx.ErrNoRows error.
func (api *API) ScanOne(dst interface{}, rows pgx.Rows) error {
	err := api.dbqueryAPI.ScanOne(dst, NewRowsAdapter(rows))
	if NotFound(err) {
		return errors.WithStack(pgx.ErrNoRows)
	}
	return errors.WithStack(err)
}

// ScanRow is a wrapper around the dbscan.ScanRow function.
// See dbscan.ScanRow for details.
func (api *API) ScanRow(dst interface{}, rows pgx.Rows) error {
	err := api.dbqueryAPI.ScanRow(dst, NewRowsAdapter(rows))
	return errors.WithStack(err)
}

// RowsAdapter makes pgx.Rows compliant with the dbscan.Rows interface.
// See dbscan.Rows for details.
type RowsAdapter struct {
	pgx.Rows
}

// NewRowsAdapter returns a new RowsAdapter instance.
func NewRowsAdapter(rows pgx.Rows) *RowsAdapter {
	return &RowsAdapter{Rows: rows}
}

// Columns implements the dbscan.Rows.Columns method.
func (ra RowsAdapter) Columns() ([]string, error) {
	columns := make([]string, len(ra.Rows.FieldDescriptions()))
	for i, fd := range ra.Rows.FieldDescriptions() {
		columns[i] = string(fd.Name)
	}
	return columns, nil
}

// Close implements the dbscan.Rows.Close method.
func (ra RowsAdapter) Close() error {
	ra.Rows.Close()
	return nil
}
