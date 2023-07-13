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

// FileManifest 保存文件清单，记录每一次提交的所有文件
type FileManifest struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Digest   string // 文件清单哈希
	CommitID string `gorm:"type:varchar(64), unique"`
}

type FileManifests []*FileManifest

// FileBlob 保存文件blob
type FileBlob struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Digest   string // 文件哈希
	CommitID string `gorm:"type:varchar(64);index"`
	FileName string
}

type FileBlobs []*FileBlob

//// FileIdentity 唯一文件表，根据哈希值记录文件实际存储地址
//type FileIdentity struct {
//	ID       int64  `gorm:"primaryKey;autoIncrement"`
//	Digest   string `gorm:"unique"` // 文件哈希
//	Location string // 文件实际存储地址
//	FileName string // 文件实际存储名字
//}
