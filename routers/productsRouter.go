package routers

import (
	"config/models"
	"config/repositories"
	"config/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RegisterProducts(productsGroup *gin.RouterGroup, productsRepo *repositories.ProductsRepository, permissionsHelper *utilities.PermissionsHelper) {
	productsGroup.GET("/:productCode", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "product-view", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		product, err := productsRepo.GetOne(c.Param("productCode"))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if product == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, product)
	})
	productsGroup.GET("/", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "product-search", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		products, err := productsRepo.GetMany()
		if err != nil {
			log.Error().Err(err).Msg("error retrieving products")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, products)
	})
	productsGroup.POST("/", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "product-create", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			log.Warn().Msg("failed to bind request body to models.Product")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := productsRepo.Create(&product); err != nil {
			log.Error().Err(err).Msg("Error creating product")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusCreated)
	})
	productsGroup.PUT("/:productCode", func(c *gin.Context) {
		permissionsResult := checkPermissions(c, "product-edit", permissionsHelper)
		if permissionsResult != http.StatusOK {
			c.AbortWithStatus(permissionsResult)
			return
		}
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			log.Warn().Msg("failed to bind request body to models.Product")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if product.ProductCode != c.Param("productCode") {
			log.Warn().Msg("product code in request body does not match URL")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := productsRepo.Update(&product); err != nil {
			log.Error().Err(err).Msgf("Error updating product %s", c.Param("productCode"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
}
