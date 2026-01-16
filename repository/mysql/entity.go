package mysql

import (
	"fmt"

	"github.com/a5932016/go-ddd-example/model"
	"github.com/a5932016/go-ddd-example/util/mGorm"
	"gorm.io/gorm"
)

func GetEntityDB(db *gorm.DB, opt model.EntityOption) *gorm.DB {
	mDB := mGorm.New(db)

	if opt.Keyword != nil {
		keyword := fmt.Sprint(*opt.Keyword, "%")
		mDB.DB = mDB.DB.Where("name LIKE ?", keyword)
	}

	mDB = mDB.OrderWithFilter(opt.SortBy)

	if opt.Offset != nil && opt.Limit != nil {
		db = db.Limit(*opt.Limit)
		db = db.Offset(*opt.Offset)
	}

	return mDB.DB
}
