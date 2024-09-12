package obsoletetests

import (
	"k8s.io/apimachinery/pkg/util/sets"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
)

type OCPObsoleteTestManager struct{}

type obsoleteTestIdentifier struct {
	name  string
	suite string
}

var obsoleteTests = sets.New[obsoleteTestIdentifier](
	[]obsoleteTestIdentifier{
		// Removed in alert refactor by TRT https://github.com/openshift/origin/pull/28332
		{
			name:  "[sig-arch] Check if alerts are firing during or after upgrade success",
			suite: "Cluster upgrade",
		},
		// The test has been removed in cucushift, and migrate to ginkgo with different format
		// https://github.com/openshift/cucushift/pull/9640
		{
			name:  "OCP-12158:APIServer Specify ResourceQuota on project",
			suite: "remote registry related scenarios",
		},
		// These tests are unreliable so we're obsoleting them temporarily while they get fixed.
		// They are created by a script in gather-extra, but this is unreliable:
		// 	https://github.com/openshift/release/blob/4a07b67554d760fc75db4431f38094e82d60d57b/ci-operator/step-registry/gather/extra/gather-extra-commands.sh#L646
		//
		//		1. Some jobs skip gather on success, so these tests only get created on job failure
		//		2. Some jobs have multiple post steps, and execution stops if an earlier step fails
		//
		{name: "operator conditions authentication", suite: "Operator results"},
		{name: "operator conditions baremetal", suite: "Operator results"},
		{name: "operator conditions cloud-controller-manager", suite: "Operator results"},
		{name: "operator conditions cloud-credential", suite: "Operator results"},
		{name: "operator conditions cluster-api", suite: "Operator results"},
		{name: "operator conditions cluster-autoscaler", suite: "Operator results"},
		{name: "operator conditions config-operator", suite: "Operator results"},
		{name: "operator conditions console", suite: "Operator results"},
		{name: "operator conditions control-plane-machine-set", suite: "Operator results"},
		{name: "operator conditions csi-snapshot-controller", suite: "Operator results"},
		{name: "operator conditions dns", suite: "Operator results"},
		{name: "operator conditions etcd", suite: "Operator results"},
		{name: "operator conditions image-registry", suite: "Operator results"},
		{name: "operator conditions ingress", suite: "Operator results"},
		{name: "operator conditions insights", suite: "Operator results"},
		{name: "operator conditions kube-apiserver", suite: "Operator results"},
		{name: "operator conditions kube-controller-manager", suite: "Operator results"},
		{name: "operator conditions kube-scheduler", suite: "Operator results"},
		{name: "operator conditions kube-storage-version-migrator", suite: "Operator results"},
		{name: "operator conditions machine-api", suite: "Operator results"},
		{name: "operator conditions machine-approver", suite: "Operator results"},
		{name: "operator conditions machine-config", suite: "Operator results"},
		{name: "operator conditions marketplace", suite: "Operator results"},
		{name: "operator conditions monitoring", suite: "Operator results"},
		{name: "operator conditions network", suite: "Operator results"},
		{name: "operator conditions node-tuning", suite: "Operator results"},
		{name: "operator conditions olm", suite: "Operator results"},
		{name: "operator conditions openshift-apiserver", suite: "Operator results"},
		{name: "operator conditions openshift-controller-manager", suite: "Operator results"},
		{name: "operator conditions openshift-samples", suite: "Operator results"},
		{name: "operator conditions operator-lifecycle-manager", suite: "Operator results"},
		{name: "operator conditions operator-lifecycle-manager-catalog", suite: "Operator results"},
		{name: "operator conditions operator-lifecycle-manager-packageserver", suite: "Operator results"},
		{name: "operator conditions service-ca", suite: "Operator results"},
		{name: "operator conditions storage", suite: "Operator results"},
	}...)

func (*OCPObsoleteTestManager) IsObsolete(test *v1.TestInfo) bool {
	return obsoleteTests.Has(obsoleteTestIdentifier{name: test.Name, suite: test.Suite})
}
