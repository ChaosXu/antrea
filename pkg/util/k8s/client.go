// Copyright 2019 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import (
	apiextensionclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/klog/v2"
	aggregatorclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"

	crdclientset "antrea.io/antrea/pkg/client/clientset/versioned"
	legacycrdclientset "antrea.io/antrea/pkg/legacyclient/clientset/versioned"

	netdefclient "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned/typed/k8s.cni.cncf.io/v1"
)

// CreateClients creates kube clients from the given config.
func CreateClients(config componentbaseconfig.ClientConnectionConfiguration, kubeAPIServerOverride string) (
	clientset.Interface, aggregatorclientset.Interface, crdclientset.Interface, apiextensionclientset.Interface, error) {
	kubeConfig, err := createRestConfig(config, kubeAPIServerOverride)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	client, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	aggregatorClient, err := aggregatorclientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// Create client for crd operations
	crdClient, err := crdclientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// Create client for crd manipulations
	apiExtensionClient, err := apiextensionclientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return client, aggregatorClient, crdClient, apiExtensionClient, nil
}

// CreateLegacyCRDClient creates legacyCRD client from the given config.
func CreateLegacyCRDClient(config componentbaseconfig.ClientConnectionConfiguration, kubeAPIServerOverride string) (legacycrdclientset.Interface, error) {
	kubeConfig, err := createRestConfig(config, kubeAPIServerOverride)
	if err != nil {
		return nil, err
	}

	legacyCrdClient, err := legacycrdclientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return legacyCrdClient, nil
}

// CreateNetworkAttachDefClient creates net-attach-def client handle from the given config.
func CreateNetworkAttachDefClient(config componentbaseconfig.ClientConnectionConfiguration, kubeAPIServerOverride string) (netdefclient.K8sCniCncfIoV1Interface, error) {
	kubeConfig, err := createRestConfig(config, kubeAPIServerOverride)
	if err != nil {
		return nil, err
	}

	netAttachDefClient, err := netdefclient.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return netAttachDefClient, nil

}

func createRestConfig(config componentbaseconfig.ClientConnectionConfiguration, kubeAPIServerOverride string) (*rest.Config, error) {
	var kubeConfig *rest.Config
	var err error

	if len(config.Kubeconfig) == 0 {
		klog.Info("No kubeconfig file was specified. Falling back to in-cluster config")
		kubeConfig, err = rest.InClusterConfig()
	} else {
		kubeConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: config.Kubeconfig},
			&clientcmd.ConfigOverrides{}).ClientConfig()
	}

	if len(kubeAPIServerOverride) != 0 {
		kubeConfig.Host = kubeAPIServerOverride
	}

	if err != nil {
		return nil, err
	}

	kubeConfig.AcceptContentTypes = config.AcceptContentTypes
	kubeConfig.ContentType = config.ContentType
	kubeConfig.QPS = config.QPS
	kubeConfig.Burst = int(config.Burst)

	return kubeConfig, nil

}
