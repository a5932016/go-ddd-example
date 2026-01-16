package filters

import (
	"regexp"
	"strings"

	"github.com/a5932016/go-ddd-example/util/pb"
)

// StringFilter string filter
type StringFilter struct {
	Is         *string  `json:"is,omitempty"`
	IsNot      *string  `json:"is_not,omitempty"`
	StartsWith *string  `json:"starts_with,omitempty"`
	EndWith    *string  `json:"end_with,omitempty"`
	In         []string `json:"in,omitempty"`
	NotIn      []string `json:"not_in,omitempty"`
	Like       *string  `json:"like,omitempty"`
	LikeSlice  []string `json:"like[],omitempty"`
}

// ToSQL implement filter adapter
func (f StringFilter) ToSQL(column string) (queries []string, args []interface{}) {
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
	if f.StartsWith != nil {
		queries = append(queries, column+" Like ?")
		args = append(args, *f.StartsWith+"%")
	}
	if f.EndWith != nil {
		queries = append(queries, column+" Like ?")
		args = append(args, "%"+*f.EndWith)
	}
	if f.Like != nil {
		queries = append(queries, column+" Like ?")
		args = append(args, "%"+*f.Like+"%")
	}
	if f.LikeSlice != nil {
		for _, like := range f.LikeSlice {
			queries = append(queries, column+" Like ?")
			args = append(args, "%"+like+"%")
		}
	}
	return
}

// ToQueryString to query string
func (f StringFilter) ToQueryString() map[string]string {
	query := map[string]string{}
	if f.Is != nil {
		query["is"] = *f.Is
	}
	if f.IsNot != nil {
		query["is_not"] = *f.IsNot
	}
	if f.StartsWith != nil {
		query["starts_with"] = *f.StartsWith
	}
	if f.In != nil {
		query["in"] = strings.Join(f.In, ",")
	}
	if f.NotIn != nil {
		query["not_in"] = strings.Join(f.NotIn, ",")
	}
	if f.Like != nil {
		query["like"] = *f.Like
	}
	return query
}

func NewStringFilterFromPb(p *pb.StringFilter) *StringFilter {
	if p == nil {
		return nil
	}

	s := StringFilter{}
	if p.Is != nil {
		s.Is = p.Is
	}
	if p.IsNot != nil {
		s.IsNot = p.IsNot
	}
	if p.StartsWith != nil {
		s.StartsWith = p.StartsWith
	}
	if p.EndWith != nil {
		s.EndWith = p.EndWith
	}
	if p.In != nil {
		s.In = p.In
	}
	if p.NotIn != nil {
		s.NotIn = p.NotIn
	}
	if p.Like != nil {
		s.Like = p.Like
	}
	if p.Likes != nil {
		s.LikeSlice = p.Likes
	}

	return &s
}

func CutString(s string, length int) []string {
	if length <= 0 {
		return []string{""} // Return nil for invalid length or empty input for consistency
	}

	// Remove HTML tags
	re := regexp.MustCompile("<[^>]*>")
	s = re.ReplaceAllString(s, "")

	// Trim spaces after removing HTML tags to avoid leading/trailing spaces in chunks
	s = strings.TrimSpace(s)

	var chunks []string
	for len(s) > 0 {
		// If the remaining string is shorter than the length, take the whole string
		if len(s) <= length {
			chunks = append(chunks, s)
			break
		}

		// Find the next chunk boundary without cutting a word in the middle
		boundary := length
		if boundary < len(s) && s[boundary] != ' ' && (boundary+1 < len(s) && s[boundary+1] != ' ') {
			// If we're in the middle of a word, move back to the nearest space
			for boundary > 0 && s[boundary-1] != ' ' {
				boundary--
			}
		}

		// If we didn't find a space, just cut at the length
		if boundary == 0 {
			boundary = length
		}

		chunks = append(chunks, s[:boundary])
		s = strings.TrimSpace(s[boundary:]) // Trim leading spaces for the next chunk
	}

	return chunks
}
