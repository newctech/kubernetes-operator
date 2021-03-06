/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"testing"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	ecsfake "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned/fake"
	informers "github.com/gosoon/kubernetes-operator/pkg/client/informers/externalversions"
	"github.com/gosoon/kubernetes-operator/pkg/enum"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/controller"
)

var alwaysReady = func() bool { return true }

func TestProcessOperateFinished(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	kubernetesClusterClient := ecsfake.NewSimpleClientset()
	informerFactory := informers.NewSharedInformerFactory(kubernetesClusterClient, controller.NoResyncPeriodFunc())
	ecsv1Controller := NewController(fakeClient, kubernetesClusterClient,
		informerFactory.Ecs().V1().KubernetesClusters())
	ecsv1Controller.kubernetesClusterSynced = alwaysReady

	testCases := []*ecsv1.KubernetesCluster{
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase: enum.Running,
			},
		},
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-2",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase: enum.Failed,
			},
		},
	}

	for _, test := range testCases {
		_, err := kubernetesClusterClient.EcsV1().KubernetesClusters("").Create(test)
		if err != nil {
			t.Fatalf("error injecting ecs add: %v", err)
		}

		err = ecsv1Controller.processOperateFinished(test)
		if !assert.Equal(t, nil, err) {
			t.Fatalf("expected: %v but get %v", nil, err)
		}
	}
}

func TestProcessOperateFailed(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	kubernetesClusterClient := ecsfake.NewSimpleClientset()
	informerFactory := informers.NewSharedInformerFactory(kubernetesClusterClient, controller.NoResyncPeriodFunc())
	ecsv1Controller := NewController(fakeClient, kubernetesClusterClient,
		informerFactory.Ecs().V1().KubernetesClusters())
	ecsv1Controller.kubernetesClusterSynced = alwaysReady

	testCases := []*ecsv1.KubernetesCluster{
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase: enum.Running,
			},
		},
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-2",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase: enum.Failed,
			},
		},
	}

	for _, test := range testCases {
		_, err := kubernetesClusterClient.EcsV1().KubernetesClusters("").Create(test)
		if err != nil {
			t.Fatalf("error injecting ecs add: %v", err)
		}

		err = ecsv1Controller.processOperateFailed(test)
		if !assert.Equal(t, nil, err) {
			t.Fatalf("expected: %v but get %v", nil, err)
		}
	}
}

func TestProcessOperateCreating(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	kubernetesClusterClient := ecsfake.NewSimpleClientset()
	informerFactory := informers.NewSharedInformerFactory(kubernetesClusterClient, controller.NoResyncPeriodFunc())
	ecsv1Controller := NewController(fakeClient, kubernetesClusterClient,
		informerFactory.Ecs().V1().KubernetesClusters())
	ecsv1Controller.kubernetesClusterSynced = alwaysReady

	testCases := []*ecsv1.KubernetesCluster{
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase: enum.Running,
			},
		},
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-2",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase: enum.Failed,
			},
		},
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-3",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase: enum.Creating,
			},
		},
	}

	for _, test := range testCases {
		//informerFactory.Ecs().V1().KubernetesClusters().Informer().GetIndexer().Add(test)
		//ecsv1Controller.enqueueKubernetesCluster(test)
		//ecsv1Controller.workqueue.AddRateLimited(test.Namespace + "/" + test.Name)
		_, err := kubernetesClusterClient.EcsV1().KubernetesClusters("").Create(test)
		if err != nil {
			t.Fatalf("error injecting ecs add: %v", err)
		}
		err = ecsv1Controller.processOperateCreating(test)
		if !assert.Equal(t, nil, err) {
			t.Fatalf("expected: %v but get %v", nil, err)
		}
	}
}

func TestProcessOperateNew(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	kubernetesClusterClient := ecsfake.NewSimpleClientset()
	informerFactory := informers.NewSharedInformerFactory(kubernetesClusterClient, controller.NoResyncPeriodFunc())
	ecsv1Controller := NewController(fakeClient, kubernetesClusterClient,
		informerFactory.Ecs().V1().KubernetesClusters())
	ecsv1Controller.kubernetesClusterSynced = alwaysReady

	testCases := []*ecsv1.KubernetesCluster{
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase:  enum.New,
				Reason: "...",
			},
		},
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-2",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase:  enum.New,
				Reason: "",
			},
		},
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-3",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase:  enum.Creating,
				Reason: "",
			},
		},
		&ecsv1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-4",
			},
			Status: ecsv1.KubernetesClusterStatus{
				Phase:  enum.Creating,
				Reason: "...",
			},
		},
	}

	for _, test := range testCases {
		_, err := kubernetesClusterClient.EcsV1().KubernetesClusters("").Create(test)
		if err != nil {
			t.Fatalf("error injecting ecs add: %v", err)
		}
		err = ecsv1Controller.processOperateNew(test)
		if !assert.Equal(t, nil, err) {
			t.Fatalf("expected: %v but get %v", nil, err)
		}
	}
}
