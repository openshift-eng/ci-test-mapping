package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
)

type VerifyFlags struct {
	JiraURL string
}

func NewVerifyFlags() *VerifyFlags {
	return &VerifyFlags{
		JiraURL: "https://redhat.atlassian.net/rest/api/2/issue/createmeta/OCPBUGS/issuetypes/",
	}
}

var verifyCmd = &cobra.Command{
	Use:   "jira-verify",
	Short: "Verify all JIRA components are represented in ci-test-mapping",
	Run: func(cmd *cobra.Command, args []string) {
		f := NewVerifyFlags()

		logrus.Info("fetching jira components from jira...")
		bearerToken := os.Getenv("JIRA_TOKEN")
		if len(bearerToken) == 0 {
			basicToken := os.Getenv("JIRA_TOKEN_BASIC")
			if len(basicToken) == 0 {
				cmd.Usage() // nolint:errcheck
				logrus.Fatal("jira token required")
			}
		}
		jiraComponents, err := getJiraComponents(f.JiraURL)
		if err != nil {
			logrus.WithError(err).Fatal("could not fetch jira components")
		}

		logrus.Info("verifying all components have a mapping")
		reg := registry.NewComponentRegistry()
		mappedJiraComponents := sets.New[string]()
		for _, c := range reg.Components {
			mappedJiraComponents.Insert(c.JiraComponents()...)
		}

		for _, component := range jiraComponents {
			if !mappedJiraComponents.Has(component) {
				logrus.WithFields(logrus.Fields{
					"path":    "pkg/" + getPackagePath(component),
					"package": getPackageName(component),
				}).Warningf("no mapping for jira component %q", component)
			}
		}

		jiraComponentsSet := sets.New[string](jiraComponents...)
		for _, mappedComponent := range mappedJiraComponents.UnsortedList() {
			if !jiraComponentsSet.Has(mappedComponent) && mappedComponent != "" {
				logrus.Errorf("unknown component %q not found in jira", mappedComponent)
			}
		}

		logrus.Info("done!")
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
