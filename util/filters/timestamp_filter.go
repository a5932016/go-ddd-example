package filters

type TimestampFilter struct {
	After     *int64  `json:"after,omitempty"`
	Before    *int64  `json:"before,omitempty"`
	On        *int64  `json:"on,omitempty"`
	Between   []int64 `json:"between,omitempty"`
	IsPresent *bool   `json:"is_present,omitempty"`
}

func (t TimestampFilter) ToSQL(column string) (queries []string, args []interface{}) {
	queries, args = []string{}, []interface{}{}
	if len(column) == 0 {
		return
	}
	if t.After != nil {
		queries = append(queries, column+" > to_timestamp(?)")
		args = append(args, *t.After)
	}
	if t.Before != nil {
		queries = append(queries, column+" < to_timestamp(?)")
		args = append(args, *t.Before)
	}
	if t.On != nil {
		queries = append(queries, column+" = to_timestamp(?)")
		args = append(args, *t.On)
	}
	if t.Between != nil && len(t.Between) == 2 {
		queries = append(queries, column+" BETWEEN to_timestamp(?) AND to_timestamp(?)")
		args = append(args, (t.Between)[0], (t.Between)[1])
	}
	return
}
