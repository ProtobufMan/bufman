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

var (
	testRepositoryName     = "testRepo"
	testRepositoryFullName = testUsername + "/" + testRepositoryName
	testRepository         *registryv1alpha1.Repository
)

func TestCreateRepositoryByFullName(t *testing.T) {
	defer TestDeleteToken(t)

	client := newTestRepositoryClient(t)
	req := connect.NewRequest(&registryv1alpha1.CreateRepositoryByFullNameRequest{
		FullName:   testRepositoryFullName,
		Visibility: registryv1alpha1.Visibility_VISIBILITY_PUBLIC,
	})
	resp, err := client.CreateRepositoryByFullName(context.Background(), req)
	if err != nil && connect.CodeOf(err) != connect.CodeAlreadyExists {
		t.Log(connect.CodeOf(err))
		t.Error("create repo by full name error")
		return
	}
	if connect.CodeOf(err) == connect.CodeAlreadyExists {
		t.Log("repo already exists")
		return
	}

	t.Log(resp.Msg.String())
}

func TestGetRepositoryByFullName(t *testing.T) {
	defer TestDeleteToken(t)
	TestCreateRepositoryByFullName(t)

	client := newTestRepositoryClient(t)

	req := connect.NewRequest(&registryv1alpha1.GetRepositoryByFullNameRequest{
		FullName: testRepositoryFullName,
	})
	resp, err := client.GetRepositoryByFullName(context.Background(), req)
	if err != nil {
		t.Error("get repo error", connect.CodeOf(err))
		return
	}

	testRepository = resp.Msg.GetRepository()
	t.Log(resp.Msg.String())
}

func TestGetRepository(t *testing.T) {
	defer TestDeleteToken(t)

	TestGetRepositoryByFullName(t)

	client := newTestRepositoryClient(t)

	req := connect.NewRequest(&registryv1alpha1.GetRepositoryRequest{
		Id: testRepository.GetId(),
	})
	resp, err := client.GetRepository(context.Background(), req)
	if err != nil {
		t.Error("get repo error")
	}

	testRepository = resp.Msg.GetRepository()
}

func TestListRepositories(t *testing.T) {
	defer TestDeleteToken(t)

	TestGetRepositoryByFullName(t)

	client := newTestRepositoryClient(t)

	req := connect.NewRequest(&registryv1alpha1.ListRepositoriesRequest{
		PageSize: 10,
		Reverse:  false,
	})
	_, err := client.ListRepositories(context.Background(), req)
	if err != nil {
		t.Error("list repo error")
	}
}

func TestUpdateRepositorySettingsByName(t *testing.T) {
	defer TestDeleteToken(t)

	TestGetRepositoryByFullName(t)
	client := newTestRepositoryClient(t)

	desc := "this is a test repo"
	req := connect.NewRequest(&registryv1alpha1.UpdateRepositorySettingsByNameRequest{
		OwnerName:      testUsername,
		RepositoryName: testRepositoryName,
		Visibility:     registryv1alpha1.Visibility_VISIBILITY_PRIVATE,
		Description:    &desc,
	})
	_, err := client.UpdateRepositorySettingsByName(context.Background(), req)
	if err != nil {
		t.Log(connect.CodeOf(err))
		t.Error("list repo error")
	}
}

func TestDeleteRepositoryByFullName(t *testing.T) {
	defer TestDeleteToken(t)

	TestGetRepositoryByFullName(t)
	client := newTestRepositoryClient(t)

	req := connect.NewRequest(&registryv1alpha1.DeleteRepositoryByFullNameRequest{
		FullName: testRepositoryFullName,
	})
	_, err := client.DeleteRepositoryByFullName(context.Background(), req)
	if err != nil {
		t.Log(connect.CodeOf(err))
		t.Error("list repo error")
	}
}

func newTestRepositoryClient(t *testing.T) registryv1alpha1connect.RepositoryServiceClient {
	TestCreateToken(t)
	var client registryv1alpha1connect.RepositoryServiceClient
	client = registryv1alpha1connect.NewRepositoryServiceClient(http.DefaultClient, "http://localhost:39099", interceptors.WithAuthHeaderInterceptor(testToken))

	return client
}
