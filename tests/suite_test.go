/*
Copyright 2019-2020 The OpenEBS Authors

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

package tests

import (
	"flag"

	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	ns "github.com/openebs/dynamic-nfs-provisioner/pkg/kubernetes/api/core/v1/namespace"
	"github.com/openebs/maya/tests"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	kubeConfigPath              string
	namespace                   = "nfs-tests-ns"
	namespaceObj                *corev1.Namespace
	err                         error
	NFSProvisionerLabelSelector = "openebs.io/component-name=openebs-nfs-provisioner"
	OpenEBSNamespace            = "openebs"
)

func TestSource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test application deployment")
}

func init() {
	flag.StringVar(&kubeConfigPath, "kubeconfig", os.Getenv("KUBECONFIG"), "path to kubeconfig to invoke kubernetes API calls")
}

var ops *tests.Operations

var _ = BeforeSuite(func() {

	ops = tests.NewOperations(tests.WithKubeConfigPath(kubeConfigPath))

	By("waiting for openebs-nfs-provisioner pod to come into running state")
	provPodCount := ops.GetPodRunningCountEventually(
		string(OpenEBSNamespace),
		string(NFSProvisionerLabelSelector),
		1,
	)
	Expect(provPodCount).To(Equal(1))

	By("building a namespace")
	namespaceObj, err = ns.NewBuilder().
		WithGenerateName(namespace).
		APIObject()
	Expect(err).ShouldNot(HaveOccurred(), "while building namespace {%s}", namespaceObj.GenerateName)

	By("creating above namespace")
	namespaceObj, err = ops.NSClient.Create(namespaceObj)
	Expect(err).To(BeNil(), "while creating namespace {%s}", namespaceObj.GenerateName)

})

var _ = AfterSuite(func() {

	By("deleting namespace")
	err = ops.NSClient.Delete(namespaceObj.Name, &metav1.DeleteOptions{})
	Expect(err).To(BeNil(), "while deleting namespace {%s}", namespaceObj.Name)

})
