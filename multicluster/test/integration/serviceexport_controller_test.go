// Copyright 2021 Antrea Authors.
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

package integration

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8smcsv1alpha1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"

	mcsv1alpha1 "antrea.io/antrea/multicluster/apis/multicluster/v1alpha1"
	"antrea.io/antrea/multicluster/controllers/multicluster/common"
)

// This file contains test cases for below basic scenarios:
//  * Create ResourceExports when a ServiceExport is created.
//  * Update ResourceExport when exported Service is updated.
//  * Update ServiceExport status when the Service doesn't exist
//  * Update ResourceExport when the Endpoints has new Endpoints
//  * Delete ResourceExport when the ServiceExport is deleted

var _ = Describe("ServiceExport controller", func() {
	svcSpec := corev1.ServiceSpec{
		Ports: svcPorts,
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-svc",
			Namespace: testNamespace,
		},
		Spec: svcSpec,
	}
	svcNamespacedName := types.NamespacedName{
		Namespace: svc.Namespace,
		Name:      svc.Name,
	}
	epNamespacedName := svcNamespacedName

	ep := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-svc",
			Namespace: testNamespace,
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					addr1,
				},
				NotReadyAddresses: []corev1.EndpointAddress{
					addr2,
				},
				Ports: epPorts,
			},
		},
	}
	svcExport := &k8smcsv1alpha1.ServiceExport{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-svc",
			Namespace: testNamespace,
		},
	}
	svcExportNoService := &k8smcsv1alpha1.ServiceExport{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-svc-dont-exist",
			Namespace: testNamespace,
		},
	}
	svcResExportName := LocalClusterID + "-" + svc.Namespace + "-" + svc.Name + "-service"
	epResExportName := LocalClusterID + "-" + ep.Namespace + "-" + ep.Name + "-endpoints"

	expectedEpResExport := &mcsv1alpha1.ResourceExport{
		ObjectMeta: metav1.ObjectMeta{
			Name:      epResExportName,
			Namespace: LeaderNamespace,
		},
		Spec: mcsv1alpha1.ResourceExportSpec{
			ClusterID: LocalClusterID,
			Name:      ep.Name,
			Namespace: ep.Namespace,
			Kind:      common.EndpointsKind,
		},
	}

	ctx := context.Background()
	It("Should create ResourceExports when new ServiceExport for ClusterIP Service is created", func() {
		By("By exposing a ClusterIP type of Service")
		expectedEpResExport.Spec.Endpoints = &mcsv1alpha1.EndpointsExport{
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							IP: "192.168.17.11",
						},
					},
					Ports: epPorts,
				},
			},
		}
		Expect(k8sClient.Create(ctx, svc)).Should(Succeed())
		Expect(k8sClient.Create(ctx, ep)).Should(Succeed())
		Expect(k8sClient.Create(ctx, svcExport)).Should(Succeed())
		var err error
		latestSvc := &corev1.Service{}
		err = k8sClient.Get(ctx, svcNamespacedName, latestSvc)
		Expect(err).ToNot(HaveOccurred())
		svcResExport := &mcsv1alpha1.ResourceExport{}
		epResExport := &mcsv1alpha1.ResourceExport{}
		Eventually(func() bool {
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: LeaderNamespace, Name: svcResExportName}, svcResExport)
			return err == nil
		}, timeout, interval).Should(BeTrue())
		Expect(svcResExport.ObjectMeta.Labels["sourceKind"]).Should(Equal("Service"))
		Expect(svcResExport.Spec.Service.ServiceSpec.ClusterIP).Should(Equal(latestSvc.Spec.ClusterIP))
		Expect(len(svcResExport.Spec.Service.ServiceSpec.Ports)).Should(Equal(len(svcPorts)))
		Eventually(func() bool {
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: LeaderNamespace, Name: epResExportName}, epResExport)
			return err == nil
		}, timeout, interval).Should(BeTrue())
		Expect(epResExport.ObjectMeta.Labels["sourceKind"]).Should(Equal("Endpoints"))
		Expect(epResExport.Spec).Should(Equal(expectedEpResExport.Spec))
	})

	It("Should update existing ResourceExport when existing Service is updated", func() {
		By("By update Service's ports")
		newPorts := []corev1.ServicePort{
			{
				Name:     "udp88",
				Protocol: "UDP",
				Port:     88,
			},
		}
		latestSvc := &corev1.Service{}
		err := k8sClient.Get(ctx, svcNamespacedName, latestSvc)
		Expect(err).ToNot(HaveOccurred())
		latestSvc.Spec.Ports = newPorts
		Expect(k8sClient.Update(ctx, latestSvc)).Should(Succeed())

		svcResExport := &mcsv1alpha1.ResourceExport{}
		time.Sleep(5 * time.Second)
		Eventually(func() bool {
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: LeaderNamespace, Name: svcResExportName}, svcResExport)
			return err == nil
		}, timeout, interval).Should(BeTrue())
		Expect(svcResExport.ObjectMeta.Labels["sourceKind"]).Should(Equal("Service"))
		Expect(svcResExport.Spec.Service.ServiceSpec.ClusterIP).Should(Equal(latestSvc.Spec.ClusterIP))
		Expect(len(svcResExport.Spec.Service.ServiceSpec.Ports)).Should(Equal(len(newPorts)))
	})

	It("Should update existing ServiceExport status when corresponding Service doesn't exist", func() {
		By("By create a ServiceExport without a real Service")
		Expect(k8sClient.Create(ctx, svcExportNoService)).Should(Succeed())
		time.Sleep(2 * time.Second)
		latestSvcExportNoService := &k8smcsv1alpha1.ServiceExport{}
		err := k8sClient.Get(ctx, types.NamespacedName{
			Namespace: svcExportNoService.Namespace,
			Name:      svcExportNoService.Name,
		}, latestSvcExportNoService)
		Expect(err).ToNot(HaveOccurred())
		conditions := latestSvcExportNoService.Status.Conditions
		Expect(len(conditions)).Should(Equal(1))
		Expect(*conditions[0].Message).Should(Equal("the Service does not exist"))
	})

	It("Should update existing ResourceExport when corresponding Endpoints has new Endpoint", func() {
		By("By update an Endpoint with a new address")
		latestEp := &corev1.Endpoints{}
		Expect(k8sClient.Get(ctx, epNamespacedName, latestEp)).Should(Succeed())
		addresses := latestEp.Subsets[0].Addresses
		addresses = append(addresses, addr3)
		latestEp.Subsets[0].Addresses = addresses
		Expect(k8sClient.Update(ctx, latestEp)).Should(Succeed())
		time.Sleep(2 * time.Second)
		epResExport := &mcsv1alpha1.ResourceExport{}
		expectedEpResExport.Spec.Endpoints = &mcsv1alpha1.EndpointsExport{
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							IP: "192.168.17.11",
						},
						{
							IP: "192.168.17.13",
						},
					},
					Ports: epPorts,
				},
			},
		}

		var err error
		Eventually(func() bool {
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: LeaderNamespace, Name: epResExportName}, epResExport)
			return err == nil
		}, timeout, interval).Should(BeTrue())
		Expect(epResExport.ObjectMeta.Labels["sourceKind"]).Should(Equal("Endpoints"))
		Expect(epResExport.Spec).Should(Equal(expectedEpResExport.Spec))

	})

	It("Should delete existing ResourceExport when existing ServiceExport is deleted", func() {
		By("By remove a ServiceExport resource")
		err := k8sClient.Delete(ctx, svcExport)
		Expect(err).ToNot(HaveOccurred())
		resExp := &mcsv1alpha1.ResourceExport{}
		Eventually(func() bool {
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: LeaderNamespace, Name: svcResExportName}, resExp)
			return apierrors.IsNotFound(err)
		}, timeout, interval).Should(BeTrue())
		Eventually(func() bool {
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: LeaderNamespace, Name: epResExportName}, resExp)
			return apierrors.IsNotFound(err)
		}, timeout, interval).Should(BeTrue())
	})
})
