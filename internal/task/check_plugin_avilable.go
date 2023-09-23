package task

import (
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"time"
)

const (
	checkPluginAvailableGroup    = "check_plugin_available"
	checkPluginAvailableKey      = "check_plugin_available"
	checkPluginAvailableInterval = time.Hour
)

func AddCheckPluginAvailableJob() error {
	return AsyncTaskManager.AddJob(checkPluginAvailableGroup, time.Second, checkPluginAvailableKey, pullCheckPluginAvailableJob())
}

func AddNextCheckPluginAvailableJob() error {
	return AsyncTaskManager.AddJob(checkPluginAvailableGroup, checkPluginAvailableInterval, checkPluginAvailableKey, pullCheckPluginAvailableJob())
}

// 生成一个验证registry auth的异步任务
func pullCheckPluginAvailableJob() func() {
	return func() {
		logger.Infof("Running a %v job", checkPluginAvailableGroup)
		defer func() {
			_ = AddNextCheckPluginAvailableJob()
		}()
		pluginMapper := &mapper.PluginMapperImpl{}

		plugins, err := pluginMapper.FindAllUnavailable()
		if err != nil {
			return
		}

		for _, plugin := range plugins {
			_ = AddPullPluginImageJob(plugin.PluginID)
		}
	}
}
