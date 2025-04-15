package components

import (
	"reflect"
	"strings"
	"testing"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/storage"
	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
)

func TestIdentifyTest(t *testing.T) {
	componentRegistry := registry.NewComponentRegistry()

	tests := []struct {
		before           func() error
		name             string
		testInfo         *v1.TestInfo
		wantError        string
		wantComponent    string
		wantCapabilities []string
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
		})
	}
}
