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
	testDockerRepoName     = "testDockerRepoName"
	testDockerRepoAddress  = "input your docker repo address"
	testDockerRepoUsername = "input your docker repo username"
	testDockerRepoPassword = "input your docker repo password"
)

func TestCreateDockerRepo(t *testing.T) {
	client := newTestDockerRepoClient(t)

	req := connect.NewRequest(&registryv1alpha1.CreateDockerRepoRequest{
		Name:     testDockerRepoName,
		Address:  testDockerRepoAddress,
		Username: testDockerRepoUsername,
		Password: testDockerRepoPassword,
		Note:     "test docker repo",
	})
	resp, err := client.CreateDockerRepo(context.Background(), req)
	if err != nil && connect.CodeOf(err) != connect.CodeAlreadyExists {
		t.Log(connect.CodeOf(err))
		t.Error("create docker repo  error")
		return
	}
	if connect.CodeOf(err) == connect.CodeAlreadyExists {
		t.Log("docker repo already exists")
		return
	}

	t.Log(resp.Msg.String())
}

func newTestDockerRepoClient(t *testing.T) registryv1alpha1connect.DockerRepoServiceClient {
	TestCreateToken(t)
	var client registryv1alpha1connect.DockerRepoServiceClient
	client = registryv1alpha1connect.NewDockerRepoServiceClient(http.DefaultClient, "http://bufman.io", interceptors.WithAuthHeaderInterceptor(testToken))

	return client
}
