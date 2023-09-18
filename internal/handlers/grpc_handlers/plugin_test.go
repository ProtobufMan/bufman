package grpc_handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/interceptors"
	"github.com/bufbuild/connect-go"
	"net/http"
	"testing"
)

func TestCreatePlugin(t *testing.T) {
	client := newTestPluginClient(t)

	req := connect.NewRequest(&registryv1alpha1.CreateCuratedPluginRequest{
		Owner:          testUsername,
		Name:           "protoc-gen-go",
		Version:        "v1.0.0",
		ImageDigest:    "sha256:ed64818305c5ea32d62f61507bd534e1b67025db9032b4716c0bd43af8d62181",
		Description:    "this is a test go plugin",
		Revision:       1,
		Visibility:     1,
		ImageName:      "zomato/protoc-gen",
		DockerRepoName: testDockerRepoName,
	})
	resp, err := client.CreateCuratedPlugin(context.Background(), req)
	if err != nil && connect.CodeOf(err) != connect.CodeAlreadyExists {
		t.Log(connect.CodeOf(err))
		t.Error("create plugin  error")
		return
	}
	if connect.CodeOf(err) == connect.CodeAlreadyExists {
		t.Log("plugin already exists")
		return
	}

	t.Log(resp.Msg.String())
}

func newTestPluginClient(t *testing.T) registryv1alpha1connect.PluginCurationServiceClient {
	TestCreateToken(t)
	var client registryv1alpha1connect.PluginCurationServiceClient
	client = registryv1alpha1connect.NewPluginCurationServiceClient(http.DefaultClient, "http://bufman.io", interceptors.WithAuthHeaderInterceptor(testToken))

	return client
}
