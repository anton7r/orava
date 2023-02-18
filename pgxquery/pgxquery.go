package pgxquery

import "github.com/anton7r/orava/dbquery"

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

