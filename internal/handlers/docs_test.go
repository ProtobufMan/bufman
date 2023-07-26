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

func TestDocServiceHandler(t *testing.T) {
	TestCreateToken(t)
	docClient := registryv1alpha1connect.NewDocServiceClient(http.DefaultClient, "http://bufman.io", interceptors.WithAuthHeaderInterceptor("50251bb640c7bc5c"))

	resp1, err := docClient.GetSourceDirectoryInfo(context.Background(), connect.NewRequest(&registryv1alpha1.GetSourceDirectoryInfoRequest{
		Owner:      "dawnzzz",
		Repository: "storeapis",
		Reference:  "",
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(resp1.Msg.String())

	resp2, err := docClient.GetModulePackages(context.Background(), connect.NewRequest(&registryv1alpha1.GetModulePackagesRequest{
		Owner:      "dawnzzz",
		Repository: "storeapis",
		Reference:  "",
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(resp2.Msg.String())

	resp3, err := docClient.GetPackageDocumentation(context.Background(), connect.NewRequest(&registryv1alpha1.GetPackageDocumentationRequest{
		Owner:       "dawnzzz",
		Repository:  "storeapis",
		Reference:   "",
		PackageName: "store.v1",
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(resp3.Msg.String())

	resp4, err := docClient.GetModuleDocumentation(context.Background(), connect.NewRequest(&registryv1alpha1.GetModuleDocumentationRequest{
		Owner:      "dawnzzz",
		Repository: "storeapis",
		Reference:  "",
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(resp4.Msg.String())

	resp5, err := docClient.GetSourceFile(context.Background(), connect.NewRequest(&registryv1alpha1.GetSourceFileRequest{
		Owner:      "dawnzzz",
		Repository: "storeapis",
		Reference:  "",
		Path:       "stroe/v1/store.proto",
	}))
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(string(resp5.Msg.GetContent()))
}
