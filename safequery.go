package safequery

import (
	"database/sql"
	"regexp"
	"strconv"

	"github.com/lib/pq"
)

type Query struct {
	query string
	args  []any
}

var matchDollar = regexp.MustCompile(`\s\$\$?[A-Za-z0-9]+`)

func New(text string, args ...any) *Query {
	q := &Query{}
	return q.Add(text, args...)
}

func (q *Query) Add(text string, args ...any) *Query {
	named := make(map[string]any)
	for _, arg := range args {
		if namedArg, ok := arg.(sql.NamedArg); ok {
			named[namedArg.Name] = namedArg.Value
		}
	}

	q.query += matchDollar.ReplaceAllStringFunc(text, func(match string) string {
		match = match[2:]
		var identifier bool
		if match[0] == '$' {
			identifier = true
			match = match[1:]
		}

		var value any
		if '0' <= match[0] && match[0] <= '9' {
			index, err := strconv.Atoi(match)
			if err != nil {
				panic(err)
			}
			value = args[index-1]
		} else {
			if param, ok := named[match]; ok {
				value = param
			} else {
				panic("no named param " + match)
			}
		}

		if identifier {
			identifier := pq.QuoteIdentifier(value.(string))
			return " " + identifier
		} else {
			q.args = append(q.args, value)
			return " $" + strconv.Itoa(len(q.args))
		}
	})

	return q
}

func (q *Query) Query() string {
	return q.query
}

func (q *Query) Args() []any {
	return q.args
}
