// Package sqlxext - provides handy extensions to sqlx default functions.
package sqlxext

import (
    "context"

    "github.com/jmoiron/sqlx"
)

// NamedGetContext gets and scans one row it to dst like sqlx.GetContext does, but accepts query with named parameters.
// Nil inParams is allowed, unlike dest. Use ExecSQLForRows (NamedSelectContext) if you need to discard row(-s) or Exec
// the statement.
func NamedGetContext(
        ctx context.Context,
        db sqlx.ExtContext,
        dest interface{}, query string, inParams interface{},
) error {
    query, args, err := safeBindNamed(db, query, inParams)
    if err != nil {
        return err
    }

    return sqlx.GetContext(ctx, db, dest, query, args...)
}

// NamedSelectContext select all rows and scans them to dst like sqlx.SelectContext does, but accepts query with named
// parameters. Nil dest or nil inParams are allowed. If dest is nil, it works like sqlx.NamedExecContext.
func NamedSelectContext(
        ctx context.Context,
        db sqlx.ExtContext,
        dest interface{}, query string, inParams interface{},
) error {
    query, args, err := safeBindNamed(db, query, inParams)
    if err != nil {
        return err
    }
    if dest == nil {
        dest = &[]interface{}{}
    }

    return sqlx.SelectContext(ctx, db, dest, query, args...)
}

// ExecSQLForRows execute db request and write results to dest. It is sqlx.SelectContext under the hood, but
// accepts query with named parameters. Nil dest or nil inParams are allowed. If dest is nil, it works like
// sqlx.NamedExecContext.
func ExecSQLForRows(
        ctx context.Context,
        db sqlx.ExtContext,
        dest interface{}, query string, inParams interface{},
) error {
    return NamedSelectContext(ctx, db, dest, query, inParams)
}

// safeBindNamed safe binding params method (may take a nil as input parameter).
// nolint:unnamedResult // unneeded unnamed check
func safeBindNamed(db sqlx.ExtContext, query string, inParams interface{}) (string, []interface{}, error) {
    if inParams == nil {
        inParams = struct{}{}
    }
    // transform named parameters to '?'
    query, args, err := sqlx.Named(query, inParams)
    if err != nil {
        return "", nil, err
    }
    // map parameters '?' by number of arguments in IN-arg, if its present
    query, args, err = sqlx.In(query, args...)
    if err != nil {
        return "", nil, err
    }
    // transform parameters to $1, $2, etc - driver format
    return db.Rebind(query), args, nil
}
