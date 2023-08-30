package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/interceptors"
	"github.com/bufbuild/connect-go"
	"net/http"
	"testing"
)

func TestSearchServiceHandler(t *testing.T) {
	TestCreateToken(t)
	searchClient := registryv1alpha1connect.NewSearchServiceClient(http.DefaultClient, "http://bufman.io", interceptors.WithAuthHeaderInterceptor("626b2ea6acd3214e"))

	resp1, err := searchClient.SearchTag(context.Background(), connect.NewRequest(&registryv1alpha1.SearchTagRequest{
		RepositoryOwner: "dawnzzz",
		RepositoryName:  "storeapis",
		Query:           "v1",
		PageSize:        10,
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(resp1.Msg.String())

	resp2, err := searchClient.SearchUser(context.Background(), connect.NewRequest(&registryv1alpha1.SearchUserRequest{
		Query:    "test",
		PageSize: 10,
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(resp2.Msg.String())

	resp3, err := searchClient.SearchCurationPlugin(context.Background(), connect.NewRequest(&registryv1alpha1.SearchCuratedPluginRequest{
		Query:    "go",
		PageSize: 10,
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(resp3.Msg.String())

	resp4, err := searchClient.SearchRepository(context.Background(), connect.NewRequest(&registryv1alpha1.SearchRepositoryRequest{
		Query:    "storeapis",
		PageSize: 10,
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(resp4.Msg.String())

	resp5, err := searchClient.SearchLastCommitByContent(context.Background(), connect.NewRequest(&registryv1alpha1.SearchLastCommitByContentRequest{
		Query:    "PhoneType Pet",
		PageSize: 10,
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(resp5.Msg.String())
}
