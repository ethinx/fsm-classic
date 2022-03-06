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

package collector

import (
	"fmt"
	"github.com/flomesh-io/fsm/pkg/commons"
	"github.com/flomesh-io/fsm/pkg/repo"
	routepkg "github.com/flomesh-io/fsm/pkg/route"
	"github.com/flomesh-io/fsm/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	"net/http"
	"sync"
	"time"
)

type Collector struct {
	imu              sync.RWMutex // For ingress repo
	smu              sync.RWMutex // For service repo
	address          string
	ingressStore     map[string]routepkg.IngressRoute
	serviceStore     map[string]routepkg.ServiceRoute
	ingressHashStore map[string]string
	serviceHashStore map[string]string
	ingressVersion   string
	serviceVersion   string
	crontab          *cron.Cron
	router           *gin.Engine

	repoClient *repo.PipyRepoClient
}

// FIXME: if a cluster is unlinked to the control plane, should remove all its base repos from pipy

func NewCollector(listenAddr string, repoAddr string) *Collector {
	c := &Collector{
		address:          listenAddr,
		ingressStore:     make(map[string]routepkg.IngressRoute),
		serviceStore:     make(map[string]routepkg.ServiceRoute),
		ingressHashStore: make(map[string]string),
		serviceHashStore: make(map[string]string),
		repoClient:       repo.NewRepoClient(repoAddr),
	}
	c.router = c.newEngine()

	c.crontab = cron.New(cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	c.crontab.AddFunc("@every 10s", c.syncToRepo)

	return c
}

func (c *Collector) newEngine() *gin.Engine {
	router := gin.Default()

	router.POST(ServicePath, c.service)
	router.POST(IngressPath, c.ingress)

	router.GET(HealthPath, health)
	router.GET(ReadyPath, health)

	return router
}

func (c *Collector) Run() error {
	klog.Infof("Starting service-collector server...")

	go c.crontab.Start()

	// defaults, listens on 8080
	if len(c.address) == 0 {
		return c.router.Run()
	}

	return c.router.Run(c.address)
}

func (c *Collector) service(context *gin.Context) {
	var serviceRoutes routepkg.ServiceRoute
	if err := context.BindJSON(&serviceRoutes); err != nil {
		context.JSON(http.StatusBadRequest, &commons.Response{
			Success: false,
			Result:  err.Error(),
		})
		return
	}

	// TODO: reduce unnecessary updating
	klog.V(5).Infof("Received ServiceRoute: %#v", serviceRoutes)
	if serviceRoutes.UID != "" {
		c.smu.Lock()

		oldHash := c.serviceHashStore[serviceRoutes.UID]
		hash := serviceRoutes.Hash
		if oldHash != hash {
			c.serviceStore[serviceRoutes.UID] = serviceRoutes
			c.serviceHashStore[serviceRoutes.UID] = hash
		}

		c.smu.Unlock()
	}
	klog.V(5).Infof("Size of service repo = %d", len(c.serviceStore))

	context.JSON(http.StatusOK, &commons.Response{Success: true})
}

func (c *Collector) ingress(context *gin.Context) {
	var ingressRoutes routepkg.IngressRoute
	if err := context.BindJSON(&ingressRoutes); err != nil {
		context.JSON(http.StatusBadRequest, &commons.Response{
			Success: false,
			Result:  err.Error(),
		})
		return
	}

	// TODO: reduce unnecessary updating
	klog.V(5).Infof("Received IngressRoute: %#v", ingressRoutes)
	if ingressRoutes.UID != "" {
		c.imu.Lock()

		oldHash := c.ingressHashStore[ingressRoutes.UID]
		hash := ingressRoutes.Hash
		if oldHash != hash {
			c.ingressStore[ingressRoutes.UID] = ingressRoutes
			c.ingressHashStore[ingressRoutes.UID] = hash
		}

		c.imu.Unlock()
	}
	klog.V(5).Infof("Size of ingress repo = %d", len(c.ingressStore))

	context.JSON(http.StatusOK, &commons.Response{Success: true})
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func (c *Collector) syncToRepo() {
	klog.V(5).Infof("Syncing routes to repo ======> %s", time.Now().String())
	go c.syncIngress()
	go c.syncService()
}

func (c *Collector) syncIngress() {
	klog.V(5).Infof("Syncing ingresses to pipy repo ...")
	c.imu.RLock()
	defer c.imu.RUnlock()

	version := util.SimpleHash(c.ingressHashStore)
	klog.V(5).Infof("Old hash = %q, new hash = %q", c.ingressVersion, version)
	if version == c.ingressVersion {
		klog.V(5).Infof("Ingress info doesn't change, ignore sync ...")
		return
	}
	c.ingressVersion = version

	// TODO: check if the ingress route item is updated,
	//   then generate the batch item accordingly
	var batches []repo.Batch
	for uid, ing := range c.ingressStore {
		klog.V(5).Infof("UID %q -- size: %d", uid, len(ing.Routes))
		batch := repo.Batch{
			Basepath: fmt.Sprintf("/%s/%s/%s/%s/ingress", ing.Region, ing.Zone, ing.Group, ing.Cluster),
			Items:    []repo.BatchItem{},
		}
		// Generate router.json
		router := repo.Router{Routes: repo.RouterEntry{}}
		// Generate balancer.json
		balancer := repo.Balancer{Services: repo.BalancerEntry{}}

		for _, r := range ing.Routes {
			// router
			//router.Routes[] = append(router.Routes, routerEntry(r))
			router.Routes[routerKey(r)] = repo.ServiceInfo{Service: r.ServiceName, Rewrite: r.Rewrite}

			// balancer
			//balancer.Services = append(balancer.Services, balancerEntry(r))
			balancer.Services[r.ServiceName] = upstream(r)
		}

		batch.Items = append(batch.Items, ingressBatchItems(router, balancer)...)
		if len(batch.Items) > 0 {
			batches = append(batches, batch)
		}
	}

	if len(batches) > 0 {
		klog.V(5).Infof("Updating pipy repo, ingress batches: %#v", batches)
		go c.repoClient.Batch(batches)
	}
}

func routerKey(r routepkg.IngressRouteEntry) string {
	return fmt.Sprintf("%s%s", r.Host, r.Path)
}

func upstream(r routepkg.IngressRouteEntry) repo.Upstream {
	return repo.Upstream{
		Balancer: repo.RoundRobinLoadBalancer,
		Sticky:   false,
		Targets:  transformTargets(r.Upstreams),
	}
}

//func routerEntry(r routepkg.IngressRouteEntry) repo.RouterEntry {
//	entry := repo.RouterEntry{}
//	entry[fmt.Sprintf("%s%s", r.Host, r.Path)] = repo.ServiceInfo{
//		Service: r.ServiceName,
//	}
//	return entry
//}

//func balancerEntry(r routepkg.IngressRouteEntry) repo.BalancerEntry {
//	entry := repo.BalancerEntry{}
//	entry[r.ServiceName] = repo.Upstream{
//		Balancer: repo.RoundRobinLoadBalancer,
//		Sticky:   false,
//		Targets:  transformTargets(r.Upstreams),
//	}
//	return entry
//}

func transformTargets(endpoints []routepkg.EndpointEntry) []string {
	if len(endpoints) == 0 {
		return []string{}
	}

	targets := sets.String{}
	for _, ep := range endpoints {
		targets.Insert(fmt.Sprintf("%s:%d", ep.IP, ep.Port))
	}

	return targets.List()
}

func ingressBatchItems(router repo.Router, balancer repo.Balancer) []repo.BatchItem {
	routerItem := repo.BatchItem{
		Path:     "/config",
		Filename: "router.json",
		Content:  router,
	}
	balancerItem := repo.BatchItem{
		Path:     "/config",
		Filename: "balancer.json",
		Content:  balancer,
	}
	return []repo.BatchItem{routerItem, balancerItem}
}

func (c *Collector) syncService() {
	klog.V(5).Infof("Syncing services to pipy repo ...")
	c.smu.RLock()
	defer c.smu.RUnlock()

	version := util.SimpleHash(c.serviceHashStore)
	klog.V(5).Infof("Old hash = %q, new hash = %q", c.serviceVersion, version)
	if version == c.serviceVersion {
		klog.V(5).Infof("Services info doesn't change, ignore sync ...")
		return
	}
	c.serviceVersion = version

	regs := make(map[string]repo.ServiceRegistry)
	for uid := range c.serviceStore {
		regs[uid] = repo.ServiceRegistry{Services: repo.ServiceRegistryEntry{}}
	}

	for uid, sr := range c.serviceStore {
		klog.V(5).Infof("UID %q -- size: %d", uid, len(sr.Routes))

		for _, route := range sr.Routes {
			serviceName := servicePortName(route)
			exportServiceName := exportServicePortName(route)
			internal := internalAddr(route)
			external := externalAddr(route)

			for ruid := range regs {
				// merge the service by serviceName
				if uid == ruid {
					regs[ruid].Services[serviceName] = append(regs[ruid].Services[serviceName], internal...)
				} else {
					//  1. when merging the service, should use exported name
					//  2. if the exported flag of a serive is false, it should ONLY show in the registry it belongs to
					//  3. ONLY if the service is exposed through Ingress, it has an external address
					if len(external) > 0 && route.Export {
						regs[ruid].Services[exportServiceName] = append(regs[ruid].Services[exportServiceName], external)
					}
				}
			}
		}
	}

	var batches []repo.Batch
	for uid, sr := range c.serviceStore {
		batch := repo.Batch{
			Basepath: fmt.Sprintf("/%s/%s/%s/%s/services", sr.Region, sr.Zone, sr.Group, sr.Cluster),
			Items:    []repo.BatchItem{},
		}

		item := repo.BatchItem{
			Path:     "/config",
			Filename: "registry.json",
			Content:  regs[uid],
		}

		batch.Items = append(batch.Items, item)
		if len(batch.Items) > 0 {
			batches = append(batches, batch)
		}
	}

	if len(batches) > 0 {
		klog.V(5).Infof("Updating pipy repo, service batches: %#v", batches)
		go c.repoClient.Batch(batches)
	}
}

func servicePortName(route routepkg.ServiceRouteEntry) string {
	return fmt.Sprintf("%s/%s%s", route.Namespace, route.Name, fmtPortName(route.PortName))
}

func exportServicePortName(route routepkg.ServiceRouteEntry) string {
	return fmt.Sprintf("%s/%s%s", route.Namespace, svcName(route), fmtPortName(route.PortName))
}

func svcName(route routepkg.ServiceRouteEntry) string {
	if route.Export && route.ExportName != "" {
		return route.ExportName
	}

	return route.Name
}

func fmtPortName(in string) string {
	if in == "" {
		return ""
	}
	return fmt.Sprintf(":%s", in)
}

func internalAddr(route routepkg.ServiceRouteEntry) []string {
	//entry := repo.ServiceRegistryEntry{}
	//serviceName := fmt.Sprintf("%s.%s", route.Name, route.Namespace)
	//serviceAddress := fmt.Sprintf("%s:%d", route.IP, route.Port)
	//entry[serviceName] = []string{serviceAddress}
	//
	//return entry

	result := make([]string, 0)
	for _, target := range route.Targets {
		result = append(result, target.Address)
	}

	return result
}

func externalAddr(route routepkg.ServiceRouteEntry) string {
	//entry := repo.ServiceRegistryEntry{}
	//serviceName := fmt.Sprintf("%s.%s", route.Name, route.Namespace)
	//serviceAddress := route.ExternalPath
	//entry[serviceName] = []string{serviceAddress}
	//
	//return entry
	return route.ExternalPath
}
