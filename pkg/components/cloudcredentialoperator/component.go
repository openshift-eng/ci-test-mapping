package cloudcredentialoperator

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var CloudCredentialOperatorComponent = Component{
	Component: &config.Component{
		Name:                 "Cloud Credential Operator",
		Operators:            []string{"cloud-credential"},
		DefaultJiraComponent: "Cloud Credential Operator",
		Namespaces: []string{
			"openshift-cloud-credential-operator",
		},
		TestRenames: map[string]string{
			"[Cloud Credential Operator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cloud-credential-operator":    "[bz-Cloud Credential Operator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cloud-credential-operator",
			"[Cloud Credential Operator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cloud-credential-operator": "[bz-Cloud Credential Operator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cloud-credential-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Cloud Credential Operator"},
			},
			{Suite: "CloudCredentialOperator components upgrade tests"},
			{Suite: "cloud credential operator"},
			{
				IncludeAny: []string{
					"Cloud credential",
					"upgrade should succeed: cloud-credential",
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
