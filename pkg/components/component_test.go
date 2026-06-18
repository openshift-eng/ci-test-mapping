package components

import (
	"reflect"
	"strings"
	"testing"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/storage"
	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
)

func TestJiraTagWinner(t *testing.T) {
	tests := []struct {
		name          string
		test          *v1.TestInfo
		ownerships    []*v1.TestOwnership
		wantNil       bool
		wantError     string
		wantComponent string
	}{
		{
			name: "jira tag resolves conflict to tagged component",
			test: &v1.TestInfo{
				Name: "[Jira:HyperShift] some test with cluster-network-operator",
			},
			ownerships: []*v1.TestOwnership{
				{Component: "HyperShift", JIRAComponent: "HyperShift"},
				{Component: "Networking / cluster-network-operator", JIRAComponent: "Networking / cluster-network-operator"},
			},
			wantComponent: "HyperShift",
		},
		{
			name: "no jira tag returns nil",
			test: &v1.TestInfo{
				Name: "some test without jira tag",
			},
			ownerships: []*v1.TestOwnership{
				{Component: "A", JIRAComponent: "A"},
				{Component: "B", JIRAComponent: "B"},
			},
			wantNil: true,
		},
		{
			name: "jira tag matches no claimant returns nil",
			test: &v1.TestInfo{
				Name: "[Jira:Unrelated] some test",
			},
			ownerships: []*v1.TestOwnership{
				{Component: "A", JIRAComponent: "A"},
				{Component: "B", JIRAComponent: "B"},
			},
			wantNil: true,
		},
		{
			name: "single ownership returns nil (no conflict to resolve)",
			test: &v1.TestInfo{
				Name: "[Jira:A] some test",
			},
			ownerships: []*v1.TestOwnership{
				{Component: "A", JIRAComponent: "A"},
			},
			wantNil: true,
		},
		{
			name: "multiple jira tags matching multiple claimants returns error",
			test: &v1.TestInfo{
				Name:  "[Jira:A][Jira:B] some test",
				Suite: "test-suite",
			},
			ownerships: []*v1.TestOwnership{
				{Component: "A", JIRAComponent: "A"},
				{Component: "B", JIRAComponent: "B"},
			},
			wantError: "at most one Jira component claimant",
		},
		{
			name: "jira tag match is case-insensitive",
			test: &v1.TestInfo{
				Name: "[Jira:hypershift] some test",
			},
			ownerships: []*v1.TestOwnership{
				{Component: "HyperShift", JIRAComponent: "HyperShift"},
				{Component: "Other", JIRAComponent: "Other"},
			},
			wantComponent: "HyperShift",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jiraTagWinner(tt.test, tt.ownerships)
			if tt.wantError != "" {
				if err == nil {
					t.Fatalf("jiraTagWinner() error = nil, want error containing %q", tt.wantError)
				}
				if !strings.Contains(err.Error(), tt.wantError) {
					t.Errorf("jiraTagWinner() error = %q, want error containing %q", err.Error(), tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("jiraTagWinner() unexpected error: %v", err)
			}
			if tt.wantNil {
				if got != nil {
					t.Errorf("jiraTagWinner() = %+v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("jiraTagWinner() = nil, want a winner")
			}
			if got.Component != tt.wantComponent {
				t.Errorf("jiraTagWinner().Component = %q, want %q", got.Component, tt.wantComponent)
			}
		})
	}
}

func TestIdentifyTest(t *testing.T) {
	componentRegistry := registry.NewComponentRegistry()

	tests := []struct {
		before           func() error
		name             string
		testInfo         *v1.TestInfo
		wantError        string
		wantComponent    string
		wantCapabilities []string
		wantID           string
		after            func() error
	}{
		{
			name:             "identifies the correct component and capability",
			testInfo:         &v1.TestInfo{Name: "[sig-storage][Feature:foobar] component with feature"},
			wantComponent:    "Storage",
			wantCapabilities: []string{"foobar"},
		},
		{
			name:             "identifies the correct component with default capability",
			testInfo:         &v1.TestInfo{Name: "[sig-storage] component with unknown capability"},
			wantComponent:    "Storage",
			wantCapabilities: []string{"Other"},
		},
		{
			name:             "handles unknown capability",
			testInfo:         &v1.TestInfo{Name: "[sig-something] what even is this"},
			wantComponent:    "Unknown",
			wantCapabilities: []string{"Other"},
		},
		{
			name:      "detects duplicate owners without priority",
			testInfo:  &v1.TestInfo{Name: "[sig-storage] A storage test"},
			wantError: "unable to resolve conflict",
			before: func() error {
				componentRegistry.Register("storage2", &storage.StorageComponent)
				return nil
			},
		},
		{
			name:             "categorizes capability based on variants",
			testInfo:         &v1.TestInfo{Name: "[sig-qe] test with variants", Variants: []string{"Procedure:automated-release"}},
			wantCapabilities: []string{"LEVEL0"},
		},
		{
			name: "identifies the correct testID for kubernetes renamed test",
			testInfo: &v1.TestInfo{
				Name: "[sig-network] DNS should provide DNS for the cluster [Conformance] [Suite:openshift/conformance/parallel/minimal] [Suite:k8s]",
			},
			wantID: "[sig-network] DNS should provide DNS for the cluster [Conformance] [Skipped:Proxy] [Suite:openshift/conformance/parallel/minimal] [Suite:k8s]",
		},
		{
			name: "identifies the correct testID for origin renamed test",
			testInfo: &v1.TestInfo{
				Name: "[sig-arch][Early] Managed cluster should [apigroup:config.openshift.io] start all core operators [Suite:openshift/conformance/parallel]",
			},
			wantID: "[sig-arch][Early] Managed cluster should [apigroup:config.openshift.io] start all core operators [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
		},
	}
	ti := NewTestIdentifier(componentRegistry, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				if err := tt.before(); err != nil {
					t.Fatalf("before() failed: %+v", err)
				}
			}

			testOwnership, err := ti.Identify(tt.testInfo)

			if err != nil {
				if tt.wantError == "" {
					t.Fatalf("IdentifyTest() returned unexpected err: %+v", err)
				}
				if !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("IdentifyTest() did not return expected err %q: %+v", tt.wantError, err)
				}
			} else if tt.wantError != "" {
				t.Fatalf("IdentifyTest() did not return expected err: %+v", err)
			}

			if tt.wantComponent != "" && testOwnership.Component != tt.wantComponent {
				t.Errorf("IdentifyTest() gotComponent = %v, want %v", testOwnership.Component, tt.wantComponent)
			}
			if tt.wantCapabilities != nil && !reflect.DeepEqual(testOwnership.Capabilities, tt.wantCapabilities) {
				t.Errorf("IdentifyTest() gotCapabilities = %v, want %v", testOwnership.Capabilities, tt.wantCapabilities)
			}
			if tt.wantID != "" && testOwnership.ID != tt.wantID {
				t.Errorf("IdentifyTest() gotID = %s, want %s", testOwnership.ID, tt.wantID)
			}
		})
	}
}
