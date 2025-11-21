package downloader

import (
	"RomManager/internal/api/romm"
	"RomManager/internal/config"
	"fmt"
)

type Downloader struct {
	roms   map[int]romm.Rom
	config *config.Config
}

func NewDownloader(c *config.Config) *Downloader {
	return &Downloader{
		roms:   make(map[int]romm.Rom),
		config: c,
	}
}

func (d *Downloader) AddRom(rom romm.Rom) {
	if _, exists := d.roms[rom.ID]; !exists {
		d.roms[rom.ID] = rom
	}
}

func (d *Downloader) Run() error {
	for _, rom := range d.roms {
		err := d.downloadFile(rom)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Downloader) downloadFile(rom romm.Rom) error {
	fmt.Printf("Downloading %s...\n", rom.Name)
	return nil
}
