package tekton_bundle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kubevirt/kubevirt-tekton-tasks-operator/pkg/operands"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	numberOfClusterTasks    = 11
	numberOfServiceAccounts = 9
	numberOfRoleBindings    = 9
	numberOfClusterRoles    = 9
	numberOfPipelines       = 2
	numberOfConfigMaps      = 2
)

var _ = Describe("Tekton bundle", func() {

	It("should return correct pipeline folder path on okd", func() {
		path := getPipelineBundlePath(true)
		Expect(path).To(Equal("/data/tekton-pipelines/okd/"))
	})

	It("should return correct pipeline folder path on kubernetes", func() {
		path := getPipelineBundlePath(false)
		Expect(path).To(Equal("/data/tekton-pipelines/kubernetes/"))
	})

	It("should return correct task path on okd", func() {
		path := getTasksBundlePath(true)
		Expect(path).To(Equal("/data/tekton-tasks/okd/kubevirt-tekton-tasks-okd-" + operands.TektonTasksVersion + ".yaml"))
	})

	It("should return correct task path on kubernetes", func() {
		path := getTasksBundlePath(false)
		Expect(path).To(Equal("/data/tekton-tasks/kubernetes/kubevirt-tekton-tasks-kubernetes-" + operands.TektonTasksVersion + ".yaml"))
	})

	It("should load correct files and convert them", func() {
		path, _ := os.Getwd()

		taskPath := filepath.Join(path, "test-bundle-files/test-tasks/test-tasks.yaml")
		pipelinePath := filepath.Join(path, "test-bundle-files/test-pipelines/")

		tektonObjs, err := decodeObjectsFromFiles(taskPath, pipelinePath)
		Expect(err).ToNot(HaveOccurred(), "it should not throw error")
		Expect(tektonObjs.ClusterTasks).To(HaveLen(numberOfClusterTasks), "number of tasks should equal")
		Expect(tektonObjs.ServiceAccounts).To(HaveLen(numberOfServiceAccounts), "number of service accounts should equal")
		Expect(tektonObjs.RoleBindings).To(HaveLen(numberOfRoleBindings), "number of role bindings should equal")
		Expect(tektonObjs.ClusterRoles).To(HaveLen(numberOfClusterRoles), "number of cluster roles should equal")
		Expect(tektonObjs.Pipelines).To(HaveLen(numberOfPipelines), "number of pipelines should equal")
		Expect(tektonObjs.ConfigMaps).To(HaveLen(numberOfConfigMaps), "number of config maps should equal")
	})
})

func TestTektonBundle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tekton Bundle Suite")
}
