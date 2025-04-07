package openshiftapiserver

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var OpenshiftApiserverComponent = Component{
	Component: &config.Component{
		Name:                 "openshift-apiserver",
		Operators:            []string{"openshift-apiserver", "openshift-lifecycle-manager", "openshift-lifecycle-manager-catalog", "openshift-lifecycle-manager-packageserver", "platform-operators-aggregated"},
		DefaultJiraComponent: "openshift-apiserver",
		Namespaces: []string{
			"openshift-apiserver",
			"openshift-apiserver-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-openshift-apiserver"},
			},
			{
				IncludeAll: []string{"single second disruption", "openshift-api"},
			},
			{
				IncludeAll: []string{"openshift-api-", "-connections"},
			},
			{Suite: "Projects"},
			{Suite: "project list tests"},
			{Suite: "project permissions"},
			{
				IncludeAny: []string{
					"upgrade should succeed: openshift-apiserver",
				},
			},
		},
		TestRenames: map[string]string{
			"[openshift-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-apiserver":                                    "[bz-openshift-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-apiserver",
			"[openshift-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-apiserver-operator":                           "[bz-openshift-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-apiserver-operator",
			"[openshift-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-apiserver":                                 "[bz-openshift-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-apiserver",
			"[openshift-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-apiserver-operator":                        "[bz-openshift-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-apiserver-operator",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-openshift-api connection/new":                              "[sig-network] there should be nearly zero single second disruptions for cache-openshift-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-openshift-api connection/new target=api-int":               "[sig-network] there should be nearly zero single second disruptions for cache-openshift-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-openshift-api connection/new target=service-network":       "[sig-network] there should be nearly zero single second disruptions for cache-openshift-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/service-network/cache-openshift-api connection/new":              "[sig-network] there should be nearly zero single second disruptions for cache-openshift-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-openshift-api connection/reused":                           "[sig-network] there should be nearly zero single second disruptions for cache-openshift-api-reused-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-openshift-api connection/reused target=api-int":            "[sig-network] there should be nearly zero single second disruptions for cache-openshift-api-reused-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-openshift-api connection/reused target=service-network":    "[sig-network] there should be nearly zero single second disruptions for cache-openshift-api-reused-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/service-network/cache-openshift-api connection/reused":           "[sig-network] there should be nearly zero single second disruptions for cache-openshift-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-openshift-api connection/new":                           "[sig-network] there should be reasonably few single second disruptions for cache-openshift-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-openshift-api connection/new target=api-int":            "[sig-network] there should be reasonably few single second disruptions for cache-openshift-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-openshift-api connection/new target=service-network":    "[sig-network] there should be reasonably few single second disruptions for cache-openshift-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/service-network/cache-openshift-api connection/new":           "[sig-network] there should be reasonably few single second disruptions for cache-openshift-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-openshift-api connection/reused":                        "[sig-network] there should be reasonably few single second disruptions for cache-openshift-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-openshift-api connection/reused target=api-int":         "[sig-network] there should be reasonably few single second disruptions for cache-openshift-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-openshift-api connection/reused target=service-network": "[sig-network] there should be reasonably few single second disruptions for cache-openshift-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/service-network/cache-openshift-api connection/reused":        "[sig-network] there should be reasonably few single second disruptions for cache-openshift-api-reused-connections",

			"[sig-scheduling][Early] The openshift-apiserver pods [apigroup:apps.openshift.io][apigroup:authorization.openshift.io][apigroup:build.openshift.io][apigroup:image.openshift.io][apigroup:project.openshift.io][apigroup:quota.openshift.io][apigroup:route.openshift.io][apigroup:security.openshift.io][apigroup:template.openshift.io] should be scheduled on different nodes [Skipped:SingleReplicaTopology] [Suite:openshift/conformance/parallel]": "[sig-scheduling][Early] The openshift-apiserver pods [apigroup:apps.openshift.io][apigroup:authorization.openshift.io][apigroup:build.openshift.io][apigroup:image.openshift.io][apigroup:project.openshift.io][apigroup:quota.openshift.io][apigroup:route.openshift.io][apigroup:security.openshift.io][apigroup:template.openshift.io] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
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
