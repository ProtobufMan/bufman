package grpc_handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/bufbuild/connect-go"
	"net/http"
	"testing"
)

var (
	userClient   registryv1alpha1connect.UserServiceClient
	testUsername = "test1"
	testPassword = "123456"
)

func init() {
	userClient = registryv1alpha1connect.NewUserServiceClient(http.DefaultClient, "http://bufman.io")
}

func TestCreateUser(t *testing.T) {
	req := connect.NewRequest(&registryv1alpha1.CreateUserRequest{
		Username: testUsername,
		Password: testPassword,
	})
	resp, err := userClient.CreateUser(context.Background(), req)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Logf("%#v", resp.Msg.String())
		t.Log("CreateUser Success")
	}

	_, err = userClient.CreateUser(context.Background(), req)
	if connect.CodeOf(err) == connect.CodeAlreadyExists {
		t.Log("Test Duplicated User Name Success, resp:", err.Error())
	} else {
		t.Error("Test Duplicated User Name Failed")
	}
}

func TestGetUser(t *testing.T) {
	req1 := connect.NewRequest(&registryv1alpha1.GetUserByUsernameRequest{
		Username: testUsername,
	})

	resp1, err := userClient.GetUserByUsername(context.Background(), req1)
	if err != nil {
		t.Error("get user by user name failed")
		return
	}
	t.Logf("%#v", resp1.Msg.String())

	userID := resp1.Msg.GetUser().GetId()
	req2 := connect.NewRequest(&registryv1alpha1.GetUserRequest{
		Id: userID,
	})
	resp2, err := userClient.GetUser(context.Background(), req2)
	if err != nil {
		t.Error("get user failed")
		return
	}

	t.Logf("%#v", resp2.Msg.String())
}

func TestListUsers(t *testing.T) {
	req := connect.NewRequest(&registryv1alpha1.ListUsersRequest{
		PageSize: 10,
		Reverse:  false,
	})

	resp, err := userClient.ListUsers(context.Background(), req)
	if err != nil {
		t.Error("list users failed")
		return
	}
	t.Logf("%#v", resp.Msg.String())
}
