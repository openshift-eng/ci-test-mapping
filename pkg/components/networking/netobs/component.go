package networkingnetobs

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var NetObsComponent = Component{
	Component: &config.Component{
		Name:                 "Networking / NetObs",
		Operators:            []string{},
		DefaultJiraComponent: "Networking / NetObs",
		Matchers: []config.ComponentMatcher{
			{Suite: `"(OCP-50532, OCP-50531, OCP-50530, OCP-59408 NETOBSERV) Netflow Table view tests"`},
			{Suite: `(OCP-50532, OCP-50531, OCP-50530 NETOBSERV) Netflow Table view tests`},
			{Suite: `(OCP-53591 NETOBSERV) Netflow Topology view features`},
			{Suite: `(OCP-60701 NETOBSERV) Connection tracking test`},
			{Suite: `(OCP-66141 NETOBSERV) PacketDrop test`},
			{Suite: `(OCP-54839 NETOBSERV) Netflow Overview page tests`},
			{Suite: `(OCP-56222 NETOBSERV) Quick Filters test`},
			{Suite: `(OCP-50532, OCP-50531, OCP-50530, OCP-59408 NETOBSERV) Netflow Table view tests`},
			{Suite: `(OCP-61893 NETOBSERV) NetObserv dashboards tests`},
			{Suite: "Flows Tracking and Monitoring for Network Analytics"},
			{Suite: "Network_Observability"},
			{Suite: "netflow table page features"},
			{Suite: "groups, edges, labels, badges"},
			{Suite: "PacketDrop features"},
			{Suite: "NetObserv dashboards tests"},
			{Suite: "Connection tracking netflow table page features"},
			{Suite: "NETOBSERV Performances"},
			{Suite: "NETOBSERV dashboards tests"},
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
