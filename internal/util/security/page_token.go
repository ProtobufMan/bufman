package security

import (
	"errors"
	"github.com/ProtobufMan/bufman/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type PageTokenChaim struct {
	PageOffset int
	jwt.RegisteredClaims
}

func GenerateNextPageToken(lastPageOffset, lastPageSize, lastDataLength int) (string, error) {
	if lastDataLength < lastPageSize {
		// 已经查询完了
		return "", nil
	}

	nextPageOffset := lastPageOffset + lastDataLength
	// 定义 token 的过期时间
	now := time.Now()
	expireTime := now.Add(config.Properties.BufMan.PageTokenExpireTime)

	// 创建一个自定义的 Claim
	chaim := &PageTokenChaim{
		PageOffset: nextPageOffset,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "bufman",
		},
	}

	// 使用 JWT 签名算法生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, chaim)

	// 将 token 进行加盐加密
	tokenString, err := token.SignedString([]byte(config.Properties.BufMan.PageTokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParsePageToken(tokenString string) (*PageTokenChaim, error) {
	if tokenString == "" {
		return &PageTokenChaim{
			PageOffset: 0,
		}, nil
	}

	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &PageTokenChaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Properties.BufMan.PageTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*PageTokenChaim); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid page token")
	}
}
