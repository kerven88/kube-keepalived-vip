/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8s

import (
	"os"
	"testing"

	api "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestGetNodeIP(t *testing.T) {
	fKNodes := []struct {
		cs *testclient.Clientset
		n  string
		ea string
	}{
		// empty node list
		{testclient.NewSimpleClientset(), "demo", ""},

		// node not exist
		{testclient.NewSimpleClientset(&api.NodeList{Items: []api.Node{{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: "demo",
			},
			Status: api.NodeStatus{
				Addresses: []api.NodeAddress{
					{
						Type:    api.NodeInternalIP,
						Address: "10.0.0.1",
					},
				},
			},
		}}}), "notexistnode", ""},

		// node  exist
		{testclient.NewSimpleClientset(&api.NodeList{Items: []api.Node{{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: "demo",
			},
			Status: api.NodeStatus{
				Addresses: []api.NodeAddress{
					{
						Type:    api.NodeInternalIP,
						Address: "10.0.0.1",
					},
				},
			},
		}}}), "demo", "10.0.0.1"},

		// search the correct node
		{testclient.NewSimpleClientset(&api.NodeList{Items: []api.Node{
			{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "demo1",
				},
				Status: api.NodeStatus{
					Addresses: []api.NodeAddress{
						{
							Type:    api.NodeInternalIP,
							Address: "10.0.0.1",
						},
					},
				},
			},
			{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "demo2",
				},
				Status: api.NodeStatus{
					Addresses: []api.NodeAddress{
						{
							Type:    api.NodeInternalIP,
							Address: "10.0.0.2",
						},
					},
				},
			},
		}}), "demo2", "10.0.0.2"},

		// get NodeExternalIP
		{testclient.NewSimpleClientset(&api.NodeList{Items: []api.Node{{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: "demo",
			},
			Status: api.NodeStatus{
				Addresses: []api.NodeAddress{
					{
						Type:    api.NodeInternalIP,
						Address: "10.0.0.1",
					}, {
						Type:    api.NodeExternalIP,
						Address: "10.0.0.2",
					},
				},
			},
		}}}), "demo", "10.0.0.2"},

		// get NodeInternalIP
		{testclient.NewSimpleClientset(&api.NodeList{Items: []api.Node{{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: "demo",
			},
			Status: api.NodeStatus{
				Addresses: []api.NodeAddress{
					{
						Type:    api.NodeExternalIP,
						Address: "",
					}, {
						Type:    api.NodeInternalIP,
						Address: "10.0.0.2",
					},
				},
			},
		}}}), "demo", "10.0.0.2"},
	}

	for _, fk := range fKNodes {
		address := GetNodeIP(fk.cs, fk.n)
		if address != fk.ea {
			t.Errorf("expected %s, but returned %s", fk.ea, address)
		}
	}
}

func TestGetPodDetails(t *testing.T) {
	// POD_NAME & POD_NAMESPACE not exist
	os.Setenv("POD_NAME", "")
	os.Setenv("POD_NAMESPACE", "")
	_, err1 := GetPodDetails(testclient.NewSimpleClientset())
	if err1 == nil {
		t.Errorf("expected an error but returned nil")
	}

	// POD_NAME not exist
	os.Setenv("POD_NAME", "")
	os.Setenv("POD_NAMESPACE", api.NamespaceDefault)
	_, err2 := GetPodDetails(testclient.NewSimpleClientset())
	if err2 == nil {
		t.Errorf("expected an error but returned nil")
	}

	// POD_NAMESPACE not exist
	os.Setenv("POD_NAME", "testpod")
	os.Setenv("POD_NAMESPACE", "")
	_, err3 := GetPodDetails(testclient.NewSimpleClientset())
	if err3 == nil {
		t.Errorf("expected an error but returned nil")
	}

	// POD not exist
	os.Setenv("POD_NAME", "testpod")
	os.Setenv("POD_NAMESPACE", api.NamespaceDefault)
	_, err4 := GetPodDetails(testclient.NewSimpleClientset())
	if err4 == nil {
		t.Errorf("expected an error but returned nil")
	}

	// success to get PodInfo
	fkClient := testclient.NewSimpleClientset(
		&api.PodList{Items: []api.Pod{{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "testpod",
				Namespace: api.NamespaceDefault,
				Labels: map[string]string{
					"first":  "first_label",
					"second": "second_label",
				},
			},
		}}},
		&api.NodeList{Items: []api.Node{{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: "demo",
			},
			Status: api.NodeStatus{
				Addresses: []api.NodeAddress{
					{
						Type:    api.NodeInternalIP,
						Address: "10.0.0.1",
					},
				},
			},
		}}})

	epi, err5 := GetPodDetails(fkClient)
	if err5 != nil {
		t.Errorf("expected a PodInfo but returned error")
		return
	}

	if epi == nil {
		t.Errorf("expected a PodInfo but returned nil")
	}
}
