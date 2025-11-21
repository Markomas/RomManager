package db

import (
	"RomManager/internal/config"
	"RomManager/internal/db/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	config *config.Config
	db     *gorm.DB
}

func New(c *config.Config) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(c.System.DBPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&entity.RommDownloadJob{})
	if err != nil {
		return nil, err
	}

	return &DB{config: c, db: db}, nil
}

func (d *DB) Close() {

}

func (d *DB) GetRommDownloadJobs() ([]entity.RommDownloadJob, error) {
	var job []entity.RommDownloadJob
	err := d.db.Order("completed ASC, progress ASC, id DESC").Find(&job).Error
	return job, err
}
