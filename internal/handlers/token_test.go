package handlers

import (
	"context"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/interceptors"
	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"testing"
	"time"
)

var (
	testToken  string
	testTokens []*registryv1alpha.Token
)

func TestCreateToken(t *testing.T) {
	var tokenClient registryv1alphaconnect.TokenServiceClient
	tokenClient = registryv1alphaconnect.NewTokenServiceClient(http.DefaultClient, "http://localhost:39099")

	TestCreateUser(t)

	// CreateToken
	createCreateTokenReq := connect.NewRequest(&registryv1alpha.CreateTokenRequest{
		Note:       "test note",
		ExpireTime: timestamppb.New(time.Now().Add(24 * time.Hour)),
		Username:   testUsername,
		Password:   testPassword,
	})
	var createTokenResp *connect.Response[registryv1alpha.CreateTokenResponse]
	var err error
	createTokenResp, err = tokenClient.CreateToken(context.Background(), createCreateTokenReq)
	if err != nil {
		t.Error("create token error")
		return
	}

	testToken = createTokenResp.Msg.GetToken()
}

func TestCheckAuth(t *testing.T) {
	tokenClient := registryv1alphaconnect.NewTokenServiceClient(http.DefaultClient, "http://localhost:39099", interceptors.WithAuthHeaderInterceptor("wrong token"))
	ListTokensReq := connect.NewRequest(&registryv1alpha.ListTokensRequest{
		PageSize:   10,
		PageOffset: 0,
		Reverse:    false,
	})
	_, err := tokenClient.ListTokens(context.Background(), ListTokensReq)
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Error("auth check error")
		return
	}
}

func TestListTokens(t *testing.T) {
	TestCreateToken(t)
	// ListTokens
	var tokenClient registryv1alphaconnect.TokenServiceClient
	tokenClient = registryv1alphaconnect.NewTokenServiceClient(http.DefaultClient, "http://localhost:39099", interceptors.WithAuthHeaderInterceptor(testToken))

	ListTokensReq := connect.NewRequest(&registryv1alpha.ListTokensRequest{
		PageSize:   10,
		PageOffset: 0,
		Reverse:    false,
	})
	ListTokensReq = connect.NewRequest(&registryv1alpha.ListTokensRequest{
		PageSize:   10,
		PageOffset: 0,
		Reverse:    false,
	})
	listTokensResp, err := tokenClient.ListTokens(context.Background(), ListTokensReq)
	if err != nil {
		t.Error("list tokens error")
		return
	}
	for i := 0; i < len(listTokensResp.Msg.GetTokens()); i++ {
		t.Log(listTokensResp.Msg.GetTokens()[i].String())
	}
	testTokens = listTokensResp.Msg.GetTokens()
}

func TestGetToken(t *testing.T) {
	TestCreateToken(t)
	TestListTokens(t)

	var tokenClient registryv1alphaconnect.TokenServiceClient
	tokenClient = registryv1alphaconnect.NewTokenServiceClient(http.DefaultClient, "http://localhost:39099", interceptors.WithAuthHeaderInterceptor(testToken))

	for i := 0; i < len(testTokens); i++ {
		// get token
		getTokenReq := connect.NewRequest(&registryv1alpha.GetTokenRequest{
			TokenId: testTokens[i].GetId(),
		})
		_, err := tokenClient.GetToken(context.Background(), getTokenReq)
		if err != nil {
			t.Error("get token error", err.Error())
			return
		}
	}
}

func TestDeleteToken(t *testing.T) {
	TestCreateToken(t)
	TestListTokens(t)

	var tokenClient registryv1alphaconnect.TokenServiceClient
	tokenClient = registryv1alphaconnect.NewTokenServiceClient(http.DefaultClient, "http://localhost:39099", interceptors.WithAuthHeaderInterceptor(testToken))

	for i := 0; i < len(testTokens); i++ {
		deleteTokenReq := connect.NewRequest(&registryv1alpha.DeleteTokenRequest{
			TokenId: testTokens[i].GetId(),
		})
		_, err := tokenClient.DeleteToken(context.Background(), deleteTokenReq)
		if err != nil {
			t.Error("delete token error")
			return
		}
	}
}
