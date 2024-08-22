package main

import (
	"context"
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/NidzamuddinMuzakki/movies-abishar/common/util"
	"github.com/NidzamuddinMuzakki/movies-abishar/config"
	commonCache "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/cache"
	"github.com/NidzamuddinMuzakki/movies-abishar/repository"

	"github.com/NidzamuddinMuzakki/movies-abishar/service"
	// Import services here
	serviceHealth "github.com/NidzamuddinMuzakki/movies-abishar/service/health"

	// Import deliveries here

	httpDelivery "github.com/NidzamuddinMuzakki/movies-abishar/handler"
	httpDeliveryHealth "github.com/NidzamuddinMuzakki/movies-abishar/handler/health"

	// Import cmd here
	cmdHttp "github.com/NidzamuddinMuzakki/movies-abishar/cmd/http"

	// Import common lib here

	commonDs "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/data_source"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/logger"

	commonPanicRecover "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/middleware/gin/panic_recovery"

	commonRegistry "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"

	commonTime "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/time"
	commonValidator "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/validator"

	// Import third parties here
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	_ "github.com/spf13/viper/remote"
)

func main() {
	ctx := context.Background()

	// Start Init //
	loc, err := time.LoadLocation(commonTime.LoadTimeZoneFromEnv())
	if err != nil {
		panic(err)
	}
	time.Local = loc
	// Configuration
	config.Init()
	// Logger
	logger.Init(logger.Config{
		AppName: config.Cold.AppName,
		Debug:   config.Hot.AppDebug,
	})
	// Validator
	validator := commonValidator.New()
	// Sentry

	// Database
	// - Master
	master, err := commonDs.NewDB(&commonDs.Config{
		Driver:                config.Cold.DBMysqlMasterDriver,
		Host:                  config.Cold.DBMysqlMasterHost,
		Port:                  config.Cold.DBMysqlMasterPort,
		DBName:                config.Cold.DBMysqlMasterDBName,
		User:                  config.Cold.DBMysqlMasterUser,
		Password:              config.Cold.DBMysqlMasterPassword,
		SSLMode:               config.Cold.DBMysqlMasterSSLMode,
		MaxOpenConnections:    config.Cold.DBMysqlMasterMaxOpenConnections,
		MaxLifeTimeConnection: config.Cold.DBMysqlMasterMaxLifeTimeConnection,
		MaxIdleConnections:    config.Cold.DBMysqlMasterMaxIdleConnections,
		MaxIdleTimeConnection: config.Cold.DBMysqlMasterMaxIdleTimeConnection,
	})
	if err != nil {
		panic(err)
	}
	// - Slave
	slave, err := commonDs.NewDB(&commonDs.Config{
		Driver:                config.Cold.DBMysqlSlaveDriver,
		Host:                  config.Cold.DBMysqlSlaveHost,
		Port:                  config.Cold.DBMysqlSlavePort,
		DBName:                config.Cold.DBMysqlSlaveDBName,
		User:                  config.Cold.DBMysqlSlaveUser,
		Password:              config.Cold.DBMysqlSlavePassword,
		SSLMode:               config.Cold.DBMysqlSlaveSSLMode,
		MaxOpenConnections:    config.Cold.DBMysqlSlaveMaxOpenConnections,
		MaxLifeTimeConnection: config.Cold.DBMysqlSlaveMaxLifeTimeConnection,
		MaxIdleConnections:    config.Cold.DBMysqlSlaveMaxIdleConnections,
		MaxIdleTimeConnection: config.Cold.DBMysqlSlaveMaxIdleTimeConnection,
	})
	if err != nil {
		panic(err)
	}
	// Activity Log Client

	// Panic Recovery
	panicRecoveryMiddleware := commonPanicRecover.NewPanicRecovery(
		validator,
		commonPanicRecover.WithConfigEnv(config.Cold.AppEnv),
	)

	caches, err := commonCache.NewCache(
		commonCache.WithDriver(commonCache.RedisDriver),
		commonCache.WithHost(config.Cold.RedisHost),
		commonCache.WithDatabase("0"),
		commonCache.WithUsername(config.Cold.RedisUsername),
		commonCache.WithPassword(config.Cold.RedisPassword),
	)
	fmt.Println(err, "err cache")
	if err != nil {
		panic(err)
	}

	// Registry
	common := commonRegistry.NewRegistry(

		commonRegistry.WithValidator(validator),

		commonRegistry.WithPanicRecoveryMiddleware(panicRecoveryMiddleware),
		commonRegistry.WithCache(caches),
	)
	// End Init //

	// Start Clients //
	// ...
	// End Clients //

	// Start Repositories //
	masterUtilTx := util.NewTransactionRunner(master)
	usersRepository := repository.NewUsersRepository(common, master, slave)
	moviesRepository := repository.NewMoviesRepository(common, master, slave)

	repoRegistry := repository.NewRegistryRepository(masterUtilTx, usersRepository, moviesRepository)
	// End Repositories //

	// Start Services //
	usersService := service.NewUsersService(common, repoRegistry)
	moviesService := service.NewMoviesService(common, repoRegistry)
	healthService := serviceHealth.NewHealth(master, slave)
	serviceRegistry := service.NewRegistry(
		healthService,
		usersService,
		moviesService,
	)
	// End Deliveries //

	// Start Deliveries //
	healthDelivery := httpDeliveryHealth.NewHealth(common, healthService)
	usersDelivery := httpDelivery.NewUsers(common, serviceRegistry)
	moviesDelivery := httpDelivery.NewMovies(common, serviceRegistry)
	registryDelivery := httpDelivery.NewRegistry(healthDelivery, usersDelivery, moviesDelivery)
	// End Deliveries //

	//

	// Start HTTP Server //
	httpServer := cmdHttp.NewServer(
		common,
		registryDelivery,
	)
	httpServer.Serve(ctx)
	// End HTTP Server //
}
