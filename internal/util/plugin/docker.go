package plugin

import (
	"bytes"
	"context"
	"errors"
	"github.com/ProtobufMan/bufman-cli/private/pkg/protoencoding"
	"github.com/ProtobufMan/bufman/internal/config"
	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"google.golang.org/protobuf/types/pluginpb"
	"io"
	"os"
)

type Docker interface {
	Close() error
	PullImage(ctx context.Context, refStr string) error
	GenerateCode(ctx context.Context, pluginName, refStr string, request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error)
}

type docker struct {
	cli           *client.Client
	serverAddress string
	username      string
	password      string
}

func NewDocker(serverAddress, username, password string) (Docker, error) {
	cli, err := config.NewDockerClient()
	if err != nil {
		return nil, err
	}

	return &docker{
		cli:           cli,
		serverAddress: serverAddress,
		username:      username,
		password:      password,
	}, nil
}

func (d *docker) Close() error {
	if d.cli != nil {
		_ = d.cli.Close()
	}

	return nil
}

func (d *docker) GenerateCode(ctx context.Context, pluginName, refStr string, request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	// 查询镜像是否已经拉取
	imageList, err := d.ListImage(ctx, refStr)
	if err != nil {
		return nil, err
	}

	if len(imageList) == 0 {
		// 拉取镜像
		err = d.PullImage(ctx, refStr)
		if err != nil {
			return nil, err
		}
	}

	// 查询容器
	var containerID string
	containerList, err := d.ListContainer(ctx, refStr)
	if err != nil {
		return nil, err
	}
	if len(containerList) == 0 {
		// 创造container
		containerID, err = d.CreateContainer(ctx, pluginName, refStr)
		if err != nil {
			return nil, err
		}
	} else {
		// 已有容器
		containerID = containerList[0].ID
	}

	// 接管输入输出
	hijackedResponse, err := d.AttachContainer(ctx, containerID)
	if err != nil {
		return nil, err
	}
	defer hijackedResponse.Close()

	// 启动容器，生成代码
	response, err := d.startContainer(ctx, containerID, hijackedResponse, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *docker) ListImage(ctx context.Context, refStr string) ([]types.ImageSummary, error) {
	return d.cli.ImageList(ctx, types.ImageListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("reference", refStr)),
	})
}

func (d *docker) ListContainer(ctx context.Context, refStr string) ([]types.Container, error) {
	containerList, err := d.cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("ancestor", refStr), filters.Arg("status", "created")),
	})
	if err != nil {
		return nil, err
	}

	if len(containerList) > 0 {
		return containerList, nil
	}

	return d.cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("ancestor", refStr), filters.Arg("status", "exited")),
	})
}

func (d *docker) PullImage(ctx context.Context, refStr string) error {
	responseBody, err := d.cli.ImagePull(ctx, refStr, types.ImagePullOptions{
		PrivilegeFunc: func() (string, error) {
			authenticateOKBody, err := d.cli.RegistryLogin(ctx, registry.AuthConfig{
				Username:      d.username,
				Password:      d.password,
				ServerAddress: d.serverAddress,
			})
			if err != nil {
				return "", err
			}

			return authenticateOKBody.IdentityToken, nil
		},
	})
	if err != nil {
		return nil
	}

	out := streams.NewOut(os.Stdout)
	err = jsonmessage.DisplayJSONMessagesToStream(responseBody, out, nil)
	if err != nil {
		return err
	}

	return nil
}

func (d *docker) CreateContainer(ctx context.Context, pluginName, image string) (string, error) {
	createResponse, err := d.cli.ContainerCreate(ctx, &container.Config{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		OpenStdin:    true,
		StdinOnce:    true,
		Image:        image,
		Entrypoint:   []string{pluginName},
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	return createResponse.ID, nil
}

func (d *docker) AttachContainer(ctx context.Context, containerID string) (types.HijackedResponse, error) {
	return d.cli.ContainerAttach(ctx, containerID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		// DetachKeys: containerID,
	})
}

func (d *docker) startContainer(ctx context.Context, containerID string, hijackedResponse types.HijackedResponse, request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	// 读取CodeGeneratorRequest作为输入
	requestData, err := protoencoding.NewWireMarshaler().Marshal(request)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(hijackedResponse.Conn, bytes.NewReader(requestData))
	if err != nil {
		return nil, err
	}
	_ = hijackedResponse.CloseWrite()

	// 执行
	err = d.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	// 读取输出
	responseBuffer := bytes.NewBuffer(nil)
	errBuffer := bytes.NewBuffer(nil)

	_, err = stdcopy.StdCopy(responseBuffer, errBuffer, hijackedResponse.Reader)
	if err != nil {
		return nil, err
	}

	if errBuffer.Len() != 0 {
		return nil, errors.New(errBuffer.String())
	}

	// 转换为CodeGeneratorResponse
	response := &pluginpb.CodeGeneratorResponse{}
	err = protoencoding.NewWireUnmarshaler(nil).Unmarshal(responseBuffer.Bytes(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
