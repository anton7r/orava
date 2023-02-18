package pgxquery_test

import (
	"context"

	"github.com/anton7r/orava/dbquery"
	"github.com/anton7r/orava/pgxquery"
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var (
	testDB  *pgxpool.Pool
	ctx     = context.Background()
	testAPI *dbscan.API
)

func getAPI() (*dbquery.API, error) {
	dbqueryAPI, err := pgxscan.NewDBScanAPI(dbscan.WithLexer(':', dbscan.SequentialDollarDelim))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	api, err := pgxquery.NewAPI(dbqueryAPI)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return api, nil
}
