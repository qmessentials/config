package main

import (
	"config/repositories"
	"config/routers"
	"config/utilities"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
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

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
		panic(err)
	}
	defer conn.Close(context.Background())
	tableVersionsRepository := repositories.NewTableVersionsRepository(conn)
	if err = tableVersionsRepository.Migrate(); err != nil {
		panic(err)
	}
	configSettingsRepository := repositories.NewConfigSettingsRepository(conn)
	if err = configSettingsRepository.Migrate(); err != nil {
		panic(err)
	}
	unitsRepository := repositories.NewUnitsRepository(conn)
	if err = unitsRepository.Migrate(); err != nil {
		panic(err)
	}
	productsRepository := repositories.NewProductsRepository(conn)
	if err = productsRepository.Migrate(); err != nil {
		panic(err)
	}
	testsRepository := repositories.NewTestsRepository(conn)
	if err = testsRepository.Migrate(); err != nil {
		panic(err)
	}
	if err = utilities.BootstrapConfig(configSettingsRepository, unitsRepository); err != nil {
		panic(err)
	}

	redis := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDRESS"), Password: "", DB: 0}) //default DB
	cacheService := utilities.NewCacheService(redis)
	authClient := utilities.NewAuthClient(os.Getenv("AUTH_SERVICE_ENDPOINT"))
	permissionsHelper := utilities.NewPermissionHelper(authClient, cacheService)

	routers.RegisterProducts(r.Group("/products"), productsRepository, permissionsHelper)
	routers.RegisterTests(r.Group("/tests"), testsRepository, permissionsHelper)
	routers.RegisterUnits(r.Group("/units"), unitsRepository, permissionsHelper)

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
