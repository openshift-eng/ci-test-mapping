package networkingrouter

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var RouterComponent = Component{
	Component: &config.Component{
		Name:                 "Networking / router",
		Operators:            []string{"ingress"},
		DefaultJiraComponent: "Networking / router",
		Namespaces: []string{
			"openshift-ingress",
			"openshift-ingress-canary",
			"openshift-ingress-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Routing"},
			},
			{
				SIG:        "sig-network",
				IncludeAll: []string{"Feature:Router"},
			},
			{
				SIG:        "sig-network-edge",
				IncludeAll: []string{"Feature:Router"},
			},
			{
				IncludeAll: []string{"ingress-to-", "disruption"},
			},
			{
				IncludeAll: []string{"openshift-ingress"},
			},
			{
				IncludeAll: []string{"via cluster ingress"},
			},
			{
				IncludeAll: []string{"Cluster frontend ingress"},
			},
			{
				IncludeAny: []string{
					"[sig-arch] Managed cluster should [apigroup:apps.openshift.io] should expose cluster services outside the cluster [apigroup:route.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
					"[sig-arch] Managed cluster should expose cluster services outside the cluster [apigroup:route.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Idling with a single service and ReplicationController should idle the service and ReplicationController properly [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many UDP senders (by continuing to drop all packets on the floor) [Serial] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should work with TCP (when fully idled) [Skipped:Network/OVNKubernetes] [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should work with TCP (when fully idled) [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should work with UDP [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling should handle many UDP senders (by continuing to drop all packets on the floor) [Serial] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling should work with TCP (when fully idled) [Skipped:Network/OVNKubernetes] [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling should work with UDP [Suite:openshift/conformance/parallel]",
					"[sig-networking] should not see excessive FailedToUpdateEndpointSlices Error updating Endpoint Slices",
				},
			},
			{Suite: "Network_Edge"},
			{Suite: "Routing and DNS related scenarios"},
			{Suite: "Test Ingress API logging options"},
			{Suite: "Testing HTTP Headers related scenarios"},
			{Suite: "Testing Ingress Operator related scenarios"},
			{Suite: "Testing haproxy rate limit related features"},
			{Suite: "Testing haproxy router"},
			{Suite: "Testing ingress to route object"},
			{Suite: "Testing route"},
			{Suite: "Testing timeout route"},
			{Suite: "Testing websocket features"},
			{Suite: "Testing wildcard routes"},
			{Suite: "route related"},
			{Suite: "Testing abrouting"},
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
