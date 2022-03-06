/*
 * The NEU License
 *
 * Copyright (c) 2022.  flomesh.io
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
 * of the Software, and to permit persons to whom the Software is furnished to do
 * so, subject to the following conditions:
 *
 * (1)The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * (2)If the software or part of the code will be directly used or used as a
 * component for commercial purposes, including but not limited to: public cloud
 *  services, hosting services, and/or commercial software, the logo as following
 *  shall be displayed in the eye-catching position of the introduction materials
 * of the relevant commercial services or products (such as website, product
 * publicity print), and the logo shall be linked or text marked with the
 * following URL.
 *
 * LOGO : http://flomesh.cn/assets/flomesh-logo.png
 * URL : https://github.com/flomesh-io
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package controller

import (
	"fmt"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	networkingv1informers "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"time"
)

type Ingressv1Handler interface {
	OnIngressv1Add(ingress *networkingv1.Ingress)
	OnIngressv1Update(oldIngress, ingress *networkingv1.Ingress)
	OnIngressv1Delete(ingress *networkingv1.Ingress)
	OnIngressv1Synced()
}

type Ingressv1Controller struct {
	Informer     cache.SharedIndexInformer
	Lister       Ingressv1Lister
	HasSynced    cache.InformerSynced
	eventHandler Ingressv1Handler
}

type Ingressv1Lister struct {
	cache.Store
}

func (l *Ingressv1Lister) ByKey(key string) (*networkingv1.Ingress, error) {
	s, exists, err := l.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("no object matching key %q in local store", key)
	}
	return s.(*networkingv1.Ingress), nil
}

func NewIngressv1ControllerWithEventHandler(ingressInformer networkingv1informers.IngressInformer, resyncPeriod time.Duration, handler Ingressv1Handler) *Ingressv1Controller {
	informer := ingressInformer.Informer()

	result := &Ingressv1Controller{
		HasSynced: informer.HasSynced,
		Informer:  informer,
		Lister: Ingressv1Lister{
			Store: informer.GetStore(),
		},
	}

	informer.AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    result.handleAddIngress,
			UpdateFunc: result.handleUpdateIngress,
			DeleteFunc: result.handleDeleteIngress,
		},
		resyncPeriod,
	)

	if handler != nil {
		result.eventHandler = handler
	}

	return result
}

func (c *Ingressv1Controller) Run(stopCh <-chan struct{}) {
	klog.InfoS("Starting ingress config controller")

	if !cache.WaitForNamedCacheSync("ingress v1 config", stopCh, c.HasSynced) {
		return
	}

	if c.eventHandler != nil {
		klog.V(3).Info("Calling handler.OnIngressv1Synced()")
		c.eventHandler.OnIngressv1Synced()
	}
}

func (c *Ingressv1Controller) handleAddIngress(obj interface{}) {
	ing, ok := obj.(*networkingv1.Ingress)
	if !ok {
		runtime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
		return
	}

	if c.eventHandler != nil {
		klog.V(4).Info("Calling handler.OnIngressv1Add")
		c.eventHandler.OnIngressv1Add(ing)
	}
}

func (c *Ingressv1Controller) handleUpdateIngress(oldObj, newObj interface{}) {
	oldIng, ok := oldObj.(*networkingv1.Ingress)
	if !ok {
		runtime.HandleError(fmt.Errorf("unexpected object type: %v", oldObj))
		return
	}
	ing, ok := newObj.(*networkingv1.Ingress)
	if !ok {
		runtime.HandleError(fmt.Errorf("unexpected object type: %v", newObj))
		return
	}

	if c.eventHandler != nil {
		klog.V(4).Info("Calling handler.OnIngressv1Update")
		c.eventHandler.OnIngressv1Update(oldIng, ing)
	}
}

func (c *Ingressv1Controller) handleDeleteIngress(obj interface{}) {
	ing, ok := obj.(*networkingv1.Ingress)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
			return
		}
		if ing, ok = tombstone.Obj.(*networkingv1.Ingress); !ok {
			runtime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
			return
		}

		if c.eventHandler != nil {
			klog.V(4).Info("Calling handler.OnIngressv1Delete")
			c.eventHandler.OnIngressv1Delete(ing)
		}
	}
}
