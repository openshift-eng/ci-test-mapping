package clusterversionoperator

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ClusterVersionOperatorComponent = Component{
	Component: &config.Component{
		Name:                 "Cluster Version Operator",
		Operators:            []string{"version"},
		DefaultJiraComponent: "Cluster Version Operator",
		Namespaces: []string{
			"openshift-cluster-version",
		},
		TestRenames: map[string]string{
			"[sig-cluster-lifecycle] cluster upgrade should complete in 105.00 minutes":                                                     "[sig-cluster-lifecycle] cluster upgrade should complete in a reasonable time",
			"[sig-cluster-lifecycle] cluster upgrade should complete in 120.00 minutes":                                                     "[sig-cluster-lifecycle] cluster upgrade should complete in a reasonable time",
			"[sig-cluster-lifecycle] cluster upgrade should complete in 210.00 minutes":                                                     "[sig-cluster-lifecycle] cluster upgrade should complete in a reasonable time",
			"[sig-cluster-lifecycle] cluster upgrade should complete in 240.00 minutes":                                                     "[sig-cluster-lifecycle] cluster upgrade should complete in a reasonable time",
			"[sig-cluster-lifecycle] cluster upgrade should complete in 75.00 minutes":                                                      "[sig-cluster-lifecycle] cluster upgrade should complete in a reasonable time",
			"[sig-cluster-lifecycle] cluster upgrade should complete in 90.00 minutes":                                                      "[sig-cluster-lifecycle] cluster upgrade should complete in a reasonable time",
			"[Cluster Version Operator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cluster-version":    "[bz-Cluster Version Operator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cluster-version",
			"[Cluster Version Operator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cluster-version": "[bz-Cluster Version Operator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cluster-version",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAny: []string{
					"cluster-version-operator",
					"bz-Cluster Version Operator",
				},
			},
			{
				IncludeAll: []string{"upgrade"},
				// There are other component's cases with `upgrade` substring in the name, excluded them
				ExcludeAny: []string{
					"Cluster_Infrastructure",
					"Cloud credential",
					"OCP-13016",
					"OCP-22612",
					"OCP-41804",
					"OCP-64657",
					"OAP cert-manager",
					"samples/openshift-controller-manager/image-registry operators",
					"etcd-operator and cluster works well after upgrade",
					"Seccomp part of SCC policy should be kept and working after upgrade",
					"Make sure multiple resources work well after upgrade",
					"upgrade should succeed: rhel",
					"upgrade should succeed: node",
					"upgrade should succeed: authentication",
					"upgrade should succeed: baremetal",
					"upgrade should succeed: cloud-controller-manager",
					"upgrade should succeed: cloud-credential",
					"upgrade should succeed: cluster-autoscaler",
					"upgrade should succeed: config-operator",
					"upgrade should succeed: console",
					"upgrade should succeed: control-plane-machine-set",
					"upgrade should succeed: csi-snapshot-controller",
					"upgrade should succeed: dns",
					"upgrade should succeed: etcd",
					"upgrade should succeed: image-registry",
					"upgrade should succeed: ingress",
					"upgrade should succeed: insights",
					"upgrade should succeed: kube-apiserver",
					"upgrade should succeed: kube-controller-manager",
					"upgrade should succeed: kube-scheduler",
					"upgrade should succeed: kube-storage-version-migrator",
					"upgrade should succeed: machine-api",
					"upgrade should succeed: machine-approver",
					"upgrade should succeed: machine-config",
					"upgrade should succeed: marketplace",
					"upgrade should succeed: monitoring",
					"upgrade should succeed: network",
					"upgrade should succeed: node-tuning",
					"upgrade should succeed: olm",
					"upgrade should succeed: openshift-apiserver",
					"upgrade should succeed: openshift-controller-manager",
					"upgrade should succeed: openshift-samples",
					"upgrade should succeed: operator-lifecycle-manager",
					"upgrade should succeed: operator-lifecycle-manager-catalog",
					"upgrade should succeed: operator-lifecycle-manager-packageserver",
					"upgrade should succeed: service-ca",
					"upgrade should succeed: storage",
					"upgrade should succeed:  ",
				},
				// Let others claim upgrade tests (i.e. for their component)
				Priority: -10,
			},
			{
				// all cvo QE cases include "OTA cvo should"
				// the junit report of Prow upgrade ci is like pattern "upgrade should succeed: $UPGRADE_FAILURE_TYPE"
				IncludeAny: []string{
					"OTA cvo should",
					"upgrade should succeed: overall",
					"upgrade should succeed: cvo",
					"upgrade should succeed: rollback",
					"upgrade should succeed: admin_ack",
				},
			},
			{
				IncludeAny: []string{
					"[sig-arch] ClusterOperators [apigroup:config.openshift.io] should define at least one namespace in their lists of related objects  [Suite:openshift/conformance/parallel]",
					"[sig-arch] ClusterOperators [apigroup:config.openshift.io] should define at least one namespace in their lists of related objects [Suite:openshift/conformance/parallel]",
					"[sig-arch] ClusterOperators [apigroup:config.openshift.io] should define at least one related object that is not a namespace [Suite:openshift/conformance/parallel]",
					"[sig-arch] ClusterOperators [apigroup:config.openshift.io] should define valid related objects [Suite:openshift/conformance/parallel]",
					"[sig-arch] ClusterOperators should define at least one namespace in their lists of related objects [Suite:openshift/conformance/parallel]",
					"[sig-arch] ClusterOperators should define at least one related object that is not a namespace [Suite:openshift/conformance/parallel]",
					"[sig-arch] ClusterOperators should define valid related objects [Suite:openshift/conformance/parallel]",
					"[sig-cluster-lifecycle] TestAdminAck should succeed [apigroup:config.openshift.io] [Suite:openshift/conformance/parallel]",
					"[sig-cluster-lifecycle] TestAdminAck should succeed [Suite:openshift/conformance/parallel]",
				},
			},
		},
	},
}

func (c *Component) IdentifyTest(test *v1.TestInfo) (*v1.TestOwnership, error) {
	if matcher := c.FindMatch(test); matcher != nil {
		jira := matcher.JiraComponent
		if jira == "" {
			jira = c.DefaultJiraComponent
		}
		return &v1.TestOwnership{
			Name:          test.Name,
			Component:     c.Name,
			JIRAComponent: jira,
			Priority:      matcher.Priority,
			Capabilities:  append(matcher.Capabilities, identifyCapabilities(test)...),
		}, nil
	}

	return nil, nil
}

func (c *Component) StableID(test *v1.TestInfo) string {
	// Look up the stable name for our test in our renamed tests map.
	if stableName, ok := c.TestRenames[test.Name]; ok {
		return stableName
	}
	return test.Name
}

func (c *Component) JiraComponents() (components []string) {
	components = []string{c.DefaultJiraComponent}
	for _, m := range c.Matchers {
		components = append(components, m.JiraComponent)
	}

	return components
}
