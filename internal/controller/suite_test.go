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

package controller

import (
	"context"
	"fmt"
	atobv1 "github.com/DokiDoki1103/atob-controller/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"testing"
	"time"
)

func TestControllers2(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(atobv1.AddToScheme(scheme))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:  scheme,
		Metrics: metricsserver.Options{BindAddress: ":8090"},
	})
	if err != nil {
		t.Fatalf("Error creating manager: %v", err)
	}

	go func() {
		// 明确启动缓存
		if err := mgr.GetCache().Start(context.Background()); err != nil {
			t.Fatalf("Error starting cache: %v", err)
		}
	}()
	time.Sleep(2 * time.Second)

	atob := &atobv1.Atob{}
	err = mgr.GetClient().Get(context.Background(), client.ObjectKey{
		Namespace: "default",
		Name:      "atob-sample",
	}, atob)
	fmt.Println(atob)
	if err != nil {
		t.Fatalf("Error getting resource: %v", err)
	}

	// 对获取到的资源进行更新
	atob.Status.Status = "Updated status"
	err = mgr.GetClient().Status().Update(context.Background(), atob)
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	// 添加断言，验证更新后的状态是否正确
	if atob.Status.Status != "Updated status" {
		t.Errorf("Unexpected status after update: %s", atob.Status)
	}
}
