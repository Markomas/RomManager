package save_state_sync

import (
	"RomManager/internal/api"
	"RomManager/internal/config"
	"RomManager/internal/db"
	"RomManager/internal/db/entity"
	"fmt"
)

type SaveStateSync struct {
	config *config.Config
	api    *api.Romm
	db     *db.DB
}

func (s SaveStateSync) Sync() error {
	roms := s.db.GetAllRoms()
	for _, rom := range roms {
		s.downloadRemoteSaveState(rom)
	}

	return nil
}

func NewSaveStateSync(c *config.Config, a *api.Romm, d *db.DB) *SaveStateSync {
	return &SaveStateSync{
		config: c,
		api:    a,
		db:     d,
	}
}

func (s SaveStateSync) downloadRemoteSaveState(rom entity.Rom) {
	saveStates, err := s.api.GetSaveStates(rom.RommId, rom.PlatformID)
	if err != nil {
		fmt.Println("Error fetching save states for rom %d: %v", rom.RommId, err)
		return
	}

	for _, saveState := range saveStates {
		if saveState.MissingFromFS == true {
			continue
		}

		_, err = s.api.DownloadSaveStateToTmp(saveState, rom)
		if err != nil {
			return
		}

		//saveState := &entity.SaveState{
		//	RommID:           rom.RommId,
		//	FileName:         saveState.FileName,
		//	LocalPath:        saveState.LocalPath,
		//	VersionUpdatedAt: saveState.VersionUpdatedAt,
		//	Md5Hash:          saveState.Md5Hash,
		//}
		//state, err := s.db.GetSaveState(saveState.ID)

	}
}
