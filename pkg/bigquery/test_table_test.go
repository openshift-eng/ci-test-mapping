package bigquery

import (
	"testing"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
)

func TestBuildSuitesFilter(t *testing.T) {
	tests := []struct {
		name   string
		config *v1.Config
		want   string
	}{
		{
			name:   "no filter",
			config: &v1.Config{},
			want:   "1=1",
		},
		{
			name: "include exact only",
			config: &v1.Config{
				IncludeSuites: []string{"openshift-tests", "CNV-lp-interop"},
			},
			want: "(testsuite IN ('openshift-tests','CNV-lp-interop'))",
		},
		{
			name: "include patterns only",
			config: &v1.Config{
				IncludeSuitePatterns: []string{"lp-interop%", "e2e-%"},
			},
			want: "(testsuite LIKE 'lp-interop%' OR testsuite LIKE 'e2e-%')",
		},
		{
			name: "include exact and patterns",
			config: &v1.Config{
				IncludeSuites:         []string{"openshift-tests"},
				IncludeSuitePatterns:  []string{"lp-interop%"},
			},
			want: "(testsuite IN ('openshift-tests') OR testsuite LIKE 'lp-interop%')",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TestTableManager{config: tt.config}
			got := tm.buildSuitesFilter()
			if got != tt.want {
				t.Errorf("buildSuitesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildExcludeSuitesFilter(t *testing.T) {
	tests := []struct {
		name   string
		config *v1.Config
		want   string
	}{
		{
			name:   "no filter",
			config: &v1.Config{},
			want:   "",
		},
		{
			name: "exclude exact only",
			config: &v1.Config{
				ExcludeSuites: []string{"noise-suite"},
			},
			want: "AND testsuite NOT IN ('noise-suite')",
		},
		{
			name: "exclude patterns only",
			config: &v1.Config{
				ExcludeSuitePatterns: []string{"tmp-%", "skip-%"},
			},
			want: "AND testsuite NOT LIKE 'tmp-%' AND testsuite NOT LIKE 'skip-%'",
		},
		{
			name: "exclude exact and patterns",
			config: &v1.Config{
				ExcludeSuites:        []string{"exact-skip"},
				ExcludeSuitePatterns: []string{"skip-%"},
			},
			want: "AND testsuite NOT IN ('exact-skip') AND testsuite NOT LIKE 'skip-%'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TestTableManager{config: tt.config}
			got := tm.buildExcludeSuitesFilter()
			if got != tt.want {
				t.Errorf("buildExcludeSuitesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
