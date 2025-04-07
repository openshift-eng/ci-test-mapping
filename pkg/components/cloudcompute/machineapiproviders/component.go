package cloudcomputemachineapiproviders

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var MachineAPIProvidersComponent = Component{
	Component: &config.Component{
		Name:                 "Cloud Compute / Machine API Providers",
		Operators:            []string{},
		DefaultJiraComponent: "Cloud Compute / Machine API Providers",
		Matchers: []config.ComponentMatcher{
			{
				IncludeAny: []string{
					"[sig-cluster-lifecycle] Cluster_Infrastructure MAPI",
					"[sig-cluster-lifecycle] Cluster_Infrastructure Upgrade",
					"upgrade should succeed: machine-api",
				},
				Priority: 1,
			},
			{Suite: "Alerting for machine-api"},
			{Suite: "Machine-api components upgrade tests"},
			{Suite: "UPI GCP Tests"},
			{Suite: "Machine misc features testing"},
			{Suite: "Machine features testing"},
			{Suite: "AWS machine specific features testing"},
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
