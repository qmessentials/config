package main

import (
	"config/models"
	"config/routers"
	"config/utilities"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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

	redis := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDRESS"), Password: "", DB: 0}) //default DB
	cacheService := utilities.NewCacheService(redis)
	authClient := utilities.NewAuthClient(os.Getenv("AUTH_SERVICE_ENDPOINT"))
	permissionsHelper := utilities.NewPermissionHelper(authClient, cacheService)

	routers.RegisterProducts(r.Group("/products"), db, permissionsHelper)

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
