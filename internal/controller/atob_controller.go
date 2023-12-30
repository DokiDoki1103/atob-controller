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
	"github.com/DokiDoki1103/atob-controller/internal/docker"
	"github.com/DokiDoki1103/atob-controller/internal/gitx"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	atobv1 "github.com/DokiDoki1103/atob-controller/api/v1"
)

// AtobReconciler reconciles a Atob object
type AtobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=atob.itihey.com,resources=atobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atob.itihey.com,resources=atobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atob.itihey.com,resources=atobs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Atob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *AtobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	atob := new(atobv1.Atob)
	if err := r.Get(ctx, req.NamespacedName, atob); err != nil {
		log.FromContext(ctx).Error(err, "unable to fetch Atob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if atob.Status.Status == atobv1.Success {
		return ctrl.Result{}, nil
	}
	log.FromContext(ctx).Info(fmt.Sprintf("Reconciling Atob: %s", atob.Name))

	codePath := "data/code/" + gitx.GetRepoName(atob.Spec.Git.Url)
	logPath := "data/log/1.log"

	file, _ := os.Create(logPath)
	defer file.Close()

	r.UpdateStatus(ctx, atob, atobv1.Pulling)
	err := gitx.Default().PullOrClone(ctx, codePath, atob.Spec.Git, file)
	if err != nil {
		log.FromContext(ctx).Error(err, "unable to pull or clone git")
		r.UpdateStatus(ctx, atob, atobv1.PullFailed)
		return ctrl.Result{}, err
	}
	r.UpdateStatus(ctx, atob, atobv1.PullSuccess)
	err = docker.Default().Build(ctx, codePath, atob.Spec.Image, file)
	if err != nil {
		r.UpdateStatus(ctx, atob, atobv1.BuildFailed)
		return ctrl.Result{}, err
	}
	r.UpdateStatus(ctx, atob, atobv1.Success)
	log.FromContext(ctx).Info("success")
	return ctrl.Result{}, nil
}

func (r *AtobReconciler) UpdateStatus(ctx context.Context, atob *atobv1.Atob, status string) {
	logger := log.FromContext(ctx)
	atob.Status.Status = status
	if err := r.Status().Update(ctx, atob); err != nil {
		logger.Error(err, "unable to update Atob status: "+status)
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&atobv1.Atob{}).
		Complete(r)
}
