package services

import (
	"context"
	"errors"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TokenService interface {
	CreateToken(ctx context.Context, userName, password string, expireTime time.Time, note string) (*model.Token, e.ResponseError)
	GetToken(ctx context.Context, userID, tokenID string) (*model.Token, e.ResponseError)
	ListTokens(ctx context.Context, userID string, offset, limit int, reverse bool) (model.Tokens, e.ResponseError)
	DeleteToken(ctx context.Context, userID, tokenID string) e.ResponseError
}

type TokenServiceImpl struct {
	userMapper  mapper.UserMapper
	tokenMapper mapper.TokenMapper
}

func NewTokenService() TokenService {
	return &TokenServiceImpl{
		userMapper:  &mapper.UserMapperImpl{},
		tokenMapper: &mapper.TokenMapperImpl{},
	}
}

func (tokenService *TokenServiceImpl) CreateToken(ctx context.Context, userName, password string, expireTime time.Time, note string) (*model.Token, e.ResponseError) {
	user, err := tokenService.userMapper.FindByUserName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewPermissionDeniedError(registryv1alpha1connect.TokenServiceCreateTokenProcedure)
		}

		return nil, e.NewInternalError(registryv1alpha1connect.TokenServiceCreateTokenProcedure)
	}
	if security.EncryptPlainPassword(userName, password) != user.Password {
		// 密码不正确
		return nil, e.NewPermissionDeniedError(registryv1alpha1connect.TokenServiceCreateTokenProcedure)
	}

	token := &model.Token{
		ID:         0,
		UserID:     user.UserID,
		TokenID:    uuid.NewString(),
		TokenName:  security.GenerateToken(userName, note),
		ExpireTime: expireTime,
		Note:       note,
	}
	err = tokenService.tokenMapper.Create(token) // 创建token
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.TokenServiceCreateTokenProcedure)
	}

	return token, nil
}

func (tokenService *TokenServiceImpl) GetToken(ctx context.Context, userID, tokenID string) (*model.Token, e.ResponseError) {
	token, err := tokenService.tokenMapper.FindAvailableByTokenID(tokenID)
	if err != nil {
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("token")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.TokenServiceGetTokenProcedure)
	}
	if userID != token.UserID {
		// 不能查看其他人的token
		return nil, e.NewPermissionDeniedError(registryv1alpha1connect.TokenServiceGetTokenProcedure)
	}

	return token, nil
}

func (tokenService *TokenServiceImpl) ListTokens(ctx context.Context, userID string, offset, limit int, reverse bool) (model.Tokens, e.ResponseError) {
	tokens, err := tokenService.tokenMapper.FindAvailablePageByUserID(userID, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.TokenServiceListTokensProcedure)
	}

	return tokens, nil
}

func (tokenService *TokenServiceImpl) DeleteToken(ctx context.Context, userID, tokenID string) e.ResponseError {
	// 查询token
	token, err := tokenService.tokenMapper.FindAvailableByTokenID(tokenID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewNotFoundError("token")
		}

		return e.NewInternalError(registryv1alpha1connect.TokenServiceDeleteTokenProcedure)
	}
	if token.UserID != userID {
		// 不能删除其他人的token
		return e.NewPermissionDeniedError(registryv1alpha1connect.TokenServiceDeleteTokenProcedure)
	}

	// 删除token
	err = tokenService.tokenMapper.DeleteByTokenID(tokenID)
	if err != nil {
		return e.NewInternalError(registryv1alpha1connect.TokenServiceDeleteTokenProcedure)
	}

	return nil
}
