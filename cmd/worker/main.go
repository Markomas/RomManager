package main

import (
	"RomManager/internal/api"
	"RomManager/internal/config"
	"RomManager/internal/db"
	"RomManager/internal/save_state_sync"
	"fmt"
	"time"
)

func main() {
	c, err := config.New("config.yml")
	if err != nil {
		panic(err)
	}

	db, err := db.New(c)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	a := api.New(c)
	saveStateSync := save_state_sync.NewSaveStateSync(c, a, db)

	for true {
		err = downloadCheck(db, a)
		if err != nil {
			fmt.Println("Error processing download job:", err)
		}
		err = saveStateSync.Sync()
		if err != nil {
			fmt.Println("Error syncing save states:", err)
		}

		time.Sleep(time.Second * 60)
	}
}

func downloadCheck(db *db.DB, a *api.Romm) error {
	for {
		job, err := db.GetNextRommDownloadJob()
		if err != nil {
			return nil // No more jobs
		}

		fmt.Println("Processing job:", job.ID)
		rom, outFilePath, downloadErr := a.DownloadRomm(job.RommID, func(progress float64) {
			fmt.Printf("Downloading romm id: %d: %.2f%%\n", job.RommID, progress*100)
			err := db.UpdateRommDownloadJobProgress(job.ID, progress)
			if err != nil {
				fmt.Println("Error updating download job progress:", err)
			}
		})

		var erroras *string
		if downloadErr != nil {
			errStr := downloadErr.Error()
			erroras = &errStr
			fmt.Println("Error downloading rom:", downloadErr)
		} else {
			erroras = nil
		}
		db.UpdateRommDownloadJobProgressAsCompleted(job.ID, rom, outFilePath, erroras)
	}
}
