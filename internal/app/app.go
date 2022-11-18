package app

import (
	"fmt"
	"nextclan/validator-register/mobile-validator-scheduler-service/config"
	v1 "nextclan/validator-register/mobile-validator-scheduler-service/internal/controller/http/v1"
	usecase "nextclan/validator-register/mobile-validator-scheduler-service/internal/usecase"
	"nextclan/validator-register/mobile-validator-scheduler-service/internal/usecase/repo"
	"nextclan/validator-register/mobile-validator-scheduler-service/pkg/httpserver"
	"nextclan/validator-register/mobile-validator-scheduler-service/pkg/logger"
	mongodb "nextclan/validator-register/mobile-validator-scheduler-service/pkg/mongo"
	"nextclan/validator-register/mobile-validator-scheduler-service/pkg/redis"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	fmt.Println("Starting App...")

	mongoDb := mongodb.New(cfg.Mongo.ConnectionUri)
	deviceRepository := repo.NewDeviceRepository(mongoDb.Client.Database(cfg.Mongo.Database).Collection(cfg.Mongo.DeviceCollectionName))
	// HTTP Server
	httpServer := initializeHttp(l, cfg)

	//Scheduler
	if cfg.DeviceFilterJobs.IsRunJob {
		redisPoolJobsClient := redis.NewRedisPoolClient(cfg.Addr, cfg.Password, cfg.RedisScheduleDB)
		clearSchedule(cfg, redisPoolJobsClient)
		poolJobs := initScheduler(cfg, l, deviceRepository, redisPoolJobsClient)
		shutdownApplicationHandler(l, httpServer, mongoDb, redisPoolJobsClient, poolJobs)

	} else {
		// Shutdown
		shutdownApplicationHandler(l, httpServer, mongoDb, nil, nil)
	}
}

func initializeHttp(l *logger.Logger, cfg *config.Config) *httpserver.Server {
	handler := gin.New()
	handler.Use(cors.Default())
	v1.NewRouter(handler, l)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	return httpServer
}

func clearSchedule(cfg *config.Config, rej redis.RedisInterface) {
	client := work.NewClient(cfg.DeviceFilterJobs.ScheduleNamespace, rej.Pool())
	results, _, _ := client.ScheduledJobs(1)
	for _, v := range results {
		client.DeleteScheduledJob(v.RunAt, v.ID)
	}
	retryJobs, _, _ := client.RetryJobs(1)
	for _, v := range retryJobs {
		client.DeleteRetryJob(v.RetryAt, v.ID)
	}
}

func initScheduler(cfg *config.Config, l logger.Interface, r *repo.DeviceRepositoryMongo, rej redis.RedisInterface) *work.WorkerPool {

	poolJobs := work.NewWorkerPool(usecase.ScheduleUpdateMobileValidatorDeviceStatus{}, 1, cfg.DeviceFilterJobs.ScheduleNamespace, rej.Pool())
	poolJobs.PeriodicallyEnqueue(cfg.DeviceFilterJobs.ScheduleCron, cfg.DeviceFilterJobs.JobName)

	poolJobs.Middleware(func(c *usecase.ScheduleUpdateMobileValidatorDeviceStatus, job *work.Job, next work.NextMiddlewareFunc) error {
		deviceCache := redis.NewRedisClient(cfg.Redis.Addr, cfg.Password, cfg.Redis.RedisDeviceDB)
		c.Log = l
		c.DeviceRepository = r
		c.RedisClient = deviceCache
		c.DeviceFilterConfig = &cfg.DeviceFilterJobs
		c.Log.Info("repository", r)

		return next()
	})
	poolJobs.Middleware((*usecase.ScheduleUpdateMobileValidatorDeviceStatus).NotifyStartJob)
	poolJobs.Job(cfg.DeviceFilterJobs.JobName, (*usecase.ScheduleUpdateMobileValidatorDeviceStatus).Execute)
	poolJobs.Start()

	return poolJobs
}

func shutdownApplicationHandler(l *logger.Logger, httpServer *httpserver.Server, mongoDB *mongodb.MongoDB, redisPoolJobClient redis.RedisInterface, poolJobs *work.WorkerPool) {
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	}
	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	mongoDB.Close()

	if poolJobs != nil {
		poolJobs.Stop()
	}
	if redisPoolJobClient != nil {
		if err = redisPoolJobClient.Close(); err != nil {
			l.Error(fmt.Errorf("app - Closing - Redis Connection Poll: %w", err))
		}
	}
}
