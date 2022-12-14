package safequery_test

import (
	"database/sql"
	"testing"

	"github.com/rudderlabs/safequery"
	"github.com/stretchr/testify/require"
)

func TestMultipleAdds(t *testing.T) {
	q := safequery.New("SELECT * FROM table WHERE id = $1", 1)

	t.Log("every add should be independent of previous adds")
	q.Add(" AND name = $1 OR name = $2", "John", "Jane")

	require.Equal(t, "SELECT * FROM table WHERE id = $1 AND name = $2 OR name = $3", q.Query())
	require.Equal(t, []any{1, "John", "Jane"}, q.Args())
}

func TestDoubleDollar(t *testing.T) {
	q := safequery.New("SELECT $1, * FROM $$2 WHERE id = $3", true, "table_name", 1)

	require.Equal(t, "SELECT $1, * FROM \"table_name\" WHERE id = $2", q.Query())
	require.Equal(t, []any{true, 1}, q.Args())

	t.Run("complex example", func(t *testing.T) {
		q.Add("SELECT $1, * FROM $$2 JOIN $$3 WHERE id = $4", true, "table_name", "other table", 1)
	})
}

func TestNamedArg(t *testing.T) {
	q := safequery.New(
		"SELECT * FROM $$table WHERE id = $id",
		sql.NamedArg{Name: "table", Value: "table_name"},
		sql.NamedArg{Name: "id", Value: 1},
	).Add(" AND name = $name", sql.NamedArg{Name: "name", Value: "John"})

	require.Equal(t, `SELECT * FROM "table_name" WHERE id = $1 AND name = $2`, q.Query())
	require.Equal(t, []any{1, "John"}, q.Args())
}
