package service

import (
	"pineapple-music/internal/model"

	"gorm.io/gorm"
)

type AuditService struct {
	DB *gorm.DB
}

func (s *AuditService) Log(action, role, ip, detail string) {
	s.DB.Create(&model.AuditLog{
		Action: action,
		Role:   role,
		IP:     ip,
		Detail: detail,
	})
}
