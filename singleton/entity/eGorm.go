package entity

type EGORM interface {
	Get(dist any, id uint) error
	List(dist any, opt any) error
	Count(dist *int64, opt any) error
	Create(dist any) error
	Update(dist any, id uint) error
	Delete(id uint) error
}
