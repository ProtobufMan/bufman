package task

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/core/docker"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"time"
)

const (
	pullPluginImageGroup = "pull_plugin_image"
	pullPluginImageTry   = 3
	maxPullNum           = 3
)

func AddPullPluginImageJob(pluginID string) error {
	return AsyncTaskManager.AddJob(pullPluginImageGroup, time.Second, pluginID, pullPluginImageJob(pluginID))
}

// 生成一个验证registry auth的异步任务
func pullPluginImageJob(pluginID string) func() {

	return func() {
		logger.Infof("Running a %v job", pullPluginImageGroup)

		pluginMapper := &mapper.PluginMapperImpl{}
		dockerRepoMapper := &mapper.DockerRepoMapperImpl{}

		plugin, err := pluginMapper.FindByPluginID(pluginID)
		if err != nil {
			return
		}

		if plugin.Deprecated || plugin.IsAvailable || plugin.TryNum > maxPullNum {
			return
		}

		registry, err := dockerRepoMapper.FindByDockerRepoID(plugin.DockerRepoID)
		if err != nil {
			return
		}

		// 进行多次尝试
		for i := 0; i < pullPluginImageTry; i++ {
			// pull image
			dockerCli, err := docker.NewDockerClient(registry.Address, registry.UserName, registry.Password)
			err = dockerCli.TryPullImage(context.Background(), plugin.ImageName, plugin.ImageDigest)
			if err != nil {
				continue
			}

			// pull 成功
			_ = pluginMapper.UpdateAvailable(pluginID, true)
			return
		}

		// pull 失败
		_ = pluginMapper.UpdateAvailable(pluginID, false)
		return
	}
}
