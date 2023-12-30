package docker

import (
	"context"
	atobv1 "github.com/DokiDoki1103/atob-controller/api/v1"
	"github.com/docker/docker/client"
	"io"
	"os"
)

type DockerInterface interface {
	Build(ctx context.Context, path string, config *atobv1.ImageConfig, file *os.File) error
	Pull(ctx context.Context, config *atobv1.ImageConfig) (io.ReadCloser, error)
}

var defaultDocker *docker

type docker struct {
	dockerClient *client.Client
}

func (d *docker) Start() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	d.dockerClient = cli
	return nil
}

func (d *docker) CloseHandle() {
	d.dockerClient.Close()
}

func NewDocker() {
	// 递归创建目录 用于存放日志
	if err := os.MkdirAll("data/log/", os.ModePerm); err != nil {
	}

	// 递归创建目录 用于存放源码
	if err := os.MkdirAll("data/code/", os.ModePerm); err != nil {
	}

	defaultDocker = &docker{}
	defaultDocker.Start()
}

func Default() DockerInterface {
	return defaultDocker
}
