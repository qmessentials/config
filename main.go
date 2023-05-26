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
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
		panic("unable to connect to database")
	}
	if err = db.AutoMigrate(&models.Product{}); err != nil {
		log.Error().Err(err).Msg("Error automigrating model Product")
	}
	if err = db.AutoMigrate(&models.Test{}); err != nil {
		log.Error().Err(err).Msg("Error automigrating model Test")
	}
	if err = db.AutoMigrate(&models.Unit{}); err != nil {
		log.Error().Err(err).Msg("Error automigrating model Unit")
	}
	if err = db.AutoMigrate(&models.ConfigSetting{}); err != nil {
		log.Error().Err(err).Msg("Error automigrating model ConfigSetting")
	}
	if err = utilities.BootstrapConfig(db); err != nil {
		log.Error().Err(err).Msgf("Error bootstrapping config data")
		panic("database isn't bootstrapped")
	}

	redis := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDRESS"), Password: "", DB: 0}) //default DB
	cacheService := utilities.NewCacheService(redis)
	authClient := utilities.NewAuthClient(os.Getenv("AUTH_SERVICE_ENDPOINT"))
	permissionsHelper := utilities.NewPermissionHelper(authClient, cacheService)

	routers.RegisterProducts(r.Group("/products"), db, permissionsHelper)
	routers.RegisterTests(r.Group("/tests"), db, permissionsHelper)
	routers.RegisterUnits(r.Group("/units"), db)

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
