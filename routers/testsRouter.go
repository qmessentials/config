package routers

import (
	"config/models"
	"config/repositories"
	"config/utilities"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RegisterTests(testsGroup *gin.RouterGroup, testsRepo *repositories.TestsRepository, permissionsHelper *utilities.PermissionsHelper) {
	testsGroup.GET("/:testName", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "test-view", permissionsHelper)
		if permissionsResult != http.StatusOK {
			log.Warn().Msg("failed permission request for test-view")
			c.AbortWithStatus(permissionsResult)
			return
		}
		test, err := testsRepo.GetOne(c.Param("testName"))
		if err != nil {
			log.Error().Err(err).Msg("failed retrieving test")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if test == nil {
			c.AbortWithStatus(http.StatusNotFound)
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
		tests, err := testsRepo.GetMany()
		if err != nil {
			log.Error().Err(err).Msg("error retrieving tests")
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
		if err := c.BindJSON(&test); err != nil {
			log.Error().Err(err).Msg("request body could not be bound")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := testsRepo.Create(&test); err != nil {
			log.Error().Err(err).Msg("error creating test")
			c.AbortWithStatus(http.StatusInternalServerError)
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
			log.Warn().Msg("request body could not be bound")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if !strings.EqualFold(test.TestName, c.Param("testName")) {
			log.Warn().Msg("test name in request body does not match request name in URL")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := testsRepo.Update(&test); err != nil {
			log.Error().Err(err).Msgf("error updating test '%s'", c.Param("testName"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
}
