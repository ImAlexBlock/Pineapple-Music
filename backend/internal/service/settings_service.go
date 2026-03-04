package service

import (
	"pineapple-music/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SettingsService struct {
	DB *gorm.DB
}

func (s *SettingsService) Get(key string) string {
	var setting model.Setting
	silent := s.DB.Session(&gorm.Session{Logger: s.DB.Logger.LogMode(logger.Silent)})
	if err := silent.Where("key = ?", key).First(&setting).Error; err != nil {
		return ""
	}
	return setting.Value
}

func (s *SettingsService) Set(key, value string) error {
	return s.DB.Save(&model.Setting{Key: key, Value: value}).Error
}

func (s *SettingsService) GetAll() ([]model.Setting, error) {
	var settings []model.Setting
	err := s.DB.Find(&settings).Error
	return settings, err
}

func (s *SettingsService) Delete(key string) error {
	return s.DB.Delete(&model.Setting{}, "key = ?", key).Error
}

// GetBool returns a boolean setting with a default value.
func (s *SettingsService) GetBool(key string, defaultVal bool) bool {
	v := s.Get(key)
	if v == "" {
		return defaultVal
	}
	return v == "true" || v == "1"
}
