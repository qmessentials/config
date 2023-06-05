package routers

import (
	"config/repositories"
	"config/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RegisterUnits(unitsGroup *gin.RouterGroup, repo *repositories.UnitsRepository, permissionsHelper *utilities.PermissionsHelper) {
	unitsGroup.GET("/", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "unit-search", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		units, err := repo.GetMany()
		if err != nil {
			log.Error().Err(err).Msg("error retrieving units")
		}
		c.JSON(http.StatusOK, units)
	})
}
