package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	atobv1 "github.com/DokiDoki1103/atob-controller/api/v1"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"io"
	"os"
)

func (d docker) Build(ctx context.Context, path string, config atobv1.ImageConfig, file *os.File) error {
	tar, err := archive.TarWithOptions(path, &archive.TarOptions{})
	if err != nil {
		return err
	}
	// 构建 Docker 镜像
	res, err := d.dockerClient.ImageBuild(ctx, tar, types.ImageBuildOptions{
		NoCache:    false,
		Tags:       []string{config.Name + ":" + config.Tag},
		Dockerfile: "Dockerfile",
	})
	defer res.Body.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func encodeAuthToBase64(authConfig types.AuthConfig) (string, error) {
	authJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(authJSON), nil
}

func (d docker) Pull(ctx context.Context, config atobv1.ImageConfig) (io.ReadCloser, error) {
	authConfig := types.AuthConfig{
		Username: config.Username,
		Password: config.Password,
	}
	encodedAuth, _ := encodeAuthToBase64(authConfig)

	pullOptions := types.ImagePullOptions{
		RegistryAuth: encodedAuth,
	}

	imageName := fmt.Sprintf("%s:%s", config.Name, config.Tag)
	buf, err := d.dockerClient.ImagePull(ctx, imageName, pullOptions)
	if err != nil {
		return nil, err
	}

	ins, _, err := d.dockerClient.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return nil, err
	}

	exportPorts := make(map[string]struct{})
	for port := range ins.Config.ExposedPorts {
		exportPorts[string(port)] = struct{}{}
	}

	marshal, err := json.Marshal(ins.Config)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(marshal))
	return buf, nil
}
