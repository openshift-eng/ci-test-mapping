package components

import (
	"reflect"
	"testing"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
)

// TestE2EAnalysisReadinessRows pins the three readiness rows emitted by
// openshift-tests e2e-analysis (suite openshift-tests) to their owning
// components. Each matcher pins both the exact test name and Suite
// openshift-tests so no other suite with a similar name is claimed.
// Origin: pkg/e2eanalysis in openshift/origin via the gather-extra step,
// which hardcodes suite openshift-tests.
func TestE2EAnalysisReadinessRows(t *testing.T) {
	t.Parallel()
	rows := []struct {
		name          string
		wantComponent string
		wantPriority  int
	}{
		{"verify the cluster readiness and stability", "Cluster Version Operator", 0},
		{"ensure 1 worker node at least gets ready", "Node / Kubelet", 1},
		{"verify all machines should be in Running state", "Cloud Compute / Machine API Providers", 1},
	}
	identifier := NewTestIdentifier(registry.NewComponentRegistry(), nil)
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			ownership, err := identifier.Identify(&v1.TestInfo{Name: row.name, Suite: "openshift-tests"})
			if err != nil || ownership == nil {
				t.Fatalf("Identify() = %v, %v", ownership, err)
			}
			if ownership.Component != row.wantComponent || ownership.JIRAComponent != row.wantComponent {
				t.Errorf("Identify() = %s/%s, want %s", ownership.Component, ownership.JIRAComponent, row.wantComponent)
			}
			if ownership.Priority != row.wantPriority {
				t.Errorf("Identify() priority = %d, want %d", ownership.Priority, row.wantPriority)
			}
			if !reflect.DeepEqual(ownership.Capabilities, []string{"Other"}) {
				t.Errorf("Identify() capabilities = %v, want [Other]", ownership.Capabilities)
			}
		})
	}
}

// TestE2EAnalysisReadinessSuitePin verifies that the name+suite pin does not
// leak to other suites. A same-named test in a different suite must remain
// unclaimed (Unknown) so the claim is confined to openshift-tests.
func TestE2EAnalysisReadinessSuitePin(t *testing.T) {
	t.Parallel()
	identifier := NewTestIdentifier(registry.NewComponentRegistry(), nil)
	tests := []string{
		"verify the cluster readiness and stability",
		"ensure 1 worker node at least gets ready",
		"verify all machines should be in Running state",
	}
	for _, name := range tests {
		t.Run(name+"/other-suite", func(t *testing.T) {
			t.Parallel()
			ownership, err := identifier.Identify(&v1.TestInfo{Name: name, Suite: "other-suite"})
			if err != nil || ownership == nil {
				t.Fatalf("Identify() = %v, %v", ownership, err)
			}
			if ownership.Component != "Unknown" {
				t.Errorf("Identify() = %s, want Unknown for suite other-suite", ownership.Component)
			}
		})
		t.Run(name+"/empty-suite", func(t *testing.T) {
			t.Parallel()
			ownership, err := identifier.Identify(&v1.TestInfo{Name: name, Suite: ""})
			if err != nil || ownership == nil {
				t.Fatalf("Identify() = %v, %v", ownership, err)
			}
			if ownership.Component != "Unknown" {
				t.Errorf("Identify() = %s, want Unknown for empty suite", ownership.Component)
			}
		})
	}
}
