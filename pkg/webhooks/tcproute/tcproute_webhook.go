/*
 * MIT License
 *
 * Copyright (c) since 2021,  flomesh.io Authors.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package tcproute

import (
	flomeshadmission "github.com/flomesh-io/traffic-guru/pkg/admission"
	"github.com/flomesh-io/traffic-guru/pkg/commons"
	"github.com/flomesh-io/traffic-guru/pkg/config"
	"github.com/flomesh-io/traffic-guru/pkg/kube"
	admissionregv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/klog/v2"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

const (
	kind      = "TCPRoute"
	groups    = "gateway.networking.k8s.io"
	resources = "tcproutes"
	versions  = "v1alpha2"

	mwPath = commons.TCPRouteMutatingWebhookPath
	mwName = "mtcproute.kb.flomesh.io"
	vwPath = commons.TCPRouteValidatingWebhookPath
	vwName = "vtcproute.kb.flomesh.io"
)

func RegisterWebhooks(caBundle []byte) {
	rule := flomeshadmission.NewRule(
		[]admissionregv1.OperationType{admissionregv1.Create, admissionregv1.Update},
		[]string{groups},
		[]string{versions},
		[]string{resources},
	)

	mutatingWebhook := flomeshadmission.NewMutatingWebhook(
		mwName,
		mwPath,
		caBundle,
		nil,
		[]admissionregv1.RuleWithOperations{rule},
	)

	validatingWebhook := flomeshadmission.NewValidatingWebhook(
		vwName,
		vwPath,
		caBundle,
		nil,
		[]admissionregv1.RuleWithOperations{rule},
	)

	flomeshadmission.RegisterMutatingWebhook(mwName, mutatingWebhook)
	flomeshadmission.RegisterValidatingWebhook(vwName, validatingWebhook)
}

type TCPRouteDefaulter struct {
	k8sAPI *kube.K8sAPI
}

//var _ webhooks.Defaulter = &TCPRouteDefaulter{}

func NewDefaulter(k8sAPI *kube.K8sAPI) *TCPRouteDefaulter {
	return &TCPRouteDefaulter{
		k8sAPI: k8sAPI,
	}
}

func (w *TCPRouteDefaulter) Kind() string {
	return kind
}

func (w *TCPRouteDefaulter) SetDefaults(obj interface{}) {
	route, ok := obj.(*gwv1alpha2.TCPRoute)
	if !ok {
		return
	}

	klog.V(5).Infof("Default Webhook, name=%s", route.Name)
	klog.V(4).Infof("Before setting default values, spec=%#v", route.Spec)

	meshConfig := config.GetMeshConfig(w.k8sAPI)

	if meshConfig == nil {
		return
	}

	klog.V(4).Infof("After setting default values, spec=%#v", route.Spec)
}

type TCPRouteValidator struct {
	k8sAPI *kube.K8sAPI
}

func (w *TCPRouteValidator) Kind() string {
	return kind
}

func (w *TCPRouteValidator) ValidateCreate(obj interface{}) error {
	return doValidation(obj)
}

func (w *TCPRouteValidator) ValidateUpdate(oldObj, obj interface{}) error {
	return doValidation(obj)
}

func (w *TCPRouteValidator) ValidateDelete(obj interface{}) error {
	return nil
}

//var _ webhooks.Validator = &TCPRouteValidator{}

func NewValidator(k8sAPI *kube.K8sAPI) *TCPRouteValidator {
	return &TCPRouteValidator{
		k8sAPI: k8sAPI,
	}
}

func doValidation(obj interface{}) error {
	//route, ok := obj.(*gwv1alpha2.TCPRoute)
	//if !ok {
	//    return nil
	//}

	return nil
}
