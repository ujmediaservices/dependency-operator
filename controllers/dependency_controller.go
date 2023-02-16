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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	http "net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/gaiaslastlaugh/dependency-operator/api/v1alpha1"
)

// DependencyReconciler reconciles a Dependency object
type DependencyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dep.example.com,resources=dependencies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dep.example.com,resources=dependencies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dep.example.com,resources=dependencies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Dependency object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *DependencyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	logger := log.Log.WithValues("dependency", req.NamespacedName)

	logger.Info("Checking dependency...")
	dep := &v1alpha1.Dependency{}
	err := r.Get(ctx, req.NamespacedName, dep)
	url := dep.Spec.Url
	resp, err := http.Get(url)
	if err != nil {
		return ctrl.Result{}, err
	} else {
		if resp.StatusCode == 200 {
			logger.Info("Verified URL", "Url", url)
			dep.Status.DependencyStatus = true
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, err
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DependencyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Dependency{}).
		Complete(r)
}
