/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AtobSpec defines the desired state of Atob
type AtobSpec struct {
	// Git 源仓库配置
	Git *GitConfig `json:"git,omitempty" yaml:"git,omitempty"`
	// Image 最后要生成的镜像配置
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`
}

// GitConfig 源仓库配置
type GitConfig struct {
	// Url 仓库地址
	Url string `json:"url" yaml:"url"`
	// Username 仓库认证用户名
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	// Password 仓库认证密码
	Password string `json:"password,omitempty" yaml:"password,omitempty"`

	// Branch 分支
	Branch string `json:"branch,omitempty" yaml:"branch,omitempty"`
	// Tag 标签
	Tag string `json:"tag,omitempty" yaml:"tag,omitempty"`
}

// ImageConfig 最后要生成的镜像配置
type ImageConfig struct {
	// Name 镜像名称
	Name string `json:"name" yaml:"name"`
	// ImageTag 镜像版本
	Tag string `json:"tag" yaml:"tag"`

	// Registry 镜像仓库地址
	Registry string `json:"registry,omitempty" yaml:"registry,omitempty"`
	// Username 镜像仓库认证用户名
	Username string `json:"username,omitempty" yaml:"username,omitempty"`

	// Password 镜像仓库认证密码
	Password string `json:"password,omitempty" yaml:"password,omitempty"`

	// Push 是否自动推送该镜像
	Push bool `json:"push,omitempty" yaml:"push,omitempty"`

	// Build 是否自动build该镜像
	Build bool `json:"build,omitempty" yaml:"build,omitempty"`
}

const (
	Error            = "Error"            // 错误
	Pulling          = "Pulling"          // 正在拉取源码
	PullSuccess      = "PullSuccess"      // 拉取源码成功
	PullFailed       = "PullFailed"       // 拉取源码失败
	CheckCodeSuccess = "CheckCodeSuccess" // 检查源码成功
	Building         = "Building"         // 正在打包构建
	BuildSuccess     = "BuildSuccess"     // 打包构建成功
	BuildFailed      = "BuildFailed"      // 打包构建失败
	Success          = "Success"          // 成功
)

// AtobStatus defines the observed state of Atob
type AtobStatus struct {
	// LogPath 指定任务日志路径
	LogPath string `json:"logPath,omitempty" yaml:"logPath,omitempty"`

	// Status 任务状态
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Atob is the Schema for the atobs API
type Atob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtobSpec   `json:"spec,omitempty"`
	Status AtobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AtobList contains a list of Atob
type AtobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Atob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Atob{}, &AtobList{})
}
