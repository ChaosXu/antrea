// Copyright 2020 Antrea Authors
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

package main

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"antrea.io/antrea/pkg/clusteridentity"
	aggregator "antrea.io/antrea/pkg/flowaggregator"
	"antrea.io/antrea/pkg/flowaggregator/apiserver"
	"antrea.io/antrea/pkg/log"
	"antrea.io/antrea/pkg/signals"
	"antrea.io/antrea/pkg/util/cipher"
)

const informerDefaultResync = 12 * time.Hour

// genObservationDomainID generates an IPFIX Observation Domain ID when one is not provided by the
// user through the flow aggregator configuration. It will first try to generate one
// deterministically based on the cluster UUID (if available, with a timeout of 10s). Otherwise, it
// will generate a random one. The cluster UUID should be available if Antrea is deployed to the
// cluster ahead of the flow aggregator, which is the expectation since when deploying flow
// aggregator as a Pod, networking needs to be configured by the CNI plugin.
func genObservationDomainID(k8sClient kubernetes.Interface) uint32 {
	const retryInterval = time.Second
	const timeout = 10 * time.Second
	const defaultAntreaNamespace = "kube-system"

	clusterIdentityProvider := clusteridentity.NewClusterIdentityProvider(
		defaultAntreaNamespace,
		clusteridentity.DefaultClusterIdentityConfigMapName,
		k8sClient,
	)
	var clusterUUID uuid.UUID
	if err := wait.PollImmediate(retryInterval, timeout, func() (bool, error) {
		clusterIdentity, _, err := clusterIdentityProvider.Get()
		if err != nil {
			return false, nil
		}
		clusterUUID = clusterIdentity.UUID
		return true, nil
	}); err != nil {
		klog.Warningf(
			"Unable to retrieve cluster UUID after %v (does ConfigMap '%s/%s' exist?); will generate a random observation domain ID",
			timeout, defaultAntreaNamespace, clusteridentity.DefaultClusterIdentityConfigMapName,
		)
		clusterUUID = uuid.New()
	}
	h := fnv.New32()
	h.Write(clusterUUID[:])
	observationDomainID := h.Sum32()
	return observationDomainID
}

func run(o *Options) error {
	klog.Infof("Flow aggregator starting...")
	// Set up signal capture: the first SIGTERM / SIGINT signal is handled gracefully and will
	// cause the stopCh channel to be closed; if another signal is received before the program
	// exits, we will force exit.
	stopCh := signals.RegisterSignalHandlers()

	log.StartLogFileNumberMonitor(stopCh)

	k8sClient, err := createK8sClient()
	if err != nil {
		return fmt.Errorf("error when creating K8s client: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(k8sClient, informerDefaultResync)
	podInformer := informerFactory.Core().V1().Pods()

	var observationDomainID uint32
	if o.config.ObservationDomainID != nil {
		observationDomainID = *o.config.ObservationDomainID
	} else {
		observationDomainID = genObservationDomainID(k8sClient)
	}
	klog.Infof("Flow aggregator Observation Domain ID: %d", observationDomainID)

	var sendJSONRecord bool
	if o.format == "JSON" {
		sendJSONRecord = true
	} else {
		sendJSONRecord = false
	}

	flowAggregator := aggregator.NewFlowAggregator(
		o.externalFlowCollectorAddr,
		o.externalFlowCollectorProto,
		o.activeFlowRecordTimeout,
		o.inactiveFlowRecordTimeout,
		o.aggregatorTransportProtocol,
		o.flowAggregatorAddress,
		o.includePodLabels,
		k8sClient,
		observationDomainID,
		podInformer,
		sendJSONRecord,
	)
	err = flowAggregator.InitCollectingProcess()
	if err != nil {
		return fmt.Errorf("error when creating collecting process: %v", err)
	}
	err = flowAggregator.InitAggregationProcess()
	if err != nil {
		return fmt.Errorf("error when creating aggregation process: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go flowAggregator.Run(stopCh, &wg)

	cipherSuites, err := cipher.GenerateCipherSuitesList(o.config.APIServer.TLSCipherSuites)
	if err != nil {
		return fmt.Errorf("error generating Cipher Suite list: %v", err)
	}
	apiServer, err := apiserver.New(
		flowAggregator,
		o.config.APIServer.APIPort,
		cipherSuites,
		cipher.TLSVersionMap[o.config.APIServer.TLSMinVersion])
	if err != nil {
		return fmt.Errorf("error when creating flow aggregator API server: %v", err)
	}
	go apiServer.Run(stopCh)

	informerFactory.Start(stopCh)

	<-stopCh
	klog.Infof("Stopping flow aggregator")
	wg.Wait()
	return nil
}

func createK8sClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}
