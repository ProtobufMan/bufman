package model

import "time"

type DockerRepo struct {
	ID             int64  `gorm:"primaryKey;autoIncrement"`
	UserID         string `gorm:"type:varchar(64);uniqueIndex:uni_user_repo_name"`
	DockerRepoID   string
	DockerRepoName string    `gorm:"uniqueIndex:uni_user_repo_name"` // 名称，在上传插件时需要指定名称，用户ID和名称唯一确定一个Docker repo
	Address        string    `gorm:"not null"`                       // docker repo的地址
	UserName       string    // docker repo的用户名
	Password       string    // docker repo的登录凭证（password or token）token注意登录expire time
	CreatedTime    time.Time `gorm:"autoCreateTime"`
	UpdateTime     time.Time `gorm:"autoUpdateTime"`
	Note           string
}

type DockerRepos []*DockerRepo
