package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"createdAt,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Password     string         `json:"-"`
	Role         string         `json:"role"`
	AuthCode     int            `json:"-"`
	IsRegistered bool           `gorm:"default:0" json:"isRegistered"`
	Fishes       []Pet          `json:"-"`
}

type Pet struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Type      string  `json:"type"`
	Sex       string  `json:"sex"`
	Satiety   int     `json:"satiety"`
	LoveMeter float32 `json:"loveMeter"`
	Cost      float32 `json:"cost"`
	UserID    uint    `json:"userId"`
}
