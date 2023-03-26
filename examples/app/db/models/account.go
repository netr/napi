package models

import (
	"github.com/bxcodec/faker/v3"
	"time"
)

type Account struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `valid:"required" json:"username" gorm:"type: varchar(16); unique;"`
	Password  string    `valid:"required" json:"-" gorm:"type: varchar(32);"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoUpdateTime"`
}

func (f *Account) Make() interface{} {
	model := &Account{
		Username: faker.Username(),
		Password: faker.Password(),
	}
	return model
}
