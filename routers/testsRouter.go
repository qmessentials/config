package routers

import (
	"config/models"
	"config/repositories"
	"config/utilities"
	"net/http"
	"strconv"
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
		pageSizeString := c.Query("pageSize")
		lastKeyString := c.Query("lastKey")
		namePattern := c.Query("namePattern")
		unitTypeValues := c.QueryArray("unitType")
		var criteria models.TestCriteria
		if len(namePattern) > 0 {
			criteria.NamePattern = &namePattern
		}
		if len(unitTypeValues) > 0 {
			criteria.UnitTypeValues = &unitTypeValues
		}
		pageSize, err := strconv.Atoi(pageSizeString)
		if err != nil {
			log.Error().Err(err).Msgf("Unable to parse pageSize value of '%v' as int", pageSizeString)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		var lastKey *string
		if lastKeyString != "" {
			lastKey = &lastKeyString
		}
		tests, err := testsRepo.GetMany(pageSize, lastKey, &criteria)
		if err != nil {
			log.Error().Err(err).Msg("error retrieving tests")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if tests == nil {
			log.Error().Msg("tests repo returned nil result but no error")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		log.Info().Msgf("%v tests found", len(*tests))
		c.JSON(http.StatusOK, *tests)
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
