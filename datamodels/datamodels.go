package datamodels

import (
	"time"
)

type PK struct {
	ID uint64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

type Model struct {
	CreatedAt time.Time  `json:"created_at" form:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" form:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" form:"deleted_at" sql:"index"`
}
