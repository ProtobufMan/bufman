package interceptors

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/bufbuild/connect-go"
	"net/http"
	"strings"
)

type OptionalAuthInterceptor struct {
	tokenMapper mapper.TokenMapper
}

var optionalAuthProcedures = make(map[string]bool)

func WithOptionalAuthInterceptor(procedures ...string) connect.Option {
	for i := 0; i < len(procedures); i++ {
		optionalAuthProcedures[procedures[i]] = true
	}

	return connect.WithInterceptors(connect.UnaryInterceptorFunc(OptionalAuthInterceptor{&mapper.TokenMapperImpl{}}.WrapUnary))
}

func (optionalAuthInterceptor OptionalAuthInterceptor) WrapUnary(unaryFunc connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if !needOptionalAuth(req.Spec().Procedure) {
			// 不需要验证token
			return unaryFunc(ctx, req)
		}

		// 选择性的鉴权
		userID := optionalAuthInterceptor.optionalAuth(req.Header())

		// 如果没有携带token，或者token验证未通过，则userID = ""
		resp, err := unaryFunc(context.WithValue(ctx, constant.UserIDKey, userID), req)

		return resp, err
	}
}

func (optionalAuthInterceptor OptionalAuthInterceptor) optionalAuth(header http.Header) (userID string) {
	raw := header.Get(constant.AuthHeader)
	contents := strings.Split(raw, " ")
	if len(contents) != 2 {
		return ""
	}

	prefix, token := contents[0], contents[1]
	if prefix != constant.AuthPrefix {
		return ""
	}

	// 验证token
	tokenEntity, err := optionalAuthInterceptor.tokenMapper.FindAvailableByTokenName(token)
	if err != nil {
		return ""
	}

	return tokenEntity.UserID
}

func needOptionalAuth(procedure string) bool {
	return authProcedures[procedure]
}
