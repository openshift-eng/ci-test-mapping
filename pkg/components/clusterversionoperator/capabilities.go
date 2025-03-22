package clusterversionoperator

import (
	"regexp"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/util"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	// ClusterUpgrade is a CVO capability to upgrade the cluster
	ClusterUpgrade = "ClusterUpgrade"
	// ClusterOperators is a CVO capability to manage and monitor Cluster Operators
	ClusterOperators = "ClusterOperators"
	// Operator is a CVO capability to operate the cluster both during and outside of upgrades
	Operator = "Operator"
	// AdminAck is a CVO capability to process administrator acknowledgements
	AdminAck = "AdminAck"
)

var cvoCapabilitiesIdentifiers = map[*regexp.Regexp]string{
	regexp.MustCompile(`.*upgrade.*`): ClusterUpgrade,

	// The junit report of Prow upgrade ci is like pattern "upgrade should succeed: $UPGRADE_FAILURE_TYPE"
	regexp.MustCompile(`.*upgrade should succeed: (overall|cvo|rollback).*`): ClusterUpgrade,

	// e.g. [sig-cluster-lifecycle] pathological event should not see excessive Back-off restarting failed containers for ns/openshift-cluster-version
	regexp.MustCompile(".*ns/openshift-cluster-version.*"): Operator,
	// all invariant tests
	// e.g. [Jira:"Cluster Version Operator"] monitor test legacy-cvo-invariants collection
	regexp.MustCompile(".*monitor test.*"): Operator,
	// all QE's cvo e2e cases from ginkgo include "OTA cvo should"
	regexp.MustCompile(".*OTA cvo should.*"): Operator,

	// e.g. Cluster upgrade.[sig-cluster-lifecycle] ClusterOperators are available and not degraded after upgrade
	regexp.MustCompile(".*ClusterOperators.*"): ClusterOperators,

	regexp.MustCompile(`.*(admin ack|AdminAck).*`): AdminAck,

	regexp.MustCompile(".*upgrade should succeed: admin_ack.*"): AdminAck,
}

func identifyCapabilities(test *v1.TestInfo) []string {
	capabilities := sets.New[string](util.DefaultCapabilities(test)...)
	for matcher, capability := range cvoCapabilitiesIdentifiers {
		if matcher.MatchString(test.Name) {
			capabilities.Insert(capability)
		}
	}

	return capabilities.UnsortedList()
}
