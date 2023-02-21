package dbquery

import (
	"regexp"
	"strings"

	"github.com/georgysavva/scany/v2/dbscan"
)

type API struct {
	scanyApi        *dbscan.API
	fieldMapperFn   NameMapperFunc // more or less should be remove, becuase they already exist in the scany api
	columnSeparator string // more or less should be remove, becuase they already exist in the scany api
	structTagKey    string // more or less should be remove, becuase they already exist in the scany api
	lexer           *Lexer
}

type APIOption func(api *API)

func NewAPI(opts ...APIOption) (*API, error) {
	api := &API{
		scanyApi:        dbscan.DefaultAPI,
		fieldMapperFn:   SnakeCaseMapper,
		columnSeparator: ".",
		structTagKey:    "db",
	}

	for _, o := range opts {
		o(api)
	}

	return api, nil
}

func (api *API) NamedQueryParams(query string, arg interface{}) (string, []interface{}, error) {
	compiledQuery, argNames, err := api.lexer.Compile(query)
	if err != nil {
		return "", nil, err
	}

	args, err := api.args(arg, argNames)
	if err != nil {
		return "", nil, err
	}

	return compiledQuery, args, nil
}

// NameMapperFunc is a function type that maps a struct field name to the database column name.
type NameMapperFunc func(string) string

var (
	matchFirstCapRe = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCapRe   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// SnakeCaseMapper is a NameMapperFunc that maps struct field to snake case.
func SnakeCaseMapper(str string) string {
	snake := matchFirstCapRe.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCapRe.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// WithStructTagKey allows to use a custom struct tag key.
// The default tag key is `db`.
func WithStructTagKey(structTagKey string) APIOption {
	return func(api *API) {
		api.structTagKey = structTagKey
	}
}

// WithLexer allows to set a custom
func WithLexer(delim rune, compileDelim DriverDelim) APIOption {
	return func(api *API) {
		api.lexer = newLexer(delim, compileDelim)
	}
}

func mustNewAPI(opts ...APIOption) *API {
	api, err := NewAPI(opts...)
	if err != nil {
		panic(err)
	}
	return api
}

var DefaultAPI = mustNewAPI()
