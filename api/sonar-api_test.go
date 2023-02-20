package api_test

import (
	"drone-sonar-plugin/api"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
   - SAST_SONARQUBE_PROJECT_KEY=${DRONE_REPO_NAMESPACE}_${DRONE_REPO_NAME}
   - curl -v -L -u "$${SAST_SONARQUBE_API_TOKEN}:" -X POST "https://$${SAST_SONARQUBE_WEB_DOMAIN}/api/projects/create?name=$${SAST_SONARQUBE_PROJECT_KEY}&project=$${SAST_SONARQUBE_PROJECT_KEY}&visibility=private"
   - curl -v -L -u "$${SAST_SONARQUBE_API_TOKEN}:" -X POST "https://$${SAST_SONARQUBE_WEB_DOMAIN}/api/project_branches/rename?name=${DRONE_TARGET_BRANCH}&project=$${SAST_SONARQUBE_PROJECT_KEY}"
   - curl -v -L -u "$${SAST_SONARQUBE_API_TOKEN}:" -X POST "https://$${SAST_SONARQUBE_WEB_DOMAIN}/api/alm_settings/set_github_binding?almSetting=$${SAST_SONARQUBE_GITHUB_APP}&monorepo=false&project=$${SAST_SONARQUBE_PROJECT_KEY}&repository=${DRONE_REPO}&summaryCommentEnabled=true"
   - curl -v -L -u "$${SAST_SONARQUBE_API_TOKEN}:" -X POST "https://$${SAST_SONARQUBE_WEB_DOMAIN}/api/qualitygates/select?projectKey=$${SAST_SONARQUBE_PROJECT_KEY}&gateName=${DRONE_REPO_NAME}"
*/

var SONAR api.SonarLogin

func TestMain(m *testing.M) {
	SONAR = api.SonarLogin{
		BaseUrl: "http://172.16.32.199:9000",
		Token:   "squ_103df73c173a29e1cabc7e9390f57d3046435ca3",
	}
	exitCode := m.Run()

	fmt.Println("TestMain End")
	os.Exit(exitCode)
}

func TestGetProjects(t *testing.T) {
	// user token required
	v := url.Values{}
	resp, err := SONAR.Get("api/projects/search", v)

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
	}
	t.Logf("%v", string(bodyText))
	assert.Equal(t, http.StatusOK, resp.StatusCode, "they should be equal")
}

func TestPostCreateProject(t *testing.T) {
	// user token required
	// /api/projects/create?
	// name=$${SAST_SONARQUBE_PROJECT_KEY}&
	// project=$${SAST_SONARQUBE_PROJECT_KEY}&
	// visibility=private"

	formData := url.Values{
		"name":       {"PostCreateTestProject"},
		"project":    {"PostCreateTestProject"},
		"visibility": {"private"},
	}

	resp, err := SONAR.Post("api/projects/create", formData)
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
	}
	t.Logf("%v", string(bodyText))

	//assert.NotEqual(t, http.StatusOK, resp.Status, "they should be equal")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "they should be equal")

	resp2, err := SONAR.Post("api/projects/delete", formData)
	if err != nil {
		t.Log(err)
	}
	assert.Equal(t, http.StatusNoContent, resp2.StatusCode, "they should be equal")
}
