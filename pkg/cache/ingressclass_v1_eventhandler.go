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

package cache

import (
	ingresspipy "github.com/flomesh-io/fsm/pkg/ingress"
	networkingv1 "k8s.io/api/networking/v1"
)

func (c *Cache) OnIngressClassv1Add(class *networkingv1.IngressClass) {
	c.updateDefaultIngressClass(class, class.Name)
}

func (c *Cache) OnIngressClassv1Update(oldClass, class *networkingv1.IngressClass) {
	if oldClass.ResourceVersion == class.ResourceVersion {
		return
	}

	c.updateDefaultIngressClass(class, class.Name)
}

func (c *Cache) OnIngressClassv1Delete(class *networkingv1.IngressClass) {
	// if the default IngressClass is deleted, set the DefaultIngressClass variable to empty
	c.updateDefaultIngressClass(class, ingresspipy.NoDefaultIngressClass)
}

func (c *Cache) OnIngressClassv1Synced() {
	c.mu.Lock()
	c.ingressClassesSynced = true
	c.setInitialized(c.ingressesSynced && c.servicesSynced && c.endpointsSynced)
	c.mu.Unlock()

	c.syncRoutes()
}

func (c *Cache) updateDefaultIngressClass(class *networkingv1.IngressClass, className string) {
	isDefault, ok := class.GetAnnotations()[ingresspipy.IngressClassAnnotationKey]
	if ok && isDefault == "true" {
		ingresspipy.DefaultIngressClass = className
	}
}
