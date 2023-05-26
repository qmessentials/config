package routers

import (
	"config/models"
	"config/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RegisterTests(testsGroup *gin.RouterGroup, db *gorm.DB, permissionsHelper *utilities.PermissionsHelper) {
	testsGroup.GET("/:testName", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "test-view", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		var test models.Test
		result := db.First(&test, "test_name=?", c.Param("testName"))
		if result.Error != nil {
			if result.Error.Error() == "record not found" {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			log.Error().Err(result.Error).Msg("Error retrieving test")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, test)
	})
	testsGroup.GET("/", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "test-search", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		var tests []models.Test
		result := db.Find(&tests)
		if result.Error != nil {
			log.Error().Err(result.Error).Msg("Error retrieving tests")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, tests)
	})
	testsGroup.POST("/", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "test-create", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		var test models.Test
		err := c.BindJSON(&test)
		if err != nil {
			log.Error().Err(err).Msg("Request body could not be bound")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		result := db.Create(&test)
		if result.Error != nil {
			log.Error().Err(result.Error).Msg("Error creating test")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected != 1 {
			log.Error().Err(result.Error).Msgf("Test creation returned %d rows affected (expected 1)", result.RowsAffected)
			c.AbortWithStatus((http.StatusInternalServerError))
			return
		}
		c.Status(http.StatusCreated)
	})
	testsGroup.PUT("/:testName", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "test-edit", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		var test models.Test
		if err := c.BindJSON(&test); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if test.TestName != c.Param("testName") {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		result := db.Model(&test).Where("test_name = ?", test.TestName).Updates(models.Test{
			UnitType:           test.UnitType,
			References:         test.References,
			Standards:          test.Standards,
			AvailableModifiers: test.AvailableModifiers})
		if result.Error != nil {
			log.Error().Err(result.Error).Msgf("Error updating test '%s'", c.Param("testName"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected != 1 {
			log.Error().Err(result.Error).Msgf("Test update returned %d rows affected (expected 1)", result.RowsAffected)
			c.AbortWithStatus((http.StatusInternalServerError))
			return
		}
		c.Status(http.StatusOK)
	})
	testsGroup.DELETE("/:testName", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "test-remove", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		result := db.Where("test_name", c.Param("testName")).Delete(&models.Test{})
		if result.Error != nil {
			log.Error().Err(result.Error).Msgf("Error deleting test '%s'", c.Param("testName"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected != 1 {
			log.Error().Err(result.Error).Msgf("Test delete returned %d rows affected (expected 1)", result.RowsAffected)
			c.AbortWithStatus((http.StatusInternalServerError))
			return
		}
		c.Status(http.StatusOK)
	})
}
