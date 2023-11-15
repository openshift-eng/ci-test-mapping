package cloudcomputeotherprovider

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var OtherProviderComponent = Component{
	Component: &config.Component{
		Name:                 "Cloud Compute / Other Provider",
		Operators:            []string{"cluster-api", "machine-approver", "machine-api", "control-plane-machine-set"},
		DefaultJiraComponent: "Cloud Compute / Other Provider",
		Namespaces: []string{
			"openshift-cluster-machine-approver",
			"openshift-machine-api",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAny: []string{
					"bz-Cloud Compute",
					"bz-cluster-api",
					"bz-control-plane-machine-set",
					"Managed cluster should have same number of Machines and Nodes",
				},
			},
			{
				IncludeAll: []string{"service-load-balancer-", "disruption"},
			},
			{
				SIG:        "sig-network-edge",
				IncludeAll: []string{"Application behind service load balancer"},
			},
			{Suite: "Alerting for machine-api"},
			{Suite: "Machine-api components upgrade tests"},
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
