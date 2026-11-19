package lpocp

import (
	"regexp"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var LPocpComponent = Component{
	Component: &config.Component{
		Name:                 "LP--OCP",
		Operators:            []string{},
		DefaultJiraComponent: "LP--OCP",
		Matchers: []config.ComponentMatcher{
			{SuiteRegEx: regexp.MustCompile(`^lp-ocp-compat--OCP--`)},
			{SuiteRegEx: regexp.MustCompile(`^lp-chaos--OCP--`)},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--OPP--`), Capabilities: []string{"OPP"}},
			// Runtime-first exact Suite retained until producers migrate to
			// canonical lp-interop--OCP--OPP--* names (INTEROP-9629 migration):
			// 5 live rows currently Unknown.
			{Suite: "lp-interop--OPP--smoke", Capabilities: []string{"OPP"}},
			// Runtime-first exact Suite retained until producers migrate to
			// canonical lp-interop--OCP--OPP--* names (INTEROP-9629 migration):
			// 1 live row currently Unknown.
			{Suite: "lp-interop--OPP--backup", Capabilities: []string{"OPP"}},
			// Runtime-first exact Suite retained until producers migrate to
			// canonical lp-interop--OCP--OPP--* names (INTEROP-9629 migration):
			// 2 live rows currently Unknown.
			{Suite: "lp-interop--OPP--wait-mcp", Capabilities: []string{"OPP"}},
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
