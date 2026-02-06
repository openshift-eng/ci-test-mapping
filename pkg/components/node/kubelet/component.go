package nodekubelet

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var KubeletComponent = Component{
	Component: &config.Component{
		Name:                 "Node / Kubelet",
		Operators:            []string{},
		DefaultJiraComponent: "Node / Kubelet",
		Matchers: []config.ComponentMatcher{
			{
				SIG:      "sig-node",
				Priority: 1,
			},
			{
				IncludeAny: []string{
					"Node process segfaulted",
					"all nodes should be ready",
					"[sig-arch] [Conformance] sysctl pod should not start for sysctl not on whitelist kernel.msgmax [Suite:openshift/conformance/parallel/minimal]",
					"[sig-arch] [Conformance] sysctl pod should not start for sysctl not on whitelist net.ipv4.ip_dynaddr [Suite:openshift/conformance/parallel/minimal]",
					"[sig-arch] [Conformance] sysctl whitelists kernel.shm_rmid_forced [Suite:openshift/conformance/parallel/minimal]",
					"[sig-arch] [Conformance] sysctl whitelists net.ipv4.ip_local_port_range [Suite:openshift/conformance/parallel/minimal]",
					"[sig-arch] [Conformance] sysctl whitelists net.ipv4.ip_unprivileged_port_start [Suite:openshift/conformance/parallel/minimal]",
					"[sig-arch] [Conformance] sysctl whitelists net.ipv4.ping_group_range [Suite:openshift/conformance/parallel/minimal]",
					"[sig-arch] [Conformance] sysctl whitelists net.ipv4.tcp_syncookies [Suite:openshift/conformance/parallel/minimal]",
					"upgrade should succeed: node",
					"Seccomp part of SCC policy should be kept and working after upgrade",
					"Make sure multiple resources work well after upgrade",
					"OCP-13016",
				},
				Priority: 1,
			},
			{
				IncludeAny: []string{
					":Node ",
				},
				Priority: 1,
			},
			{Suite: "Access to Node logs"},
			{Suite: "Configuration of environment variables check"},
			{Suite: "Downward API"},
			{Suite: "InitContainers"},
			{Suite: "NODE"},
			{Suite: "Node components upgrade tests"},
			{Suite: "Node management"},
			{Suite: "Node related"},
			{Suite: "NodeSelector related tests"},
			{Suite: "Workload Secrets test"},
			{Suite: "configmap related"},
			{Suite: "containers related features"},
			{Suite: "kubelet restart and node restart"},
			{Suite: "node logs related features"},
			{Suite: "pod related features"},
			{Suite: "pods related feature"},
			{Suite: "pods related scenarios"},
			{Suite: "secrets related"},
			{Suite: "secrets related scenarios"},
			{Suite: "Access to Node logs (OCP-43996,admin)"},
			{Suite: "configMap"},
			{Suite: "Permission Data"},
			{Suite: "scenarios related with secret volume"},
		},
		TestRenames: map[string]string{
			"[sig-node] pathological event openshift-config-operator should not get probe error on readiness probe due to timeout":                                                                      "[sig-node] openshift-config-operator should not get probe error on readiness probe due to timeout",
			"[sig-node] pathological event openshift-config-operator should not get probe error on liveness probe due to timeout":                                                                       "[sig-node] openshift-config-operator should not get probe error on liveness probe due to timeout",
			"[sig-node] pathological event openshift-config-operator readiness probe should not fail due to timeout":                                                                                    "[sig-node] openshift-config-operator readiness probe should not fail due to timeout",
			"[sig-node] pathological event NodeHasSufficeintMemory condition does not occur too often":                                                                                                  "[sig-node] Test the NodeHasSufficeintMemory condition does not occur too often",
			"[sig-node] pathological event NodeHasNoDiskPressure condition does not occur too often":                                                                                                    "[sig-node] Test the NodeHasNoDiskPressure condition does not occur too often",
			"[sig-node] pathological event NodeHasSufficientPID condition does not occur too often":                                                                                                     "[sig-node] Test the NodeHasSufficientPID condition does not occur too often",
			"[sig-builds][sig-node][Feature:Builds][apigroup:build.openshift.io] zstd:chunked Image should successfully run date command [Skipped:Disconnected] [Suite:openshift/conformance/parallel]": "[sig-node][Feature:Builds][apigroup:build.openshift.io] zstd:chunked Image should successfully run date command [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
			"[sig-node] zstd:chunked Image should successfully run date command [Skipped:Disconnected] [Suite:openshift/conformance/parallel]":                                                          "[sig-builds][sig-node][Feature:Builds][apigroup:build.openshift.io] zstd:chunked Image should successfully run date command [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
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
