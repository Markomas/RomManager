package db

import (
	"RomManager/internal/config"
	"RomManager/internal/db/entity"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	config *config.Config
	db     *gorm.DB
}

func New(c *config.Config) (*DB, error) {
	if _, err := os.Stat(c.System.DBPath); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(c.System.DBPath), 0755)
		if err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(sqlite.Open(c.System.DBPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&entity.RommDownloadJob{}, &entity.Rom{})
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

func (d *DB) GetAllRoms() []entity.Rom {
	var roms []entity.Rom
	d.db.Find(&roms)
	return roms
}

func (d *DB) GetSaveState(id int) (*entity.SaveState, error) {
	saveState := &entity.SaveState{}
	return saveState, d.db.Where("romm_id =?", id).First(saveState).Error
}
