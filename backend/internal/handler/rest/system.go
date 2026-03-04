package rest

import (
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
)

func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		util.SubsonicOK(c, nil)
	}
}

func GetLicense() gin.HandlerFunc {
	return func(c *gin.Context) {
		type License struct {
			Valid bool   `xml:"valid,attr" json:"valid"`
			Email string `xml:"email,attr" json:"email"`
		}
		util.SubsonicOK(c, gin.H{
			"license": License{Valid: true, Email: "pineapple@local"},
		})
	}
}
