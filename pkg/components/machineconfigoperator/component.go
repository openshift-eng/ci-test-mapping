package machineconfigoperator

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var MachineConfigOperatorComponent = Component{
	Component: &config.Component{
		Name:                 "Machine Config Operator",
		Operators:            []string{"machine-config"},
		DefaultJiraComponent: "Machine Config Operator",
		Namespaces: []string{
			"openshift-machine-config-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Machine Config Operator"},
			},
			{
				IncludeAll: []string{"bz-machine config operator"},
			},
			{
				IncludeAll: []string{"machine-config-operator"},
			},
			{
				IncludeAny: []string{
					"machine-config-operator",
					"node count should match or exceed machine count",
					"OSUpdateStarted event should be recorded for nodes that reach OSUpdateStaged",
					"upgrade should succeed: machine-config",
				},
			},
			{
				SIG:          "sig-cluster-lifecycle",
				IncludeAll:   []string{"Pods cannot access the /config"},
				Capabilities: []string{"Config"},
			},
			{Suite: "MCO"},
			{
				SIG: "sig-mco",
			},
		},
		TestRenames: map[string]string{
			"[Machine Config Operator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-machine-config-operator":    "[bz-Machine Config Operator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-machine-config-operator",
			"[Machine Config Operator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-machine-config-operator": "[bz-Machine Config Operator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-machine-config-operator",

			// Remap MachineConfigNodes tests
			"[sig-mco][OCPFeatureGate:MachineConfigNodes] [Suite:openshift/machine-config-operator/disruptive][Disruptive]Should properly report MCN conditions on node degrade [apigroup:machineconfiguration.openshift.io] [Serial]":                     "[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:MachineConfigNodes] [Serial][Slow]Should properly report MCN conditions on node degrade [apigroup:machineconfiguration.openshift.io]",
			"[sig-mco][OCPFeatureGate:MachineConfigNodes] Should have MCN properties matching associated node properties for nodes in default MCPs [apigroup:machineconfiguration.openshift.io] [Suite:openshift/conformance/parallel]":                    "[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:MachineConfigNodes] [Suite:openshift/conformance/parallel]Should have MCN properties matching associated node properties for nodes in default MCPs [apigroup:machineconfiguration.openshift.io]",
			"[sig-mco][OCPFeatureGate:MachineConfigNodes] Should properly block MCN updates by impersonation of the MCD SA [apigroup:machineconfiguration.openshift.io] [Suite:openshift/conformance/parallel]":                                            "[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:MachineConfigNodes] [Suite:openshift/conformance/parallel]Should properly block MCN updates by impersonation of the MCD SA [apigroup:machineconfiguration.openshift.io]",
			"[sig-mco][OCPFeatureGate:MachineConfigNodes] Should properly block MCN updates from a MCD that is not the associated one [apigroup:machineconfiguration.openshift.io] [Suite:openshift/conformance/parallel]":                                 "[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:MachineConfigNodes] [Suite:openshift/conformance/parallel]Should properly block MCN updates from a MCD that is not the associated one [apigroup:machineconfiguration.openshift.io]",
			"[sig-mco][OCPFeatureGate:MachineConfigNodes] [Serial]Should have MCN properties matching associated node properties for nodes in custom MCPs [apigroup:machineconfiguration.openshift.io] [Suite:openshift/conformance/serial]":               "[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:MachineConfigNodes] [Suite:openshift/conformance/serial][Serial]Should have MCN properties matching associated node properties for nodes in custom MCPs [apigroup:machineconfiguration.openshift.io]",
			"[sig-mco][OCPFeatureGate:MachineConfigNodes] [Serial]Should properly transition through MCN conditions on rebootless node update [apigroup:machineconfiguration.openshift.io] [Suite:openshift/conformance/serial]":                           "[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:MachineConfigNodes] [Suite:openshift/conformance/serial][Serial]Should properly transition through MCN conditions on rebootless node update [apigroup:machineconfiguration.openshift.io]",
			"[sig-mco][OCPFeatureGate:MachineConfigNodes] [Serial]Should properly update the MCN from the associated MCD [apigroup:machineconfiguration.openshift.io] [Suite:openshift/conformance/serial]":                                                "[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:MachineConfigNodes] [Suite:openshift/conformance/serial][Serial]Should properly update the MCN from the associated MCD [apigroup:machineconfiguration.openshift.io]",
			"[sig-mco][OCPFeatureGate:MachineConfigNodes] [Suite:openshift/machine-config-operator/disruptive][Disruptive][Slow]Should properly create and remove MCN on node creation and deletion [apigroup:machineconfiguration.openshift.io] [Serial]": "[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:MachineConfigNodes] [Serial][Slow]Should properly create and remove MCN on node creation and deletion [apigroup:machineconfiguration.openshift.io]",

			// Remap PinnedImages tests
			"[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:PinnedImages][Disruptive] All Nodes in a Custom Pool should have the PinnedImages in PIS [apigroup:machineconfiguration.openshift.io] [Serial]":                              "[Suite:openshift/machine-config-operator/disruptive][Suite:openshift/conformance/serial][sig-mco][OCPFeatureGate:PinnedImages][OCPFeatureGate:MachineConfigNodes][Serial] All Nodes in a Custom Pool should have the PinnedImages in PIS [apigroup:machineconfiguration.openshift.io]",
			"[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:PinnedImages][Disruptive] All Nodes in a standard Pool should have the PinnedImages PIS [apigroup:machineconfiguration.openshift.io] [Serial]":                               "[Suite:openshift/machine-config-operator/disruptive][Suite:openshift/conformance/serial][sig-mco][OCPFeatureGate:PinnedImages][OCPFeatureGate:MachineConfigNodes][Serial] All Nodes in a standard Pool should have the PinnedImages PIS [apigroup:machineconfiguration.openshift.io]",
			"[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:PinnedImages][Disruptive] Invalid PIS leads to degraded MCN in a custom Pool [apigroup:machineconfiguration.openshift.io] [Serial]":                                          "[Suite:openshift/machine-config-operator/disruptive][Suite:openshift/conformance/serial][sig-mco][OCPFeatureGate:PinnedImages][OCPFeatureGate:MachineConfigNodes][Serial] Invalid PIS leads to degraded MCN in a custom Pool [apigroup:machineconfiguration.openshift.io]",
			"[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:PinnedImages][Disruptive] Invalid PIS leads to degraded MCN in a standard Pool [apigroup:machineconfiguration.openshift.io] [Serial]":                                        "[Suite:openshift/machine-config-operator/disruptive][Suite:openshift/conformance/serial][sig-mco][OCPFeatureGate:PinnedImages][OCPFeatureGate:MachineConfigNodes][Serial] Invalid PIS leads to degraded MCN in a standard Pool [apigroup:machineconfiguration.openshift.io]",
			"[Suite:openshift/machine-config-operator/disruptive][sig-mco][OCPFeatureGate:PinnedImages][Disruptive] [Slow]All Nodes in a custom Pool should have the PinnedImages even after Garbage Collection [apigroup:machineconfiguration.openshift.io] [Serial]": "[Suite:openshift/machine-config-operator/disruptive][Suite:openshift/conformance/serial][sig-mco][OCPFeatureGate:PinnedImages][OCPFeatureGate:MachineConfigNodes][Serial] All Nodes in a custom Pool should have the PinnedImages even after Garbage Collection [apigroup:machineconfiguration.openshift.io]",
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
