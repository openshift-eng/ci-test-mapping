package networkingopenshiftsdn

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var OpenshiftSdnComponent = Component{
	Component: &config.Component{
		Name:                 "Networking / openshift-sdn",
		Operators:            []string{},
		DefaultJiraComponent: "Networking / openshift-sdn",
		Matchers: []config.ComponentMatcher{
			{
				// Tests that skip a network other than SDN are assumed to belong to us.
				SIG:        "sig-network",
				IncludeAll: []string{"Skipped:Network/"},
				ExcludeAll: []string{"Skipped:Network/OpenShiftSDN"},
			},
			{
				IncludeAny: []string{
					"Build image sdn",
					"Bug 1812261: iptables is segfaulting",
				},
				Priority: 1,
			},
			{Suite: "SDN related networking scenarios"},
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
