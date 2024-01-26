package networkingovnkubernetes

import (
	"regexp"

	log "github.com/sirupsen/logrus"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

// A regular expression + its replacement string to apply to test names
type testNameReplacement struct {
	regexp      *regexp.Regexp
	replacement string
}

type Component struct {
	*config.Component
	testNameReplacements []testNameReplacement
}

var OvnKubernetesComponent = Component{
	Component: &config.Component{
		Name:                 "Networking / ovn-kubernetes",
		Operators:            []string{},
		DefaultJiraComponent: "Networking / ovn-kubernetes",
		Namespaces: []string{
			"openshift-ovn-kubernetes",
		},
		Matchers: []config.ComponentMatcher{
			{
				SIG: "sig-network",
				// Tests that skip a network other than OVN are assumed to belong to us.
				IncludeAll: []string{"Skipped:Network/"},
				ExcludeAny: []string{"Skipped:Network/OVNKubernetes", "Skipped:Network/OVNKuberenetes"},
			},
			{
				IncludeAll: []string{"ovn-kubernetes"},
				Priority:   1,
			},
		},
		TestRenames: map[string]string{
			"[Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ovn-kubernetes":    "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ovn-kubernetes",
			"[Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ovn-kubernetes": "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ovn-kubernetes",
		},
	},
	testNameReplacements: []testNameReplacement{
		{
			regexp:      regexp.MustCompile("avalibility"),
			replacement: "availability",
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
	// Apply any test name replacements
	name := test.Name
	for _, r := range c.testNameReplacements {
		name = r.regexp.ReplaceAllString(name, r.replacement)
	}

	// Look up the stable name for our test in our renamed tests map.
	if stableName, ok := c.TestRenames[name]; ok {
		name = stableName
	}
	if name != test.Name {
		log.Tracef("Mapped test %q to %q", test.Name, name)
	}
	return name
}

func (c *Component) JiraComponents() (components []string) {
	components = []string{c.DefaultJiraComponent}
	for _, m := range c.Matchers {
		components = append(components, m.JiraComponent)
	}

	return components
}
