package services

import (
	"errors"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(userName, password string) (*model.User, e.ResponseError)
	GetUser(userID string) (*model.User, e.ResponseError)
	GetUserByUsername(userName string) (*model.User, e.ResponseError)
	ListUsers(offset int, limit int, reverse bool) (model.Users, e.ResponseError)
}

type UserServiceImpl struct {
	userMapper mapper.UserMapper
}

func NewUserService() UserService {
	return &UserServiceImpl{
		userMapper: &mapper.UserMapperImpl{},
	}
}

func (userService *UserServiceImpl) CreateUser(userName, password string) (*model.User, e.ResponseError) {
	user := &model.User{
		UserID:   uuid.NewString(),
		UserName: userName,
		Password: util.EncryptPlainPassword(userName, password), // 加密明文密码
	}

	err := userService.userMapper.Create(user) // 创建用户
	if err != nil {
		// 用户重复
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, e.NewAlreadyExistsError("user name")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.UserServiceCreateUserProcedure)
	}

	return user, nil
}

func (userService *UserServiceImpl) GetUser(userID string) (*model.User, e.ResponseError) {
	user, err := userService.userMapper.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("user")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.UserServiceGetUserProcedure)
	}

	return user, nil
}

func (userService *UserServiceImpl) GetUserByUsername(userName string) (*model.User, e.ResponseError) {
	user, err := userService.userMapper.FindByUserName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("user")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.UserServiceGetUserByUsernameProcedure)
	}

	return user, nil
}

func (userService *UserServiceImpl) ListUsers(offset int, limit int, reverse bool) (model.Users, e.ResponseError) {
	users, err := userService.userMapper.FindPage(offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.UserServiceListUsersProcedure)
	}

	return users, nil
}
