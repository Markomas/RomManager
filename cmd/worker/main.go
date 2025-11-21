package main

import (
	"RomManager/internal/api"
	"RomManager/internal/config"
	"RomManager/internal/db"
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

	for true {
		job, err := db.GetNextRommDownloadJob()
		if err != nil {
			time.Sleep(time.Second * 10)
			continue
		}
		fmt.Println("Processing job:", job.ID)
		err = a.DownloadRomm(job.RommID, func(progress float64) {
			fmt.Printf("Downloading romm id: %d: %.2f%%\n", job.RommID, progress*100)
			err = db.UpdateRommDownloadJobProgress(job.ID, progress)
			if err != nil {
				fmt.Println("Error updating download job progress:", err)
			}
		})
		if err != nil {
			fmt.Println("Error downloading romm:", err)
		}
		fmt.Println("Romm download completed")
		db.UpdateRommDownloadJobProgressAsCompleted(job.ID, err)
	}
}
