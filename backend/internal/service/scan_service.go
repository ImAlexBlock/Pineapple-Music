package service

import (
	"fmt"

	"pineapple-music/internal/model"
	"pineapple-music/internal/scanner"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ScanService struct {
	DB       *gorm.DB
	MusicDir string
}

// StartScan begins an asynchronous scan.
func (s *ScanService) StartScan() (*model.ScanJob, error) {
	if scanner.IsScanning() {
		return nil, fmt.Errorf("a scan is already in progress")
	}

	job := &model.ScanJob{Status: "pending"}
	if err := s.DB.Create(job).Error; err != nil {
		return nil, err
	}

	go scanner.Scan(s.DB, s.MusicDir, job)

	return job, nil
}

// GetScanJob returns the latest scan job.
func (s *ScanService) GetLatestJob() (*model.ScanJob, error) {
	var job model.ScanJob
	silent := s.DB.Session(&gorm.Session{Logger: s.DB.Logger.LogMode(logger.Silent)})
	if err := silent.Order("id DESC").First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// GetScanJob returns a specific scan job.
func (s *ScanService) GetScanJob(id uint) (*model.ScanJob, error) {
	var job model.ScanJob
	if err := s.DB.First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}
