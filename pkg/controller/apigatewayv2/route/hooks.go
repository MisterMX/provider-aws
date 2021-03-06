/*
Copyright 2020 The Crossplane Authors.

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

package route

import (
	"context"

	svcsdk "github.com/aws/aws-sdk-go/service/apigatewayv2"
	ctrl "sigs.k8s.io/controller-runtime"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane/provider-aws/apis/apigatewayv2/v1alpha1"
	aws "github.com/crossplane/provider-aws/pkg/clients"
)

// SetupRoute adds a controller that reconciles Route.
func SetupRoute(mgr ctrl.Manager, l logging.Logger) error {
	name := managed.ControllerName(svcapitypes.RouteGroupKind)
	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&svcapitypes.Route{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(svcapitypes.RouteGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient()}),
			managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient())),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

func (*external) preObserve(context.Context, *svcapitypes.Route) error {
	return nil
}
func (*external) postObserve(_ context.Context, cr *svcapitypes.Route, _ *svcsdk.GetRouteOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	cr.SetConditions(xpv1.Available())
	return obs, nil
}

func (*external) preCreate(context.Context, *svcapitypes.Route) error {
	return nil
}

func (e *external) postCreate(_ context.Context, cr *svcapitypes.Route, res *svcsdk.CreateRouteOutput, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	// NOTE(muvaf): Route ID is chosen as external name since it's the only unique
	// identifier.
	meta.SetExternalName(cr, aws.StringValue(res.RouteId))
	cre.ExternalNameAssigned = true
	return cre, nil
}

func (*external) preUpdate(context.Context, *svcapitypes.Route) error {
	return nil
}

func (*external) postUpdate(_ context.Context, _ *svcapitypes.Route, upd managed.ExternalUpdate, err error) (managed.ExternalUpdate, error) {
	return upd, err
}
func lateInitialize(*svcapitypes.RouteParameters, *svcsdk.GetRouteOutput) error {
	return nil
}

func preGenerateGetRouteInput(_ *svcapitypes.Route, obj *svcsdk.GetRouteInput) *svcsdk.GetRouteInput {
	return obj
}

func postGenerateGetRouteInput(cr *svcapitypes.Route, obj *svcsdk.GetRouteInput) *svcsdk.GetRouteInput {
	obj.ApiId = cr.Spec.ForProvider.APIID
	obj.RouteId = aws.String(meta.GetExternalName(cr))
	return obj
}

func preGenerateCreateRouteInput(_ *svcapitypes.Route, obj *svcsdk.CreateRouteInput) *svcsdk.CreateRouteInput {
	return obj
}

func postGenerateCreateRouteInput(cr *svcapitypes.Route, obj *svcsdk.CreateRouteInput) *svcsdk.CreateRouteInput {
	obj.ApiId = cr.Spec.ForProvider.APIID
	return obj
}

func preGenerateDeleteRouteInput(_ *svcapitypes.Route, obj *svcsdk.DeleteRouteInput) *svcsdk.DeleteRouteInput {
	return obj
}

func postGenerateDeleteRouteInput(cr *svcapitypes.Route, obj *svcsdk.DeleteRouteInput) *svcsdk.DeleteRouteInput {
	obj.ApiId = cr.Spec.ForProvider.APIID
	obj.RouteId = aws.String(meta.GetExternalName(cr))
	return obj
}
