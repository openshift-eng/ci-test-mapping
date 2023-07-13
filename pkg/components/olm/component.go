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
		Operators:            []string{"marketplace", "operator-lifecycle-manager", "operator-lifecycle-manager-catalog", "operator-lifecycle-manager-packageserver"},
		DefaultJiraComponent: "OLM",
		Matchers: []config.ComponentMatcher{
			{
				IncludeAny: []string{
					"bz-OLM",
					"bz-platform-operators-aggregated",
				},
			},
			{
				Namespaces: []string{
					"openshift-marketplace",
					"openshift-operator-lifecycle-manager",
					"openshift-operators",
				},
			},
			{
				SIG: "sig-operator",
			},
			{
				IncludeAny: []string{
					"[sig-arch] ocp payload should be based on existing source OLM version should contain the source commit id [apigroup:config.openshift.io] [Suite:openshift/conformance/parallel]",
					"[sig-arch] ocp payload should be based on existing source OLM version should contain the source commit id [Suite:openshift/conformance/parallel]",
					"[sig-arch] ocp payload should be based on existing source [Serial] olm version should contain the source commit id [Suite:openshift/conformance/serial]",
					"[sig-arch] openshift-marketplace pods should not get excessive startupProbe failures",
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
