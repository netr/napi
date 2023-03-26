package napi

import (
	"gorm.io/gorm"
)

type IBaseRepository interface {
	DB() *gorm.DB
	Get(model interface{}, conds ...interface{}) (interface{}, error)
	Create(model interface{}) error
	Delete(model interface{}, id interface{}) error
	Update(model interface{}, id interface{}, values map[string]interface{}) error
	Exists(model interface{}, query interface{}, args ...interface{}) bool
}

type UpdateMap map[string]interface{}

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) IBaseRepository {
	return &BaseRepository{
		db: db,
	}
}

func (r BaseRepository) DB() *gorm.DB {
	return r.db
}

func (r BaseRepository) Get(model interface{}, conds ...interface{}) (interface{}, error) {
	tx := r.db.Find(&model, conds...)
	return model, tx.Error
}

func (r BaseRepository) Create(model interface{}) error {
	tx := r.db.Create(model)
	return tx.Error
}

func (r BaseRepository) Delete(model interface{}, id interface{}) error {
	tx := r.db.Delete(model, id)
	return tx.Error
}

func (r BaseRepository) Update(model interface{}, id interface{}, values map[string]interface{}) error {
	if tx := r.db.Model(&model).Where("id = ?", id).Updates(values).Find(&model); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r BaseRepository) Exists(model interface{}, query interface{}, args ...interface{}) bool {
	if tx := r.db.Where(query, args).First(&model); tx.Error != nil {
		return false
	} else {
		return tx.RowsAffected > 0
	}
}
