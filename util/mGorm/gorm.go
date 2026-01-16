package mGorm

import (
	"reflect"
	"strings"

	"gorm.io/gorm"

	"github.com/a5932016/go-ddd-example/util/filters"
)

type DB struct {
	*gorm.DB
}

func New(db *gorm.DB) *DB {
	return &DB{
		DB: db,
	}
}

// OrderWithFilter order with filter
func (db *DB) OrderWithFilter(filter *filters.SortFilter) *DB {
	if filter != nil {
		if len(filter.Asc) > 0 {
			db.DB = db.Order(filter.Asc + " ASC")
		} else if len(filter.Desc) > 0 {
			db.DB = db.Order(filter.Desc + " DESC")
		}
	} else {
		db.DB = db.Order("id DESC")
	}

	return db
}

// WhereWithStringFilter where with string filter
func (db *DB) WhereWithStringFilter(column string, filter *filters.StringFilter, operator string) *DB {
	if filter == nil || len(column) == 0 {
		return db
	}
	if filter.In != nil {
		db.DB = db.WhereWithOp(operator, column+" IN (?)", filter.In)
	}
	if filter.NotIn != nil {
		db.DB = db.WhereWithOp(operator, column+" NOT IN (?)", filter.NotIn)
	}
	if filter.Is != nil {
		db.DB = db.WhereWithOp(operator, column+" = ?", *filter.Is)
	}
	if filter.IsNot != nil {
		db.DB = db.WhereWithOp(operator, column+" != ?", *filter.IsNot)
	}
	if filter.StartsWith != nil {
		db.DB = db.WhereWithOp(operator, column+" Like ?", *filter.StartsWith+"%")
	}
	if filter.EndWith != nil {
		db.DB = db.WhereWithOp(operator, column+" Like ?", "%"+*filter.EndWith)
	}
	if filter.Like != nil {
		db.DB = db.WhereWithOp(operator, column+" Like ?", "%"+*filter.Like+"%")
	}
	if filter.LikeSlice != nil {
		for _, like := range filter.LikeSlice {
			db.DB = db.WhereWithOp(operator, column+" Like ?", "%"+like+"%")
		}
	}
	return db
}

// FiltersToQuery filters to mGORM Query
func (db *DB) FiltersToQuery(filters map[string]filters.Adaptor) (query Query) {
	query = Query{
		queries: []string{},
		args:    []interface{}{},
	}
	if filters == nil {
		return
	}
	for column, filter := range filters {
		vf := reflect.ValueOf(filter)
		if len(column) == 0 || vf.IsNil() {
			continue
		}
		fQueries, fArgs := filter.ToSQL(column)
		query.queries = append(query.queries, fQueries...)
		query.args = append(query.args, fArgs...)
	}

	return
}

// WhereWithNumberFilter where with number filter
func (db *DB) WhereWithNumberFilter(column string, filter *filters.NumberFilter, operator string) *DB {
	if filter == nil || len(column) == 0 {
		return db
	}
	if filter.Is != nil {
		db.DB = db.WhereWithOp(operator, column+" = ?", filter.Is)
	}
	if filter.IsNot != nil {
		db.DB = db.WhereWithOp(operator, column+" != ?", filter.Is)
	}
	if filter.In != nil {
		db.DB = db.WhereWithOp(operator, column+" IN (?)", filter.In)
	}
	if filter.NotIn != nil {
		db.DB = db.WhereWithOp(operator, column+" NOT IN (?)", filter.NotIn)
	}
	if filter.Gt != nil {
		db.DB = db.WhereWithOp(operator, column+" > ?", filter.Gt)
	}
	if filter.Gte != nil {
		db.DB = db.WhereWithOp(operator, column+" >= ?", filter.Gte)
	}
	if filter.Lt != nil {
		db.DB = db.WhereWithOp(operator, column+" < ?", filter.Lt)
	}
	if filter.Lte != nil {
		db.DB = db.WhereWithOp(operator, column+" <= ?", filter.Lte)
	}
	if filter.Between != nil && len(filter.Between) == 2 {
		db.DB = db.WhereWithOp(operator, column+" BETWEEN ? AND ?", filter.Between[0], filter.Between[1])
	}

	return db
}

// WhereWithBoolFilter where with bool filter
func (db *DB) WhereWithBoolFilter(column string, filter *filters.BooleanFilter, operator string) *DB {
	if filter == nil || len(column) == 0 {
		return db
	}
	if filter.Is != nil {
		db.DB = db.WhereWithOp(operator, column+" = ?", filter.Is)
	}

	return db
}

// WhereWithOp where condition with operator or / and
func (db *DB) WhereWithOp(operator string, query interface{}, args ...interface{}) *gorm.DB {
	opFn := db.Where
	if strings.EqualFold(operator, "or") {
		opFn = db.Or
	}

	return opFn(query, args...)
}
