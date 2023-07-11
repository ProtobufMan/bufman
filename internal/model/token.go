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
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Token struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	UserID      string    `gorm:"type:varchar(64);not null"`
	TokenID     string    `gorm:"type:varchar(64);unique; not null"`
	TokenName   string    `gorm:"type:varchar(64);type:string"`
	CreatedTime time.Time `gorm:"autoCreateTime"`
	ExpireTime  time.Time `gorm:"not null"` // token 过期时间
	Note        string
}

func (token *Token) TableName() string {
	return "tokens"
}

func (token *Token) ToProtoToken() *registryv1alpha.Token {
	if token == nil {
		return (&Token{}).ToProtoToken()
	}

	return &registryv1alpha.Token{
		Id:         token.TokenID,
		CreateTime: timestamppb.New(token.CreatedTime),
		ExpireTime: timestamppb.New(token.ExpireTime),
		Note:       token.Note,
	}
}

type Tokens []*Token

func (tokens *Tokens) ToProtoTokens() []*registryv1alpha.Token {
	ts := make([]*registryv1alpha.Token, 0, len(*tokens))
	for i := 0; i < len(*tokens); i++ {
		ts = append(ts, (*tokens)[i].ToProtoToken())
	}

	return ts
}
