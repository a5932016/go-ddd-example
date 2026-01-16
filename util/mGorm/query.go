package mGorm

import (
	"fmt"
	"strings"

	"github.com/a5932016/go-ddd-example/util/filters"
)

type Query struct {
	queries []string
	args    []interface{}
}

func (q *Query) AppendFilter(column string, f filters.Adaptor) *Query {
	fQueries, fArgs := f.ToSQL(column)
	q.queries = append(q.queries, fQueries...)
	q.args = append(q.args, fArgs...)
	return q
}

func (q *Query) Append(query string, arg ...interface{}) *Query {
	q.queries = append(q.queries, query)
	q.args = append(q.args, arg...)
	return q
}

func (q *Query) ToSQLArgs(operator string) (sql string, args []interface{}) {
	op := "AND"
	if strings.EqualFold(operator, "or") {
		op = "OR"
	}
	op = fmt.Sprintf(" %s ", op)
	sql = fmt.Sprintf(fmt.Sprintf("(%s)", strings.Join(q.queries, op)))

	return sql, q.args
}

func (q *Query) IsNil() bool {
	return len(q.args) == 0
}
