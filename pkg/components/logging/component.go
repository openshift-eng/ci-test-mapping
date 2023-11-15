package logging

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var LoggingComponent = Component{
	Component: &config.Component{
		Name:                 "Logging",
		Operators:            []string{},
		DefaultJiraComponent: "Logging",
		Matchers: []config.ComponentMatcher{
			{Suite: "Elasticsearch related tests"},
			{Suite: "Kibana related features"},
			{Suite: "LOGGING"},
			{Suite: "Logging"},
			{Suite: "Logging related features"},
			{Suite: "Logging smoke test case"},
			{Suite: "Logging view on the openshift console"},
			{Suite: "cluster log forwarder testing"},
			{Suite: "cluster logging related scenarios"},
			{Suite: "cluster-logging-operator related cases"},
			{Suite: "cluster-logging-operator related test"},
			{Suite: "elasticsearch operator related tests"},
			{Suite: "elasticsearch related tests"},
			{Suite: "elasticsearch-operator related tests"},
			{Suite: "logs related features"},
			{Suite: "Cases to test forward logs to external elasticsearch"},
			{Suite: "cluster log forwarder features"},
			{Suite: "fluentd related tests"},
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
