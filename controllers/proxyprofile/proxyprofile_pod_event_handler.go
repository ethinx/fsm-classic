/*
 * The NEU License
 *
 * Copyright (c) 2021-2022.  flomesh.io
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

package proxyprofile

import (
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

var _ handler.EventHandler = &PodRequestHandler{}

type PodRequestHandler struct {
	client.Client
}

func (p *PodRequestHandler) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	//log.Infof("ProxyProfilePodRequestHanlder - Create(), event=%#v", evt.Object)
}

func (p *PodRequestHandler) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	//log.Infof("ProxyProfilePodRequestHanlder - Delete(), event=%#v", evt)
}

func (p *PodRequestHandler) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	//log.Infof("ProxyProfilePodRequestHanlder - Generic(), event=%#v", evt)
}

func (p *PodRequestHandler) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	//log.Infof("ProxyProfilePodRequestHanlder - Update(), event=%#v", evt)
}
