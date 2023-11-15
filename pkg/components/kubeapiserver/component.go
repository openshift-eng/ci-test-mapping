package kubeapiserver

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var KubeApiserverComponent = Component{
	Component: &config.Component{
		Name:                 "kube-apiserver",
		Operators:            []string{"kube-apiserver"},
		DefaultJiraComponent: "kube-apiserver",
		Namespaces: []string{
			"default",
			"openshift",
			"kube-system",
			"openshift-config",
			"openshift-config-managed",
			"openshift-kube-apiserver",
			"openshift-kube-apiserver-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-kube-apiserver"},
			},
			{
				IncludeAll: []string{"cache-kube-api-"},
			},
			{
				IncludeAll: []string{"kube-api-", "-connections"},
			},
			{
				IncludeAll: []string{"[sig-api-machinery][Feature:APIServer]"},
			},
			{
				SIG:      "sig-api-machinery",
				Priority: -1,
			},
			{Suite: "API_Server"},
			{Suite: "CRD related"},
			{Suite: "Event related scenarios"},
			{Suite: "Roles and RoleBindings tests"},
			{Suite: "ServiceAccount and Policy Managerment"},
			{Suite: "events and logs related"},
			{Suite: "limit range related scenarios:"},
			{Suite: "test master config related steps"},
			{Suite: "REST policy related features"},
			{Suite: "senarios for checking transfer scheme"},
			{Suite: "Api proxy related cases"},
			{Suite: "REST related features"},
			{Suite: "REST features"},
			{Suite: "KUBE API server related features"},
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
