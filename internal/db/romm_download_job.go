package db

import (
	"RomManager/internal/db/entity"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm/clause"
)

func (d *DB) CreateRommDownloadJob(e *entity.RommDownloadJob) {
	d.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "romm_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "completed", "progress"}),
	}).Create(e)
}
func (d *DB) GetNextRommDownloadJob() (*entity.RommDownloadJob, error) {
	var job entity.RommDownloadJob
	d.db.Where("completed IS NULL AND (locked_till < ? OR locked_till IS NULL) and error IS NULL", time.Now()).First(&job)
	if job.ID != 0 {
		d.db.Model(&job).Update(
			"locked_till",
			time.Now().Add(time.Second*time.Duration(d.config.System.DownloadQueueTimeoutSeconds)),
		)

		return &job, nil
	}
	return nil, errors.New("no romm jobs available")
}

func (d *DB) UpdateRommDownloadJobProgress(jobId uint, progress float64) error {
	res := d.db.Model(&entity.RommDownloadJob{}).
		Where("id = ?", jobId).
		Updates(map[string]interface{}{
			"progress":    progress,
			"locked_till": time.Now().Add(time.Second * time.Duration(d.config.System.DownloadQueueTimeoutSeconds)),
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("no rows updated for job id %d", jobId)
	}

	return nil
}

func (d *DB) UpdateRommDownloadJobProgressAsCompleted(i uint, err string) {
	d.db.Model(&entity.RommDownloadJob{}).Where("id =?", i).Updates(map[string]interface{}{
		"completed":   true,
		"locked_till": nil,
		"error":       err,
	})
}
