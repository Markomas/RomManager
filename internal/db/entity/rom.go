package entity

import "gorm.io/gorm"

type Rom struct {
	gorm.Model
	Name         string
	FsName       string
	FsNameNoExt  string
	Path         string
	PlatformSlug string
	PlatformID   int
	RommId       int `gorm:"unique"`
}
