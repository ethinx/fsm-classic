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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/flomesh-io/traffic-guru/apis/proxyprofile/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ProxyProfileLister helps list ProxyProfiles.
// All objects returned here must be treated as read-only.
type ProxyProfileLister interface {
	// List lists all ProxyProfiles in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ProxyProfile, err error)
	// Get retrieves the ProxyProfile from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ProxyProfile, error)
	ProxyProfileListerExpansion
}

// proxyProfileLister implements the ProxyProfileLister interface.
type proxyProfileLister struct {
	indexer cache.Indexer
}

// NewProxyProfileLister returns a new ProxyProfileLister.
func NewProxyProfileLister(indexer cache.Indexer) ProxyProfileLister {
	return &proxyProfileLister{indexer: indexer}
}

// List lists all ProxyProfiles in the indexer.
func (s *proxyProfileLister) List(selector labels.Selector) (ret []*v1alpha1.ProxyProfile, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ProxyProfile))
	})
	return ret, err
}

// Get retrieves the ProxyProfile from the index for a given name.
func (s *proxyProfileLister) Get(name string) (*v1alpha1.ProxyProfile, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("proxyprofile"), name)
	}
	return obj.(*v1alpha1.ProxyProfile), nil
}
