package entity

import (
	"time"

	"gorm.io/gorm"
)

type SaveState struct {
	gorm.Model
	RomID            uint
	Rom              Rom `gorm:"foreignKey:RomID"`
	RommID           int `gorm:"unique"`
	FileName         string
	LocalPath        *string
	VersionUpdatedAt time.Time
	Md5Hash          string
}
