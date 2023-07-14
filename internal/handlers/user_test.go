package handlers

import (
	"context"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/bufman/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/gen/bufman/registry/v1alpha/registryv1alphaconnect"
	"github.com/bufbuild/connect-go"
	"net/http"
	"testing"
)

var (
	userClient   registryv1alphaconnect.UserServiceClient
	testUsername = "test1"
	testPassword = "123456"
)

func init() {
	userClient = registryv1alphaconnect.NewUserServiceClient(http.DefaultClient, "http://localhost:39099")
}

func TestCreateUser(t *testing.T) {
	req := connect.NewRequest(&registryv1alpha.CreateUserRequest{
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
	req1 := connect.NewRequest(&registryv1alpha.GetUserByUsernameRequest{
		Username: testUsername,
	})

	resp1, err := userClient.GetUserByUsername(context.Background(), req1)
	if err != nil {
		t.Error("get user by user name failed")
		return
	}
	t.Logf("%#v", resp1.Msg.String())

	userID := resp1.Msg.GetUser().GetId()
	req2 := connect.NewRequest(&registryv1alpha.GetUserRequest{
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
	req := connect.NewRequest(&registryv1alpha.ListUsersRequest{
		PageSize:   10,
		PageOffset: 0,
		Reverse:    false,
	})

	resp, err := userClient.ListUsers(context.Background(), req)
	if err != nil {
		t.Error("list users failed")
		return
	}
	t.Logf("%#v", resp.Msg.String())
}
