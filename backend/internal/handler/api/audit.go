package api

import (
	"net/http"

	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListAuditLogs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params util.PaginationParams
		c.ShouldBindQuery(&params)
		params.Normalize()

		var logs []model.AuditLog
		var total int64

		query := db.Model(&model.AuditLog{})

		if action := c.Query("action"); action != "" {
			query = query.Where("action = ?", action)
		}

		query.Count(&total)
		query.Offset(params.Offset).Limit(params.Limit).Order("id DESC").Find(&logs)

		c.JSON(http.StatusOK, gin.H{
			"total":  total,
			"offset": params.Offset,
			"limit":  params.Limit,
			"items":  logs,
		})
	}
}
