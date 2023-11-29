package napi

import (
	"gorm.io/gorm"
)

type IDatabaseDriver[T any] interface {
	DB() T
}

type IRepository[T any] interface {
	IDatabaseDriver[T]
	Find(model interface{}, conds ...interface{}) (interface{}, error)
	Create(model interface{}) error
	Delete(model interface{}, id interface{}) error
	Update(model interface{}, id interface{}, values map[string]interface{}) error
	Exists(model interface{}, query interface{}, args ...interface{}) bool
}

type UpdateMap map[string]interface{}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) IRepository[*gorm.DB] {
	return &GormRepository{
		db: db,
	}
}

func (r GormRepository) DB() *gorm.DB {
	return r.db
}

func (r GormRepository) Find(model interface{}, conds ...interface{}) (interface{}, error) {
	tx := r.db.Find(&model, conds...)
	return model, tx.Error
}

func (r GormRepository) Create(model interface{}) error {
	tx := r.db.Create(model)
	return tx.Error
}

func (r GormRepository) Delete(model interface{}, id interface{}) error {
	tx := r.db.Delete(model, id)
	return tx.Error
}

func (r GormRepository) Update(model interface{}, id interface{}, values map[string]interface{}) error {
	if tx := r.db.Model(&model).Where("id = ?", id).Updates(values).Find(&model); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r GormRepository) Exists(model interface{}, query interface{}, args ...interface{}) bool {
	if tx := r.db.Where(query, args).First(&model); tx.Error != nil {
		return false
	} else {
		return tx.RowsAffected > 0
	}
}
