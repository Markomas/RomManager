package entity

import (
	"time"

	"gorm.io/gorm"
)

type RommDownloadJob struct {
	gorm.Model
	Name       string
	RommID     int `gorm:"unique"`
	Completed  *bool
	Progress   *float64
	LockedTill *time.Time
	Error      string
}
