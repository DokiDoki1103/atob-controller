package gitx

import (
	"context"
	"errors"
	atobv1 "github.com/DokiDoki1103/atob-controller/api/v1"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	cssh "golang.org/x/crypto/ssh"
	"log"

	"os"
)

func sshAuth(ctx context.Context) (transport.AuthMethod, error) {
	sshKey, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa")
	if err != nil {
		return nil, err
	}
	auth, err := ssh.NewPublicKeys("git", sshKey, "")
	if err != nil {
		return nil, err
	}
	auth.HostKeyCallbackHelper.HostKeyCallback = cssh.InsecureIgnoreHostKey()
	return auth, nil
}

func basicAuth(username, password string) transport.AuthMethod {
	return &githttp.BasicAuth{
		Username: username,
		Password: password,
	}
}

func authMethod(ctx context.Context, protocol, username, password string) (transport.AuthMethod, error) {
	if protocol == "ssh" {
		auth, err := sshAuth(ctx)
		if err != nil {
			return nil, err
		}
		return auth, nil
	} else if protocol == "https" {
		return basicAuth(username, password), nil
	} else {
		return nil, errors.New("暂时不支持的类型")
	}
}
func (g *gitx) Clone(ctx context.Context, path string, config atobv1.GitConfig, file *os.File) error {
	defer file.Close()

	opts := &git.CloneOptions{
		URL:      config.Url,
		Progress: file,
		Tags:     git.NoTags,
		Depth:    1,
	}

	auth, err := authMethod(ctx, config.Endpoint.Protocol, config.Username, config.Password)

	log.Println(auth, err)
	if err != nil {
		return err
	}
	opts.Auth = auth

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = git.PlainCloneContext(ctx, path, false, opts)

	return err
}

func (g *gitx) Pull(ctx context.Context, path string, config atobv1.GitConfig, file *os.File) error {

	repo, err := git.PlainOpen(path)

	if err != nil {
		return err
	}

	if err != nil {
		return g.Clone(ctx, path, config, file)
	}

	opts := &git.PullOptions{
		SingleBranch: true,
		Depth:        1,
		Force:        true,
		Progress:     file,
	}

	auth, err := authMethod(ctx, config.Endpoint.Protocol, config.Username, config.Password)
	log.Println(auth, err)
	if err != nil {
		return err
	}
	opts.Auth = auth

	if err != nil {
		return err
	}
	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = tree.PullContext(ctx, opts)
	if err != nil {
		if err.Error() == "already up-to-date" {
			return nil
		}
		return err
	}

	return nil
}

// PullOrClone 拉取或者克隆
func (g *gitx) PullOrClone(ctx context.Context, path string, config atobv1.GitConfig, file *os.File) error {

	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := g.Clone(ctx, path, config, file)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		err := g.Pull(ctx, path, config, file)
		if err != nil {
			return err
		}
	}
	return nil

}
