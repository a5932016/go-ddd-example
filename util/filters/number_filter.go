package filters

// NumberFilter number filter
type NumberFilter struct {
	Is      *int  `json:"is,omitempty"`
	IsNot   *int  `json:"is_not,omitempty"`
	In      []int `json:"in,omitempty"`
	NotIn   []int `json:"not_in,omitempty"`
	Lt      *int  `json:"lt,omitempty"`
	Lte     *int  `json:"lte,omitempty"`
	Gt      *int  `json:"gt,omitempty"`
	Gte     *int  `json:"gte,omitempty"`
	Between []int `json:"between,omitempty"`
}

// ToSQL implement filter adapter
func (f NumberFilter) ToSQL(column string) (queries []string, args []interface{}) {
	queries, args = []string{}, []interface{}{}
	if len(column) == 0 {
		return
	}
	if f.In != nil {
		queries = append(queries, column+" IN (?)")
		args = append(args, f.In)
	}
	if f.NotIn != nil {
		queries = append(queries, column+" NOT IN (?)")
		args = append(args, f.NotIn)
	}
	if f.Is != nil {
		queries = append(queries, column+" = ?")
		args = append(args, f.Is)
	}
	if f.IsNot != nil {
		queries = append(queries, column+" != ?")
		args = append(args, f.IsNot)
	}
	if f.Gt != nil {
		queries = append(queries, column+" > ?")
		args = append(args, f.Gt)
	}
	if f.Gte != nil {
		queries = append(queries, column+" >= ?")
		args = append(args, f.Gte)
	}
	if f.Lt != nil {
		queries = append(queries, column+" < ?")
		args = append(args, f.Lt)
	}
	if f.Lte != nil {
		queries = append(queries, column+" <= ?")
		args = append(args, f.Lte)
	}
	if f.Between != nil && len(f.Between) == 2 {
		queries = append(queries, column+" BETWEEN ? AND ?")
		args = append(args, f.Between[0], f.Between[1])
	}
	return
}
