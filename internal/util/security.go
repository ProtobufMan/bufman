package util

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/ProtobufMan/bufman/internal/constant"
	"strconv"
	"time"
)

// EncryptPlainPassword 加密明文密码
func EncryptPlainPassword(userName, plainPwd string) string {
	sha := sha256.New()
	sha.Write([]byte(plainPwd))
	sha.Write([]byte(userName))
	bytes := sha.Sum(nil)

	return hex.EncodeToString(bytes)
}

// GenerateToken 生成token
func GenerateToken(username, note string) string {
	sha := sha256.New()
	// 以用户名 note 时间戳做哈希
	sha.Write([]byte(username))
	sha.Write([]byte(note))
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	sha.Write([]byte(now))

	return hex.EncodeToString(sha.Sum(nil)[:constant.TokenLength])
}

func GenerateCommitName(userName, RepositoryName string) string {
	sha := sha256.New()
	// 以用户名 note 时间戳做哈希
	sha.Write([]byte(userName))
	sha.Write([]byte(RepositoryName))
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	sha.Write([]byte(now))
	bytes := sha.Sum(nil)

	return hex.EncodeToString(bytes[:constant.CommitLength])
}
