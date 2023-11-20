package jira

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
)

func GetJiraComponents() (map[string]uint, error) {
	start := time.Now()
	log.Infof("loading jira ocpbugs component information...")
	body, err := jiraRequest("https://issues.redhat.com/rest/api/2/project/12332330/components")
	if err != nil {
		return nil, err
	}

	var components []v1.JiraComponent
	err = json.Unmarshal(body, &components)
	if err != nil {
		return nil, err
	}

	ids := make(map[string]uint)
	for _, c := range components {
		jiraID, err := strconv.ParseUint(c.ID, 10, 64)
		if err != nil {
			msg := "error parsing jira ID"
			log.WithError(err).Warn(msg)
		}

		ids[c.Name] = uint(jiraID)
	}

	log.Infof("jira ocpbugs components loaded in %+v", time.Since(start))
	return ids, nil
}

func jiraRequest(apiURL string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// For a fresh sync to a developer DB, no jira token is needed since the issues API is still open. However, if
	// we need to find cards where the trt-incident label was removed this API is not protected and returns 401
	// if tried unauthed.  So really this only affects long-lived instances of Sippy.
	//
	// WARNING: DO NOT give public-facing Sippy a personal developer token, use a service account that is not marked
	// as a Red Hat employee.
	token := os.Getenv("JIRA_TOKEN")
	if token == "" {
		log.Warningf("no token, proceeding unauthenticated (ok for most queries)")
	} else {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
