package dbquery_test

import (
	"context"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	testDB  *pgxpool.Pool
	ctx     = context.Background()
	testAPI *dbscan.API
)
