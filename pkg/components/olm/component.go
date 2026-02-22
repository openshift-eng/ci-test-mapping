package olm

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var OLMComponent = Component{
	Component: &config.Component{
		Name:                 "OLM",
		Operators:            []string{"olm", "marketplace", "operator-lifecycle-manager", "operator-lifecycle-manager-catalog", "operator-lifecycle-manager-packageserver"},
		DefaultJiraComponent: "OLM",
		Namespaces: []string{
			"openshift-marketplace",
			"openshift-operator-lifecycle-manager",
			"openshift-operators",
			"openshift-catalogd",
			"openshift-operator-controller",
			"openshift-cluster-olm-operator",
			"openshift-rukpak",
		},
		Matchers: []config.ComponentMatcher{
			{
				SIG: "sig-olmv1",
			},
			{
				IncludeAny: []string{
					"bz-OLM",
					"bz-platform-operators-aggregated",
				},
			},
			{
				SIG:      "sig-operator",
				Priority: 1,
			},
			{
				IncludeAny: []string{
					"[sig-arch] ocp payload should be based on existing source OLM version should contain the source commit id [apigroup:config.openshift.io] [Suite:openshift/conformance/parallel]",
					"[sig-arch] ocp payload should be based on existing source OLM version should contain the source commit id [Suite:openshift/conformance/parallel]",
					"[sig-arch] ocp payload should be based on existing source [Serial] olm version should contain the source commit id [Suite:openshift/conformance/serial]",
					"[sig-arch] openshift-marketplace pods should not get excessive startupProbe failures",
					"upgrade should succeed: olm",
					"upgrade should succeed: operator-lifecycle-manager",
					"upgrade should succeed: operator-lifecycle-manager-catalog",
					"upgrade should succeed: operator-lifecycle-manager-packageserver",
					"upgrade should succeed: marketplace",
				},
			},
			{Suite: "Marketplace related scenarios"},
			{Suite: "OLM"},
		},
		TestRenames: map[string]string{
			"[OLM][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-marketplace":                   "[bz-OLM][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-marketplace",
			"[OLM][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-operator-lifecycle-manager":    "[bz-OLM][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-operator-lifecycle-manager",
			"[OLM][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-operators":                     "[bz-OLM][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-operators",
			"[OLM][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-marketplace":                "[bz-OLM][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-marketplace",
			"[OLM][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-operator-lifecycle-manager": "[bz-OLM][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-operator-lifecycle-manager",
			"[OLM][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-operators":                  "[bz-OLM][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-operators",

			"[sig-scheduling][Early] The openshift-operator-lifecycle-manager pods [apigroup:packages.operators.coreos.com] should be scheduled on different nodes [Skipped:SingleReplicaTopology] [Suite:openshift/conformance/parallel]": "[sig-scheduling][Early] The openshift-operator-lifecycle-manager pods [apigroup:packages.operators.coreos.com] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
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
