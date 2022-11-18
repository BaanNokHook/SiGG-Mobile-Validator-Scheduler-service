package usecase

import (
	"context"
	"nextclan/validator-register/mobile-validator-scheduler-service/config"
	"nextclan/validator-register/mobile-validator-scheduler-service/internal/usecase/repo"
	"nextclan/validator-register/mobile-validator-scheduler-service/pkg/logger"
	"nextclan/validator-register/mobile-validator-scheduler-service/pkg/redis"
	"time"

	"github.com/gocraft/work"
)

type ScheduleUpdateMobileValidatorDeviceStatus struct {
	Log                logger.Interface
	DeviceRepository   repo.DeviceRepository
	RedisClient        redis.RedisInterface
	DeviceFilterConfig *config.DeviceFilterJobs
}

func (s *ScheduleUpdateMobileValidatorDeviceStatus) NotifyStartJob(job *work.Job, next work.NextMiddlewareFunc) error {
	s.Log.Info("Starting job: ", job.Name)
	return next()
}

func (s *ScheduleUpdateMobileValidatorDeviceStatus) Execute(job *work.Job) error {
	currentTime := time.Now().Add(time.Minute * -time.Duration(s.DeviceFilterConfig.DeviceOfflineThresholdMinute)).Unix()
	latestValidate := time.Now().Add(time.Minute * -time.Duration(s.DeviceFilterConfig.DeviceValidateThresholdMinute)).Unix()
	offlineDevice, err := s.DeviceRepository.FindDeviceByLatestSyncLte(context.TODO(), currentTime, 1000)
	if err != nil {
		s.Log.Info("query offline device err: ", err)
		return nil
	}
	s.Log.Info("offline Device: ", offlineDevice)
	for _, v := range offlineDevice {
		s.Log.Info("del redis key: ", v)
		s.RedisClient.Del(v.PublicKey)
	}
	s.Log.Info("finish del offline device")
	onlineDevice, err := s.DeviceRepository.FindDeviceByLatestSyncGte(context.TODO(), currentTime, 1000)
	if err != nil {
		s.Log.Info("query online device err: ", err)
		return nil
	}
	s.Log.Info("online Device: ", onlineDevice)
	for _, v := range onlineDevice {
		s.Log.Info("set redis key: ", v)
		if v.LatestValidateTransaction < latestValidate {
			s.RedisClient.Set(v.PublicKey, v.UserId)
		}
	}
	key, _ := s.RedisClient.GetRandomKey()
	s.Log.Info("RedisClient.GetRandomKey: ", key)
	if key == "" && len(onlineDevice) > 0 {
		for _, v := range onlineDevice {
			s.RedisClient.Set(v.PublicKey, v.UserId)
		}
	}
	s.Log.Info("Finish job: ", job.Name)
	s.RedisClient.Close()
	return nil
}
