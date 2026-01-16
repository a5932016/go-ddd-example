package filters

import "github.com/a5932016/go-ddd-example/util/pb"

type SortFilter struct {
	Asc  string `json:"asc"`
	Desc string `json:"desc"`
	// Pairs []string `json:"pairs"`
}

func NewSortFilterFromPb(p *pb.SortFilter) *SortFilter {
	if p == nil {
		return nil
	}

	return &SortFilter{
		Asc:  p.Asc,
		Desc: p.Desc,
	}
}
