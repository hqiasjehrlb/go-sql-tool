package sqltool

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
)

// Query single row into T
//
// ExtFields are columns not in struct T but returns by query.
// *T will be nil if no rows returned.
func QueryRow[T any](ctx Context, db DBConn, query string, args ...any) (result *T, ext ExtFieldMap, err error) {
	rows, exts, err := QueryRows[T](ctx, db, query, args...)
	if err != nil {
		return
	}
	if len(rows) > 0 && len(exts) > 0 {
		result = &rows[0]
		ext = exts[0]
	}
	return
}

// Query rows into []T
//
// []ExtFields are rows with columns not in struct T but returns by query.
// []T will be empty slice when no rows returned.
func QueryRows[T any](ctx Context, db DBConn, query string, args ...any) (results []T, exts []ExtFieldMap, err error) {
	results = []T{}
	exts = []ExtFieldMap{}

	data := new(T)
	rv := reflect.ValueOf(data).Elem()

	// Check if T is struct or not
	if rv.Type().Kind() != reflect.Struct {
		err = errors.New("type of T must be a struct")
		return
	}

	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	var cols []string
	cols, err = rows.Columns()
	if err != nil {
		return
	}

	addrMap := make(map[string]any)

	// Prepare column name & pointer mapping
	for i := 0; i < rv.NumField(); i++ {
		prop := rv.Field(i)
		colname := rv.Type().Field(i).Tag.Get("db")
		if colname == "" {
			colname = rv.Type().Field(i).Name
		}
		addrMap[colname] = prop.Addr().Interface()
	}

	var extCols []string
	scanDest := make([]any, len(cols))

	// Prepare scan dest array
	for i, colname := range cols {
		if _, ok := addrMap[colname]; !ok {
			extCols = append(extCols, colname)
			var v any
			addrMap[colname] = &v
		}
		scanDest[i] = addrMap[colname]
	}

	// Scan rows
	for rows.Next() {
		err = rows.Scan(scanDest...)
		if err != nil {
			return
		}
		m := make(ExtFieldMap)
		for _, colname := range extCols {
			p := addrMap[colname]
			if r := reflect.ValueOf(p); r.Kind() == reflect.Pointer {
				m[colname] = r.Elem().Interface()
			}
		}
		results = append(results, *data)
		exts = append(exts, m)
	}

	return
}

// Alias for context.Context
type Context = context.Context

// This is an interface for *sql.DB or *sql.Conn
type DBConn interface {
	QueryContext(Context, string, ...any) (*sql.Rows, error)
}
