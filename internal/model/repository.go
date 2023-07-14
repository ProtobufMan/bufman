/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/bufman/registry/v1alpha"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Repository 仓库
type Repository struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	UserID         string    `gorm:"type:varchar(64);uniqueIndex:uni_user_id_name"` // 所属用户，与仓库名组成唯一索引
	UserName       string    `gorm:"type:varchar(200);not null"`
	RepositoryID   string    `gorm:"type:varchar(64);unique;not null"`
	RepositoryName string    `gorm:"type:varchar(200);uniqueIndex:uni_user_id_name"` // 仓库名，与拥有者组成唯一索引
	CreatedTime    time.Time `gorm:"autoCreateTime"`
	UpdateTime     time.Time `gorm:"autoUpdateTime"`
	Visibility     uint8     `gorm:"default:1"` // 可见性，1:public 2:private
	Deprecated     bool      // 是否弃用
	DeprecationMsg string    // 弃用说明
	Description    string    // 描述信息

	// 拥有的draft
	DraftCommits []*Commit `gorm:"foreignKey:RepositoryID;references:RepositoryID"`
	// 拥有的tag
	Tags []*Tag `gorm:"foreignKey:RepositoryID;references:RepositoryID"`
}

func (repository *Repository) TableName() string {
	return "repositories"
}

func (repository *Repository) ToProtoRepository() *registryv1alpha.Repository {
	if repository == nil {
		return (&Repository{}).ToProtoRepository()
	}

	return &registryv1alpha.Repository{
		Id:                 repository.RepositoryID,
		CreateTime:         timestamppb.New(repository.CreatedTime),
		UpdateTime:         timestamppb.New(repository.UpdateTime),
		Name:               repository.RepositoryName,
		UserId:             repository.UserID,
		Visibility:         registryv1alpha.Visibility(repository.Visibility),
		Deprecated:         repository.Deprecated,
		DeprecationMessage: repository.DeprecationMsg,
		OwnerName:          repository.UserName,
		Description:        repository.Description,
	}
}

type Repositories []*Repository

func (repositoryEntities *Repositories) ToProtoRepositories() []*registryv1alpha.Repository {
	repositories := make([]*registryv1alpha.Repository, 0, len(*repositoryEntities))

	for i := 0; i < len(*repositoryEntities); i++ {
		repositories = append(repositories, (*repositoryEntities)[i].ToProtoRepository())
	}

	return repositories
}

type RepositoryCounts struct {
	TagsCount   int64
	DraftsCount int64
}

func (repositoryCounts *RepositoryCounts) ToProtoRepositoryCounts() *registryv1alpha.RepositoryCounts {
	if repositoryCounts == nil {
		return (&RepositoryCounts{}).ToProtoRepositoryCounts()
	}

	return &registryv1alpha.RepositoryCounts{
		TagsCount:   uint32(repositoryCounts.TagsCount),
		DraftsCount: uint32(repositoryCounts.DraftsCount),
	}
}
