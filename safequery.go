package safequery

import (
	"regexp"
	"strconv"

	"github.com/lib/pq"
)

type Query struct {
	query string
	args  []any
}

var matchDollar = regexp.MustCompile(`\s\$\$?\d+`)

func New(text string, args ...any) *Query {
	q := Query{}
	q.Add(text, args...)
	return &q
}

func (q *Query) Add(text string, args ...any) {
	shift := len(q.args)

	q.query += matchDollar.ReplaceAllStringFunc(text, func(match string) string {
		match = match[1:]

		if match[1] == '$' {
			// double dollar mode

			index, err := strconv.Atoi(match[2:])
			if err != nil {
				panic(err)
			}

			identifier := pq.QuoteIdentifier(args[index-1].(string))
			shift-- // FIXME
			args = append(args[:index-1], args[index:]...)

			return " " + identifier
		}

		index, err := strconv.Atoi(match[1:])
		if err != nil {
			panic(err)
		}

		return " $" + strconv.Itoa(shift+index)
	})
	q.args = append(q.args, args...)
}

func (q *Query) Query() string {
	return q.query
}

func (q *Query) Args() []any {
	return q.args
}
