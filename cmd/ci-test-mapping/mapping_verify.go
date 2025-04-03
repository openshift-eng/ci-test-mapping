package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/components"
)

var verifyMapFlags = NewVerifyMapFlags()
var mappingFiles = []string{
	"data/openshift-gce-devel/ci_analysis_us/component_mapping.json",
	"data/openshift-gce-devel/ci_analysis_qe/component_mapping.json",
}

type VerifyMapFlags struct {
	MappingFiles []string
}

func NewVerifyMapFlags() *VerifyMapFlags {
	return &VerifyMapFlags{}
}

func (f *VerifyMapFlags) BindFlags(fs *pflag.FlagSet) {
	fs.StringArrayVar(&f.MappingFiles, "mapping-files", mappingFiles, "Which mapping files to verify")
}

var verifyMapCmd = &cobra.Command{
	Use:   "map-verify",
	Short: "Verify mapping is correct, i.e. no tests are moving from assigned to unknown",
	Run: func(cmd *cobra.Command, args []string) {

		now := time.Now()
		logrus.Infof("verifying mappings are correct...")

		for _, mf := range verifyMapFlags.MappingFiles {
			url := fmt.Sprintf("https://raw.githubusercontent.com/openshift-eng/ci-test-mapping/main/%s", mf)

			// **** Get the current mappings
			response, err := http.Get(url)
			if err != nil {
				fmt.Println("Failed to fetch data:", err)
				return
			}
			if response.StatusCode != 200 {
				fmt.Println("Failed to fetch data:", response.StatusCode)
				return
			}

			currentMap, err := fetchMap(response.Body)
			if err != nil {
				logrus.WithError(err).Fatalf("couldn't fetch current mapping")
			}
			_ = response.Body.Close()

			// Create current mapping lookup table
			type cl struct {
				name  string
				suite string
			}
			currentComponents := make(map[cl]string)
			for _, mapping := range currentMap {
				currentComponents[cl{mapping.Name, mapping.Suite}] = mapping.Component
			}

			// **** Get the mappings from the file
			file, err := os.Open(mf)
			if err != nil {
				fmt.Println("Failed to open file:", err)
				return
			}
			newMap, err := fetchMap(file)
			if err != nil {
				logrus.WithError(err).Fatalf("couldn't read new mapping from file")
			}
			_ = file.Close()

			// Look for removed mappings
			removedMaps := make([]cl, 0)
			for _, n := range newMap {
				if n.Component != components.DefaultComponent {
					continue
				}
				locator := cl{n.Name, n.Suite}
				if old, ok := currentComponents[locator]; ok && old != components.DefaultComponent {
					removedMaps = append(removedMaps, locator)
				}
			}

			logrus.Infof("verification complete in %+v", time.Since(now))
			if len(removedMaps) == 0 {
				return
			}

			for _, removed := range removedMaps {
				logrus.WithFields(logrus.Fields{
					"Name":  removed.name,
					"Suite": removed.suite,
				}).Warningf("test moved from %q to \"Unknown\"", currentComponents[removed])
			}
			logrus.Fatalf("Components are not allowed to move to Unknown. Please assign correct ownership.")
		}
	},
}

func fetchMap(f io.Reader) ([]v1.TestOwnership, error) {
	body, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var mapping []v1.TestOwnership
	err = json.Unmarshal(body, &mapping)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return mapping, nil
}

func init() {
	verifyMapFlags.BindFlags(verifyMapCmd.Flags())
	rootCmd.AddCommand(verifyMapCmd)
}
