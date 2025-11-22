package entity

import (
	"time"

	"gorm.io/gorm"
)

type SaveState struct {
	gorm.Model
	RommID           int `gorm:"unique"`
	Rom              Rom
	FileName         string
	LocalPath        *string
	VersionUpdatedAt time.Time
	Md5Hash          string
}
