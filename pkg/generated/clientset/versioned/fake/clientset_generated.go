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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	clientset "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned"
	clusterv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/cluster/v1alpha1"
	fakeclusterv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/cluster/v1alpha1/fake"
	globaltrafficpolicyv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/globaltrafficpolicy/v1alpha1"
	fakeglobaltrafficpolicyv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/globaltrafficpolicy/v1alpha1/fake"
	multiclusterendpointv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/multiclusterendpoint/v1alpha1"
	fakemulticlusterendpointv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/multiclusterendpoint/v1alpha1/fake"
	namespacedingressv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/namespacedingress/v1alpha1"
	fakenamespacedingressv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/namespacedingress/v1alpha1/fake"
	proxyprofilev1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/proxyprofile/v1alpha1"
	fakeproxyprofilev1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/proxyprofile/v1alpha1/fake"
	serviceexportv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/serviceexport/v1alpha1"
	fakeserviceexportv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/serviceexport/v1alpha1/fake"
	serviceimportv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/serviceimport/v1alpha1"
	fakeserviceimportv1alpha1 "github.com/flomesh-io/fsm/pkg/generated/clientset/versioned/typed/serviceimport/v1alpha1/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	cs := &Clientset{tracker: o}
	cs.discovery = &fakediscovery.FakeDiscovery{Fake: &cs.Fake}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})

	return cs
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
	tracker   testing.ObjectTracker
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *Clientset) Tracker() testing.ObjectTracker {
	return c.tracker
}

var (
	_ clientset.Interface = &Clientset{}
	_ testing.FakeClient  = &Clientset{}
)

// ClusterV1alpha1 retrieves the ClusterV1alpha1Client
func (c *Clientset) ClusterV1alpha1() clusterv1alpha1.ClusterV1alpha1Interface {
	return &fakeclusterv1alpha1.FakeClusterV1alpha1{Fake: &c.Fake}
}

// GlobaltrafficpolicyV1alpha1 retrieves the GlobaltrafficpolicyV1alpha1Client
func (c *Clientset) GlobaltrafficpolicyV1alpha1() globaltrafficpolicyv1alpha1.GlobaltrafficpolicyV1alpha1Interface {
	return &fakeglobaltrafficpolicyv1alpha1.FakeGlobaltrafficpolicyV1alpha1{Fake: &c.Fake}
}

// MulticlusterendpointV1alpha1 retrieves the MulticlusterendpointV1alpha1Client
func (c *Clientset) MulticlusterendpointV1alpha1() multiclusterendpointv1alpha1.MulticlusterendpointV1alpha1Interface {
	return &fakemulticlusterendpointv1alpha1.FakeMulticlusterendpointV1alpha1{Fake: &c.Fake}
}

// NamespacedingressV1alpha1 retrieves the NamespacedingressV1alpha1Client
func (c *Clientset) NamespacedingressV1alpha1() namespacedingressv1alpha1.NamespacedingressV1alpha1Interface {
	return &fakenamespacedingressv1alpha1.FakeNamespacedingressV1alpha1{Fake: &c.Fake}
}

// ProxyprofileV1alpha1 retrieves the ProxyprofileV1alpha1Client
func (c *Clientset) ProxyprofileV1alpha1() proxyprofilev1alpha1.ProxyprofileV1alpha1Interface {
	return &fakeproxyprofilev1alpha1.FakeProxyprofileV1alpha1{Fake: &c.Fake}
}

// ServiceexportV1alpha1 retrieves the ServiceexportV1alpha1Client
func (c *Clientset) ServiceexportV1alpha1() serviceexportv1alpha1.ServiceexportV1alpha1Interface {
	return &fakeserviceexportv1alpha1.FakeServiceexportV1alpha1{Fake: &c.Fake}
}

// ServiceimportV1alpha1 retrieves the ServiceimportV1alpha1Client
func (c *Clientset) ServiceimportV1alpha1() serviceimportv1alpha1.ServiceimportV1alpha1Interface {
	return &fakeserviceimportv1alpha1.FakeServiceimportV1alpha1{Fake: &c.Fake}
}
