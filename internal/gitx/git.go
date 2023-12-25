package gitx

import (
	"context"
	atobv1 "github.com/DokiDoki1103/atob-controller/api/v1"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"os"
)

type GitConfig struct {
	Endpoint      *transport.Endpoint `json:"endpoint"`
	Url           string              `json:"url"`
	SourceCodeDir string              `json:"source_code_dir"`
	Username      string              `json:"username"`
	Password      string              `json:"password"`
	LogPath       string              `json:"log_path"`
}
type GitInterface interface {
	Clone(ctx context.Context, path string, config atobv1.GitConfig, file *os.File) error
	Pull(ctx context.Context, path string, config atobv1.GitConfig, file *os.File) error

	PullOrClone(ctx context.Context, path string, config atobv1.GitConfig, file *os.File) error
}

type gitx struct {
}

var defaultGit *gitx

func NewGit() {
	defaultGit = &gitx{}
}

func Default() GitInterface {
	return defaultGit
}
