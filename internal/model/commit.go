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
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/config"
	"github.com/ProtobufMan/bufman/internal/constant"
	modulev1alpha "github.com/ProtobufMan/bufman/internal/gen/module/v1alpha"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Commit struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	UserID         string    `gorm:"type:varchar(64);"`
	UserName       string    `gorm:"type:varchar(200);not null"`
	RepositoryID   string    `gorm:"type:varchar(64)"`
	RepositoryName string    `gorm:"type:varchar(200)"`
	CommitID       string    `gorm:"type:varchar(64);unique;not null"`
	CommitName     string    `gorm:"type:varchar(64);unique"`
	DraftName      string    `gorm:"type:varchar(20)"`
	CreatedTime    time.Time `gorm:"autoCreateTime"`
	ManifestDigest string    `gorm:"type:string;"`

	SequenceID int64

	// 文件清单
	FileManifest *FileManifest `gorm:"foreignKey:CommitID;references:CommitID"`
	// 文件blobs
	FileBlobs FileBlobs `gorm:"foreignKey:CommitID;references:CommitID"`
	// 关联的tag
	Tags Tags `gorm:"foreignKey:RepositoryID;references:RepositoryID"`
}

func (commit *Commit) TableName() string {
	return "commits"
}

func (commit *Commit) ToProtoLocalModulePin() *modulev1alpha.LocalModulePin {
	if commit == nil {
		return (&Commit{}).ToProtoLocalModulePin()
	}

	modulePin := &modulev1alpha.LocalModulePin{
		Owner:          commit.UserName,
		Repository:     commit.RepositoryName,
		Commit:         commit.CommitName,
		CreateTime:     timestamppb.New(commit.CreatedTime),
		ManifestDigest: string(manifest.DigestTypeShake256) + ":" + commit.ManifestDigest,
	}

	if commit.DraftName == "" {
		modulePin.Branch = constant.DefaultBranch
	}

	if commit.DraftName != "" {
		modulePin.DraftName = commit.DraftName
	}

	return modulePin
}

func (commit *Commit) ToProtoModulePin() *modulev1alpha.ModulePin {
	if commit == nil {
		return (&Commit{}).ToProtoModulePin()
	}

	modulePin := &modulev1alpha.ModulePin{
		Remote:         config.Properties.BufMan.ServerHost,
		Owner:          commit.UserName,
		Repository:     commit.RepositoryName,
		Commit:         commit.CommitName,
		CreateTime:     timestamppb.New(commit.CreatedTime),
		ManifestDigest: string(manifest.DigestTypeShake256) + ":" + commit.ManifestDigest,
	}

	return modulePin
}

func (commit *Commit) ToProtoRepositoryCommit() *registryv1alpha.RepositoryCommit {
	if commit == nil {
		return (&Commit{}).ToProtoRepositoryCommit()
	}

	repositoryCommit := &registryv1alpha.RepositoryCommit{
		Id:               commit.CommitID,
		CreateTime:       timestamppb.New(commit.CreatedTime),
		Name:             commit.CommitName,
		CommitSequenceId: commit.SequenceID,
		Author:           commit.UserName,
		ManifestDigest:   string(manifest.DigestTypeShake256) + ":" + commit.ManifestDigest,
	}

	if commit.DraftName != "" {
		repositoryCommit.DraftName = commit.DraftName
	}

	if len(commit.Tags) > 0 {
		repositoryCommit.Tags = commit.Tags.ToProtoRepositoryTags()
	}

	return repositoryCommit
}

type Commits []*Commit

func (commits *Commits) ToProtoRepositoryCommits() []*registryv1alpha.RepositoryCommit {
	repositoryCommits := make([]*registryv1alpha.RepositoryCommit, len(*commits))
	for i := 0; i < len(*commits); i++ {
		repositoryCommits[i] = (*commits)[i].ToProtoRepositoryCommit()
	}

	return repositoryCommits
}

func (commits *Commits) ToProtoModulePins() []*modulev1alpha.ModulePin {
	modulePins := make([]*modulev1alpha.ModulePin, len(*commits))
	for i := 0; i < len(*commits); i++ {
		modulePins[i] = (*commits)[i].ToProtoModulePin()
	}

	return modulePins
}
