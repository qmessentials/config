package routers

import (
	"config/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RegisterUnits(unitsGroup *gin.RouterGroup, db *gorm.DB) {
	unitsGroup.GET("/", func(c *gin.Context) {
		var units []models.Unit
		result := db.Find(&units)
		if result.Error != nil {
			log.Error().Err(result.Error).Msg("Error retrieving products")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, units)
	})
}
