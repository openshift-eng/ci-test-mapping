package util

import (
	"reflect"
	"strings"
	"testing"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
)

func TestExtractField(t *testing.T) {
	tests := []struct {
		name       string
		test       string
		field      string
		wantValues []string
	}{
		{
			name:       "can extract bracketed single value",
			test:       "[sig-storage] In-tree Volumes [Driver: windows-gcepd] [Testpattern: Dynamic PV (ntfs)][Feature:Windows] subPath should be able to unmount after the subpath directory is deleted [LinuxOnly] [Skipped:NoOptionalCapabilities] [Suite:openshift/conformance/parallel] [Suite:k8s]",
			field:      "Driver",
			wantValues: []string{"windows-gcepd"},
		},
		{
			name:       "can extract slash single value",
			test:       `jira/"Test Framework" validate the thing works`,
			field:      "jira",
			wantValues: []string{"Test Framework"},
		},
		{
			name:       "can extract fields case-insensitive",
			test:       `jira/"Test Framework" [Jira: Networking]`,
			field:      "JIRA",
			wantValues: []string{"Test Framework", "Networking"},
		},
		{
			name:       "can extract slash multiple value",
			test:       `jira/"Test Framework" jira/Installer validate the thing works`,
			field:      "jira",
			wantValues: []string{"Test Framework", "Installer"},
		},
		{
			name:       "handles field not present",
			test:       "[sig-storage] Foobar",
			field:      "Driver",
			wantValues: nil,
		},
		{
			name:       "can extract multiple values",
			test:       "[sig-storage] [Driver: aws] [Driver: gcp]",
			field:      "Driver",
			wantValues: []string{"aws", "gcp"},
		},
		{
			name:       "values with whitespace",
			test:       "[sig-storage] In-tree Volumes [Driver: azure-disk] [Testpattern: Dynamic PV (default fs)] subPath should support readOnly file specified in the volumeMount [LinuxOnly] [Suite:openshift/conformance/parallel] [Suite:k8s]",
			field:      "Testpattern",
			wantValues: []string{"Dynamic PV (default fs)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResults := ExtractTestField(tt.test, tt.field); !reflect.DeepEqual(gotResults, tt.wantValues) {
				t.Errorf("ExtractTestField() = %v, want %v", gotResults, tt.wantValues)
			}
		})
	}
}

func TestRemoveField(t *testing.T) {
	tests := []struct {
		name      string
		test      string
		field     string
		wantValue string
	}{
		{
			name:      "can remove first in line test field",
			test:      "[sig-storage] In-tree Volumes [Driver: windows-gcepd] [Testpattern: Dynamic PV (ntfs)][Feature:Windows] subPath should be able to unmount after the subpath directory is deleted [LinuxOnly] [Skipped:NoOptionalCapabilities] [Suite:openshift/conformance/parallel] [Suite:k8s]",
			field:     "[Driver:",
			wantValue: "[sig-storage] In-tree Volumes [Testpattern: Dynamic PV (ntfs)][Feature:Windows] subPath should be able to unmount after the subpath directory is deleted [LinuxOnly] [Skipped:NoOptionalCapabilities] [Suite:openshift/conformance/parallel] [Suite:k8s]",
		},
		{
			name:      "can remove second in line test field",
			test:      "[sig-storage] In-tree Volumes [Driver: windows-gcepd] [Testpattern: Dynamic PV (ntfs)][Feature:Windows] subPath should be able to unmount after the subpath directory is deleted [LinuxOnly] [Skipped:NoOptionalCapabilities] [Suite:openshift/conformance/parallel] [Suite:k8s]",
			field:     "[Testpattern:",
			wantValue: "[sig-storage] In-tree Volumes [Driver: windows-gcepd][Feature:Windows] subPath should be able to unmount after the subpath directory is deleted [LinuxOnly] [Skipped:NoOptionalCapabilities] [Suite:openshift/conformance/parallel] [Suite:k8s]",
		},
		{
			name:      "can remove leading test field",
			test:      "[Monitor:generation-analyzer][Jira:\"kube-apiserver\"] monitor test generation-analyzer preparation",
			field:     "[Monitor:",
			wantValue: "[Jira:\"kube-apiserver\"] monitor test generation-analyzer preparation",
		},
		{
			name:      "can remove leading with space test field",
			test:      "[Monitor:generation-analyzer] [Jira:\"kube-apiserver\"] monitor test generation-analyzer preparation",
			field:     "[Monitor:",
			wantValue: "[Jira:\"kube-apiserver\"] monitor test generation-analyzer preparation",
		},
		{
			name:      "can remove trailing test field",
			test:      "[Jira:\"apiserver-auth\"] [Monitor:legacy-authentication-invariants] monitor test legacy-authentication-invariants preparation",
			field:     "[Monitor:",
			wantValue: "[Jira:\"apiserver-auth\"] monitor test legacy-authentication-invariants preparation",
		},
		{
			name:      "can remove trailing no space test field",
			test:      "[Jira:\"apiserver-auth\"][Monitor:legacy-authentication-invariants] monitor test legacy-authentication-invariants preparation",
			field:     "[Monitor:",
			wantValue: "[Jira:\"apiserver-auth\"] monitor test legacy-authentication-invariants preparation",
		},
		{
			name:      "doesn't remove field value",
			test:      "[Jira:\"Monitoring\"] monitor test metrics-api-availability interval construction",
			field:     "[Monitor:",
			wantValue: "[Jira:\"Monitoring\"] monitor test metrics-api-availability interval construction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResults := removeTestField(tt.test, tt.field); !reflect.DeepEqual(gotResults, tt.wantValue) {
				t.Errorf("removeTestField() = %v, \nwant %v", gotResults, tt.wantValue)
			}
		})
	}
}

func TestStableIdMatch(t *testing.T) {
	tests := []struct {
		name          string
		testInfo      v1.TestInfo
		testName      string
		matchTestName string
	}{
		{
			name:          "can match stable ID",
			testInfo:      v1.TestInfo{Suite: "kubernetes/test"},
			testName:      "[Jira:\"apiserver-auth\"][Monitor:legacy-authentication-invariants] monitor test legacy-authentication-invariants preparation",
			matchTestName: "[Jira:\"apiserver-auth\"] monitor test legacy-authentication-invariants preparation",
		},
		{
			name:          "can match stable ID",
			testInfo:      v1.TestInfo{Suite: "kubernetes/test"},
			testName:      "[Monitor:generation-analyzer][Jira:\"kube-apiserver\"] monitor test generation-analyzer preparation",
			matchTestName: "[Jira:\"kube-apiserver\"] monitor test generation-analyzer preparation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stableID := StableID(&tt.testInfo, tt.testName)
			matchID := StableID(&tt.testInfo, tt.matchTestName)
			if !strings.EqualFold(stableID, matchID) {
				t.Errorf("removeTestField() = %v, \nwant %v", stableID, matchID)
			}
		})
	}
}
