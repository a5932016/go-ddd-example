package filters

import "github.com/a5932016/go-ddd-example/util/pb"

// BooleanFilter boolean filter
type BooleanFilter struct {
	Is        *bool `json:"is,omitempty"`
	IsPresent *bool `json:"is_present,omitempty"`
}

// ToSQL implement filter adapter
func (f *BooleanFilter) ToSQL(column string) (queries []string, args []interface{}) {
	queries, args = []string{}, []interface{}{}
	if len(column) == 0 {
		return
	}
	if f.Is != nil {
		queries = append(queries, column+" = ?")
		args = append(args, f.Is)
	}
	return
}

func NewBooleanFilterFromPb(p *pb.BooleanFilter) *BooleanFilter {
	if p == nil {
		return nil
	}

	b := BooleanFilter{}
	if p.Is != nil {
		b.Is = p.Is
	}
	if p.IsPresent != nil {
		b.IsPresent = p.IsPresent
	}

	return &b
}
