package lpquay

import (
	"regexp"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var LPquayComponent = Component{
	Component: &config.Component{
		Name:                 "LP--Quay",
		Operators:            []string{},
		DefaultJiraComponent: "LP--Quay",
		Matchers: []config.ComponentMatcher{
			{SuiteRegEx: regexp.MustCompile(`^lp-ocp-compat--Quay--`)},
			{SuiteRegEx: regexp.MustCompile(`^lp-chaos--Quay--`)},
			//  LP Interop Test Suites: Capability is used to group Runtime Env.
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--Quay--OPP--`), Capabilities: []string{"OPP"}},
			// Platform-infra smoke suite, routed to the existing Quay component
			// (INTEROP-9630): both the runtime-first lp-interop--OPP--smoke and
			// the canonical lp-interop--OCP--OPP--smoke forms carry the
			// quay-pull case, so the claim is scoped by suite and test name
			// and does not collide with the canonical OPP matchers above.
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--(OCP--)?OPP--smoke$`), IncludeAny: []string{"quay-pull"}},
			// Runtime-first emission from interop-opp-product-upgrade-quay, which
			// still writes the literal suite name in openshift/release. The
			// exact Suite claim outranks Cluster Version Operator's Priority -10
			// name-based "upgrade" match, so the live quay-upgrade row keeps Quay
			// ownership. Drop once that step migrates (INTEROP-9629).
			{Suite: "lp-interop--OPP--quay-upgrade", Capabilities: []string{"OPP"}},
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
