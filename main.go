package main

import (
	"config/models"
	"config/routers"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "test",
		})
	})

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/New_York", os.Getenv("PGHOST"),
		os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"), os.Getenv("PGPORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&models.Product{})
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
	}

	routers.RegisterProducts(r.Group("/products"), db)

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
