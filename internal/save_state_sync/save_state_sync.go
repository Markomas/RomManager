package save_state_sync

import (
	"RomManager/internal/api"
	"RomManager/internal/config"
	"RomManager/internal/db"
	"RomManager/internal/db/entity"
	"fmt"
	"os"
	"path/filepath"
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

		saveTmpPath, err := s.api.DownloadSaveStateToTmp(saveState, rom)
		if err != nil {
			continue
		}

		md5Hash, err := s.api.CalculateFileMd5(*saveTmpPath)
		if err != nil {
			fmt.Println("Error calculating md5 for %s: %v", *saveTmpPath, err)
			continue
		}

		fmt.Println("Downloaded save state %s to %s", saveState.FileName, *saveTmpPath)

		localSaveState := &entity.SaveState{
			RomID:            rom.ID,
			RommID:           rom.RommId,
			FileName:         saveState.FileName,
			LocalPath:        saveTmpPath,
			VersionUpdatedAt: saveState.UpdatedAt,
			Md5Hash:          md5Hash,
		}
		_, err = s.db.GetSaveStateByHash(localSaveState.Md5Hash)
		if err == nil {
			continue
		}

		destinationPath := filepath.Join(s.config.System.SaveStatesPath, filepath.Base(*saveTmpPath))
		if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
			fmt.Printf("Error creating destination directory: %v\n", err)
			continue
		}
		if err := os.Rename(*saveTmpPath, destinationPath); err != nil {
			fmt.Printf("Error moving file: %v\n", err)
			continue
		}
		s.db.AddSaveState(localSaveState)
		fmt.Printf("Saved new save state %s for rom %d\n", saveState.FileName, rom.RommId)

	}

}
