package interceptors

import (
	"context"
	"errors"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/bufbuild/connect-go"
	"net/http"
	"strings"
)

type AuthInterceptor struct {
	tokenMapper mapper.TokenMapper
}

var authProcedures = make(map[string]bool)

func WithAuthInterceptor(procedures ...string) connect.Option {
	for i := 0; i < len(procedures); i++ {
		authProcedures[procedures[i]] = true
	}

	return connect.WithInterceptors(connect.UnaryInterceptorFunc(AuthInterceptor{&mapper.TokenMapperImpl{}}.WrapUnary))
}

func (authInterceptor AuthInterceptor) WrapUnary(unaryFunc connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if !needAuth(req.Spec().Procedure) {
			// 不需要验证token
			return unaryFunc(ctx, req)
		}

		userID, authed := authInterceptor.auth(req.Header())
		if !authed {
			// 验证未通过
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
		}

		// 验证通过
		resp, err := unaryFunc(context.WithValue(ctx, constant.UserIDKey, userID), req)

		return resp, err
	}
}

func (authInterceptor AuthInterceptor) auth(header http.Header) (userID string, ok bool) {
	raw := header.Get(constant.AuthHeader)
	contents := strings.Split(raw, " ")
	if len(contents) != 2 {
		return "", false
	}

	prefix, token := contents[0], contents[1]
	if prefix != constant.AuthPrefix {
		return "", false
	}

	// 验证token
	tokenEntity, err := authInterceptor.tokenMapper.FindAvailableByTokenName(token)
	if err != nil {
		return "", false
	}

	return tokenEntity.UserID, true
}

func needAuth(procedure string) bool {
	return authProcedures[procedure]
}

type AuthHeaderInterceptor struct {
	token string
}

func WithAuthHeaderInterceptor(token string) connect.Option {
	return connect.WithInterceptors(connect.UnaryInterceptorFunc(AuthHeaderInterceptor{token: token}.WrapUnary))
}

func (authHeaderInterceptor AuthHeaderInterceptor) WrapUnary(unaryFunc connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		request.Header().Set(constant.AuthHeader, constant.AuthPrefix+" "+authHeaderInterceptor.token)

		resp, err := unaryFunc(ctx, request)

		return resp, err
	}
}
