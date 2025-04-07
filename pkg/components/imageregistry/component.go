package imageregistry

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ImageRegistryComponent = Component{
	Component: &config.Component{
		Name:                 "Image Registry",
		Operators:            []string{"image-registry"},
		DefaultJiraComponent: "Image Registry",
		Namespaces: []string{
			"openshift-image-registry",
		},
		TestRenames: map[string]string{
			"[Image Registry][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-image-registry":    "[bz-Image Registry][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-image-registry",
			"[Image Registry][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-image-registry": "[bz-Image Registry][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-image-registry",
			"[sig-imageregistry] disruption/image-registry connection/new should be available throughout the test":               "[sig-imageregistry] Image registry remains available using new connections",
			"[sig-imageregistry] disruption/image-registry connection/reused should be available throughout the test":            "[sig-imageregistry] Image registry remains available using reused connections",

			"[sig-scheduling][Early] The openshift-image-registry pods [apigroup:imageregistry.operator.openshift.io] should be scheduled on different nodes [Skipped:SingleReplicaTopology] [Suite:openshift/conformance/parallel]": "[sig-scheduling][Early] The openshift-image-registry pods [apigroup:imageregistry.operator.openshift.io] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Image Registry"},
			},
			{
				IncludeAll: []string{"image-registry-", "-connections"},
			},
			{
				SIG: "sig-imageregistry",
			},
			{
				IncludeAny: []string{
					":ImageRegistry ",
					"samples/openshift-controller-manager/image-registry operators",
					"upgrade should succeed: image-registry",
				},
				Priority: 1,
			},
			{Suite: "ImageStream Manifest"},
			{Suite: "Image_Registry"},
			{Suite: "Testing registry"},
			{Suite: "Testing image registry operator"},
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
