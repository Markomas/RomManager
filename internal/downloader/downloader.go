package downloader

import (
	"RomManager/internal/api/romm"
	"RomManager/internal/config"
	"RomManager/internal/db"
	"RomManager/internal/db/entity"
)

type Downloader struct {
	roms   map[int]romm.Rom
	config *config.Config
	db     *db.DB
}

func NewDownloader(c *config.Config, d *db.DB) *Downloader {
	return &Downloader{
		roms:   make(map[int]romm.Rom),
		config: c,
		db:     d,
	}
}

func (d *Downloader) AddRom(rom romm.Rom) {
	d.db.CreateRommDownloadJob(&entity.RommDownloadJob{RommID: rom.ID, Name: rom.Name, Completed: nil, Progress: nil})
}

func (d *Downloader) GetDownloadJobs() ([]entity.RommDownloadJob, error) {
	return d.db.GetRommDownloadJobs()
}
