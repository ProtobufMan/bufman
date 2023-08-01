package model

import (
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type DockerRepo struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	UserID         string    `gorm:"type:varchar(64);uniqueIndex:uni_user_repo_name"`
	DockerRepoID   string    `gorm:"type:varchar(64)"`
	DockerRepoName string    `gorm:"type:varchar(200);uniqueIndex:uni_user_repo_name"` // 名称，在上传插件时需要指定名称，用户ID和名称唯一确定一个Docker repo
	Address        string    `gorm:"not null"`                                         // docker repo的地址
	UserName       string    `gorm:"not null"`                                         // docker repo的用户名
	Password       string    // docker repo的登录凭证（password or token）
	CreatedTime    time.Time `gorm:"autoCreateTime"`
	UpdateTime     time.Time `gorm:"autoUpdateTime"`
	Note           string
}

func (dockerRepo *DockerRepo) ToProtoDockerRepo() *registryv1alpha1.DockerRepo {
	if dockerRepo == nil {
		return (&DockerRepo{}).ToProtoDockerRepo()
	}

	return &registryv1alpha1.DockerRepo{
		Id:         dockerRepo.DockerRepoID,
		Name:       dockerRepo.DockerRepoName,
		Address:    dockerRepo.Address,
		Username:   dockerRepo.UserName,
		CreateTime: timestamppb.New(dockerRepo.CreatedTime),
		UpdateTime: timestamppb.New(dockerRepo.UpdateTime),
		Note:       dockerRepo.Note,
	}
}

type DockerRepos []*DockerRepo

func (dockerRepos *DockerRepos) ToProtoDockerRepos() []*registryv1alpha1.DockerRepo {
	protoDockerRepos := make([]*registryv1alpha1.DockerRepo, 0, len(*dockerRepos))

	for i := 0; i < len(*dockerRepos); i++ {
		protoDockerRepos = append(protoDockerRepos, (*dockerRepos)[i].ToProtoDockerRepo())
	}

	return protoDockerRepos
}
