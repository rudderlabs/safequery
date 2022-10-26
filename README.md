# safequery

[![Tests safequery](https://github.com/rudderlabs/safequery/actions/workflows/test-safequery.yaml/badge.svg?branch=main)](https://github.com/rudderlabs/safequery/actions/workflows/test-safequery.yaml)

The goal of this library is to simplify the usage of secure query practices, with minimal abstraction.

## Incremental and conditional construction

Many times we have to create a query in multiple steps, some of them conditional. In [such conditions](https://github.com/rudderlabs/rudder-server/blob/master/jobsdb/jobsdb.go#L2522) using prepared statements can be be complex.

This library can simplify this by normalizing indices:

```go
    q := safequery.New("SELECT * FROM table WHERE ")

    if searchByID {
        q.Add("id = $1", 1)
    } else {
        q.Add("name = $1", "John")
    }

    q.Add("AND updated_at > $1", time.Now())


    // "SELECT * FROM table WHERE id = $1 AND updated_at > $2", 1, timestamp
    // OR
    // "SELECT * FROM table WHERE name = $1 AND updated_at > $2", John, timestamp
    db.QueryContext(ctx, q.Query(), q.Args()...)
```

_safequery is responsible to construct the query correctly, no need to keep tract of dollar prepared parameters indices._

## Escape identifiers

It is very common in our codebase to use variable identifiers like table names or columns. Although most databases don't support prepared statements for them, they should be at least sanitized.

```go
    q := safequery.New("SELECT * FROM $$1 WHERE id = $2", "table_name", 1)

    // "SELECT * FROM \"table_name\" WHERE id = $1"
    q.Query()

    // []any{1}
    q.Args()
```

safequery will use db specific code to sanitize the identifier.

## Named parameters

Not all SQL databases support [named parameters](https://pkg.go.dev/database/sql#NamedArg). It is a very continent feature, when dealing with complex queries. safequery can map named parameters to indexed ones for databases that don't support them.

```go
    q := safequery.New("SELECT * FROM $$table_name WHERE id = $id",
        sql.NameArgs{Name:"table_name" Value: "users"}, 
        sql.NameArgs{Name:"id" Value: 1},
    )

    // "SELECT * FROM \"users\" WHERE id = $1"
    q.Query()

    // []any{1}
    q.Args()
```
