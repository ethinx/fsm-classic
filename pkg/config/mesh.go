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

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flomesh-io/fsm/pkg/commons"
	"github.com/flomesh-io/fsm/pkg/kube"
	"github.com/flomesh-io/fsm/pkg/util"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	v1 "k8s.io/client-go/listers/core/v1"
	k8scache "k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"time"
)

type MeshConfig struct {
	IsControlPlane        bool             `json:"is-control-plane,omitempty"`
	Ingress               Ingress          `json:"ingress,omitempty"`
	GatewayApi            GatewayApi       `json:"gateway-api,omitempty"`
	RepoRootURL           string           `json:"repo-root-url,omitempty"`
	RepoPath              string           `json:"repo-path,omitempty"`
	RepoApiPath           string           `json:"repo-api-path,omitempty"`
	ServiceAggregatorAddr string           `json:"service-aggregator-addr,omitempty"`
	DefaultPipyImage      string           `json:"default-pipy-image,omitempty"`
	ProxyInitImage        string           `json:"proxy-init-image,omitempty"`
	WaitForItImage        string           `json:"wait-for-it-image,omitempty"`
	Certificate           Certificate      `json:"certificate,omitempty"`
	Cluster               Cluster          `json:"cluster,omitempty"`
	ClusterConnector      ClusterConnector `json:"cluster-connector,omitempty"`
	WebhookServiceName    string           `json:"webhook-service-name,omitempty"`
}

type Ingress struct {
	Enabled     bool   `json:"enabled,omitempty"`
	DeployName  string `json:"deploy-name,omitempty"`
	ServiceName string `json:"service-name,omitempty"`
}

type GatewayApi struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Cluster struct {
	Region string `json:"region,omitempty"`
	Zone   string `json:"zone,omitempty"`
	Group  string `json:"group,omitempty"`
	Name   string `json:"name,omitempty"`
}

type ClusterConnector struct {
	DefaultImage       string `json:"default-image,omitempty"`
	SecretMountPath    string `json:"secret-mount-path,omitempty"`
	ConfigmapName      string `json:"configmap-name,omitempty"`
	ConfigFile         string `json:"config-file,omitempty"`
	LogLevel           int32  `json:"log-level,omitempty"`
	ServiceAccountName string `json:"service-account-name,omitempty"`
}

type Certificate struct {
	Manager string `json:"manager,omitempty"`
}

type MeshConfigClient struct {
	k8sApi   *kube.K8sAPI
	cmLister v1.ConfigMapNamespaceLister
}

func NewMeshConfigClient(k8sApi *kube.K8sAPI) *MeshConfigClient {
	informerFactory := informers.NewSharedInformerFactoryWithOptions(k8sApi.Client, 60*time.Second, informers.WithNamespace(GetFsmNamespace()))
	configmapLister := informerFactory.Core().V1().ConfigMaps().Lister().ConfigMaps(GetFsmNamespace())
	configmapInformer := informerFactory.Core().V1().ConfigMaps().Informer()
	go configmapInformer.Run(wait.NeverStop)

	if !k8scache.WaitForCacheSync(wait.NeverStop, configmapInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for configmap to sync"))
	}

	return &MeshConfigClient{
		k8sApi:   k8sApi,
		cmLister: configmapLister,
	}
}

func (o *MeshConfig) RepoBaseURL() string {
	return fmt.Sprintf("%s%s", o.RepoRootURL, o.RepoPath)
}

func (o *MeshConfig) RepoApiBaseURL() string {
	return fmt.Sprintf("%s%s", o.RepoRootURL, o.RepoApiPath)
}

func (o *MeshConfig) IngressCodebasePath() string {
	return util.EvaluateTemplate(commons.IngressPathTemplate, struct {
		Region  string
		Zone    string
		Group   string
		Cluster string
	}{
		Region:  o.Cluster.Region,
		Zone:    o.Cluster.Zone,
		Group:   o.Cluster.Group,
		Cluster: o.Cluster.Name,
	}) + "/"
}

func (o *MeshConfig) ToJson() string {
	cfgBytes, err := json.Marshal(o)
	if err != nil {
		klog.Errorf("Not able to marshal MeshConfig %#v to json, %s", o, err.Error())
		return ""
	}

	return string(cfgBytes)
}

func (c *MeshConfigClient) GetConfig() *MeshConfig {
	cm := c.getConfigMap()

	if cm != nil {
		return ParseMeshConfig(cm)
	}

	return nil
}

func (c *MeshConfigClient) UpdateConfig(config *MeshConfig) {
	cm := c.getConfigMap()

	if cm == nil {
		return
	}

	cfgBytes, err := json.Marshal(config)
	if err != nil {
		klog.Errorf("Not able to marshal MeshConfig %#v to json, %s", config, err.Error())
		return
	}
	cm.Data[commons.MeshConfigJsonName] = string(cfgBytes)

	cm, err = c.k8sApi.Client.CoreV1().
		ConfigMaps(GetFsmNamespace()).
		Update(context.TODO(), cm, metav1.UpdateOptions{})

	if err != nil {
		klog.Errorf("Update ConfigMap flomesh/mesh-config error, %s", err.Error())
		return
	}

	klog.V(5).Infof("After updating, ConfigMap flomesh/mesh-config = %#v", cm)
}

func (c *MeshConfigClient) getConfigMap() *corev1.ConfigMap {
	cm, err := c.cmLister.Get(commons.MeshConfigName)

	if err != nil {
		// it takes time to sync, perhaps still not in the local store yet
		if apierrors.IsNotFound(err) {
			cm, err = c.k8sApi.Client.CoreV1().
				ConfigMaps(GetFsmNamespace()).
				Get(context.TODO(), commons.MeshConfigName, metav1.GetOptions{})

			if err != nil {
				klog.Errorf("Get ConfigMap flomesh/mesh-config from API server error, %s", err.Error())
				return nil
			}
		} else {
			klog.Errorf("Get ConfigMap flomesh/mesh-config error, %s", err.Error())
			return nil
		}
	}

	return cm
}

func ParseMeshConfig(cm *corev1.ConfigMap) *MeshConfig {
	cfgJson, ok := cm.Data[commons.MeshConfigJsonName]
	if !ok {
		klog.Error("Config file mesh_config.json not found, please check ConfigMap flomesh/mesh-config.")
		return nil
	}
	klog.V(5).Infof("Found mesh_config.json, content: %s", cfgJson)

	cfg := MeshConfig{}
	err := json.Unmarshal([]byte(cfgJson), &cfg)
	if err != nil {
		klog.Errorf("Unable to unmarshal mesh_config.json to config.MeshConfig, %s", err.Error())
		return nil
	}

	return &cfg
}
