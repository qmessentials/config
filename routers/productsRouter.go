package routers

import (
	"config/models"
	"config/utilities"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RegisterProducts(products *gin.RouterGroup, db *gorm.DB, permissionsHelper *utilities.PermissionsHelper) {
	products.GET("/:productCode", func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		log.Info().Msgf("Auth header is %s", authHeader)
		bearerPattern := regexp.MustCompile("(?i)^bearer (.*)$")
		tokens := bearerPattern.FindStringSubmatch(authHeader)
		if len(tokens) != 2 {
			log.Warn().Msg("Unauthenticated attempt to retrieve config data")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		isAllowed, err := permissionsHelper.IsAuthorized(tokens[1], "View a Product")
		if err != nil {
			log.Error().Err(err).Msg("Error checking permissions")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !isAllowed {
			log.Warn().Msgf("Failed permission check")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		var product models.Product
		result := db.First(&product, "product_code=?", c.Param("productCode"))
		if result.Error != nil {
			if result.Error.Error() == "record not found" {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			log.Error().Err(result.Error).Msg("Error retrieving product")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, product)
	})
	products.GET("/", func(c *gin.Context) {
		var products []models.Product
		result := db.Find(&products)
		if result.Error != nil {
			log.Error().Err(result.Error).Msg("Error retrieving products")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, products)
	})
	products.POST("/", func(c *gin.Context) {
		var product models.Product
		err := c.BindJSON(&product)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		result := db.Create(&product)
		if result.Error != nil {
			log.Error().Err(result.Error).Msg("Error creating product")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected != 1 {
			log.Error().Err(result.Error).Msgf("Product creation returned %d rows affected (expected 1)", result.RowsAffected)
			c.AbortWithStatus((http.StatusInternalServerError))
			return
		}
		c.Status(http.StatusCreated)
	})
	products.PUT("/:productCode", func(c *gin.Context) {
		var product models.Product
		err := c.BindJSON(&product)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if product.ProductCode != c.Param("productCode") {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		result := db.Model(&product).Where("product_code = ?", product.ProductCode).Update("description", product.Description)
		if result.Error != nil {
			log.Error().Err(result.Error).Msgf("Error updating product %s", c.Param("productCode"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected != 1 {
			log.Error().Err(result.Error).Msgf("Product update returned %d rows affected (expected 1)", result.RowsAffected)
			c.AbortWithStatus((http.StatusInternalServerError))
			return
		}
		c.Status(http.StatusOK)
	})
	products.DELETE("/:productCode", func(c *gin.Context) {
		result := db.Where("product_code", c.Param("productCode")).Delete(&models.Product{})
		if result.Error != nil {
			log.Error().Err(result.Error).Msgf("Error deleting product %s", c.Param("productCode"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected != 1 {
			log.Error().Err(result.Error).Msgf("Product delete returned %d rows affected (expected 1)", result.RowsAffected)
			c.AbortWithStatus((http.StatusInternalServerError))
			return
		}
		c.Status(http.StatusOK)
	})
}