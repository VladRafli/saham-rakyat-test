package database

import (
	"time"

	"gorm.io/gorm"
)

type Orders struct {
	ID          uint           `gorm:"primaryKey;autoIncrement;notNull" json:"id" faker:"-"`
	Name        string         `gorm:"size:255;notNull" json:"name" faker:"word" validate:"required"`
	Price       uint           `json:"price" faker:"boundary_start=1, boundary_end=1000" validate:"required,min=1"`
	ExpiredAt   time.Time      `json:"expired_at" faker:"-"`
	HistoriesID *uint          `gorm:"default:null" faker:"-"`
	CreatedAt   time.Time      `gorm:"autoCreateTime:milli" json:"created_at" faker:"-"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at" faker:"-"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" faker:"-"`
}

type Users struct {
	ID          uint           `gorm:"primaryKey;autoIncrement;notNull" json:"id" faker:"-"`
	FullName    string         `gorm:"size:255;notNull" json:"full_name" faker:"name" validate:"required"`
	FirstOrder  bool           `gorm:"default:true;notNull" json:"first_order"` // what is this? for now, I assume this is checking if user is first time ordering
	HistoriesID *uint          `gorm:"default:null" faker:"-"`
	CreatedAt   time.Time      `gorm:"autoCreateTime:milli" json:"created_at" faker:"-"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at" faker:"-"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" faker:"-"`
}

type Histories struct {
	ID           uint           `gorm:"primaryKey;autoIncrement;notNull" json:"id" faker:"-"`
	User         *Users         `gorm:"foreignKey:HistoriesID" faker:"-"`
	Orders       []Orders       `gorm:"foreignKey:HistoriesID" faker:"-"`
	Descriptions string         `gorm:"size:255;notNull" json:"descriptions" faker:"sentence"`
	CreatedAt    time.Time      `gorm:"autoCreateTime:milli" json:"created_at" faker:"-"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at" faker:"-"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" faker:"-"`
}
