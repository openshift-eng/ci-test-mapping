package managementconsole

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ManagementConsoleComponent = Component{
	Component: &config.Component{
		Name:                 "Management Console",
		Operators:            []string{"console-operator", "console"},
		DefaultJiraComponent: "Management Console",
		Namespaces: []string{
			"openshift-console",
			"openshift-console-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Management Console"},
			},
			{Suite: "Administration pages pesudo translation"},
			{Suite: "Branding check"},
			{Suite: "Debug console for pods"},
			{Suite: "Dynamic Plugins notification features"},
			{Suite: "Dynamic plugins features"},
			{Suite: "Dynamic provisioning"},
			{Suite: "Home related pages via admin console"},
			{Suite: "Notification drawer tests"},
			{Suite: "Operators Installed page test"},
			{Suite: "PDB List Page and Detail Page Test"},
			{Suite: "Projects dropdown tests"},
			{Suite: "about cluster setting page"},
			{Suite: "add idp from console"},
			{Suite: "admin console api related"},
			{Suite: "alerts browser"},
			{Suite: "command line tools page"},
			{Suite: "console configs features"},
			{Suite: "console feature on sno cluster"},
			{Suite: "console style related"},
			{Suite: "console-operator related"},
			{Suite: "console-route"},
			{Suite: "customize console related"},
			{Suite: "dark-theme related feature"},
			{Suite: "dashboards related cases"},
			{Suite: "mega menu on console"},
			{Suite: "namespace dropdown favorite test"},
			{Suite: "operand form view"},
			{Suite: "pod page"},
			{Suite: "projects related features via web"},
			{Suite: "query browser"},
			{Suite: "search page related"},
			{Suite: "tests on catalog page"},
			{Suite: "web console related upgrade check"},
			{Suite: "show shortname in console resourse badge"},
			{Suite: "user preferences related features"},
			{Suite: "overview page features"},
			{Suite: "Administration pages pesudo translation (OCP-35766,admin)"},
			{Suite: "overview cases"},
			{Suite: "environment related"},
			{Suite: "kibana web UI related cases for logging"},
			{Suite: "deployment page"},
			{Suite: "knmstate operator console plugin related features"},
			{Suite: "masthead related"},
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
