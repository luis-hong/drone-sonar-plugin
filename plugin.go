package main

import (
	"drone-sonar-plugin/api"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type (
	Config struct {
		Key   string
		Name  string
		Host  string
		Token string

		Version         string
		Branch          string
		Sources         string
		Timeout         string
		Inclusions      string
		Exclusions      string
		Level           string
		ShowProfiling   string
		BranchAnalysis  bool
		UsingProperties bool

		PullrequestKey    string
		PullrequestBranch string
		PullrequestBase   string

		DroneRepoName string
	}
	Plugin struct {
		Config Config
	}
)

func (p Plugin) Exec() error {
	/*
	  - SAST_SONARQUBE_PROJECT_KEY=${DRONE_REPO_NAMESPACE}_${DRONE_REPO_NAME}
	   - "sonar-scanner \
	    -Dsonar.projectKey=$${SAST_SONARQUBE_PROJECT_KEY} \
	    -Dsonar.sources=. \
	    -Dsonar.host.url=https://$${SAST_SONARQUBE_WEB_DOMAIN} \
	    -Dsonar.login=$${SAST_SONARQUBE_SCANNER_TOKEN} \
	    -Dsonar.github.disableInlineComments=false \
	    -Dsonar.java.binaries=. \
	    -Dsonar.java.libraries=. \
	    -Dsonar.qualitygate.wait=true \
	    -Dsonar.qualitygate.timeout=1200 \
	    -Dsonar.issuesReport.console.enable=true \
	    -Dsonar.pullrequest.provider=github \
	    -Dsonar.pullrequest.key=${DRONE_PULL_REQUEST} \
	    -Dsonar.pullrequest.github.repository=${DRONE_REPO} \
	    -Dsonar.pullrequest.base=${DRONE_TARGET_BRANCH} \
	    -Dsonar.pullrequest.branch=${DRONE_SOURCE_BRANCH}  \
	    -Dsonar.verbose=true \
	    -Dsonar.ws.timeout=120"
	*/
	args := []string{
		"-Dsonar.host.url=" + p.Config.Host,
		"-Dsonar.login=" + p.Config.Token,
	}

	SonarProjectName := strings.Join([]string{"BZV", p.Config.DroneRepoName}, "_")
	if !p.Config.UsingProperties {
		argsParameter := []string{
			"-Dsonar.projectKey=" + SonarProjectName,
			"-Dsonar.projectName=" + SonarProjectName,
			//"-Dsonar.projectVersion=" + p.Config.Version,
			"-Dsonar.sources=" + p.Config.Sources,
			"-Dsonar.ws.timeout=" + p.Config.Timeout,
			"-Dsonar.inclusions=" + p.Config.Inclusions,
			"-Dsonar.exclusions=" + p.Config.Exclusions,
			"-Dsonar.log.level=" + p.Config.Level,
			"-Dsonar.showProfiling=" + p.Config.ShowProfiling,
			//"-Dsonar.scm.provider=git",
			"-Dsonar.scm.disabled=true",
			"-Dsonar.python.version=3.6, 3.7, 3.8, 3.9",
			"-Dsonar.qualitygate.wait=True",
			"-Dsonar.qualitygate.timeout=300",
			"-Dsonar.issuesReport.console.enable=true",
			"-Dsonar.github.disableInlineComments=false",
			"-Dsonar.java.binaries=.",
			"-Dsonar.java.libraries=.",
		}
		args = append(args, argsParameter...)
	}

	// sonar pr
	// https://sonarqube.inria.fr/sonarqube/documentation/analysis/github-integration/
	if len(p.Config.PullrequestKey) > 0 {
		/*
			-Dsonar.pullrequest.provider=github \
			-Dsonar.pullrequest.key=${DRONE_PULL_REQUEST} \
			-Dsonar.pullrequest.github.repository=${DRONE_REPO} \
			-Dsonar.pullrequest.base=${DRONE_TARGET_BRANCH} \
			-Dsonar.pullrequest.branch=${DRONE_SOURCE_BRANCH}  \
		*/
		argsParameter := []string{
			"-Dsonar.pullrequest.provider=github",
			"-Dsonar.pullrequest.github.repository=" + p.Config.Name,
			"-Dsonar.pullrequest.key=" + p.Config.PullrequestKey,
			"-Dsonar.pullrequest.branch=" + p.Config.PullrequestBranch,
			"-Dsonar.pullrequest.base=" + p.Config.PullrequestBase,
		}
		args = append(args, argsParameter...)
	} else {
		//if p.Config.BranchAnalysis {
		args = append(args, "-Dsonar.branch.name=main")
		//}
	}

	sonar := api.SonarLogin{
		BaseUrl: p.Config.Host,
		Token:   p.Config.Token,
	}

	var resp *http.Response

	// create project
	resp, _ = sonar.Post("api/projects/create", url.Values{
		"name":       {SonarProjectName},
		"project":    {SonarProjectName},
		"visibility": {"private"},
	})
	if _, err := io.ReadAll(resp.Body); err != nil {
		log.Fatal("[Project Init] create project: "+SonarProjectName, err)
	} else {
		log.Println("[Project Init] create project: " + SonarProjectName + "Create...")
	}

	// main branch setting
	// api/project_branches/rename?name=${DRONE_TARGET_BRANCH}&project=$${SAST_SONARQUBE_PROJECT_KEY}"
	resp, _ = sonar.Post("api/project_branches/create", url.Values{
		"name":    {p.Config.PullrequestBase},
		"project": {SonarProjectName},
	})
	if _, err := io.ReadAll(resp.Body); err != nil {
		log.Fatal("[Project Init] main branch setting: "+p.Config.PullrequestBase, err)
	} else {
		log.Println("[Project Init] main branch setting: " + p.Config.PullrequestBase + "Create...")
	}

	// alm setting
	// api/alm_settings/set_github_binding?
	// almSetting=$${SAST_SONARQUBE_GITHUB_APP}
	// &monorepo=false
	// &project=$${SAST_SONARQUBE_PROJECT_KEY}
	// &repository=${DRONE_REPO}
	// &summaryCommentEnabled=true
	resp, _ = sonar.Post("api/alm_settings/set_github_binding", url.Values{
		"almSetting":            {"pr-app"},
		"monorepo":              {"false"},
		"project":               {SonarProjectName},
		"repository":            {p.Config.DroneRepoName},
		"summaryCommentEnabled": {"true"},
	})
	if _, err := io.ReadAll(resp.Body); err != nil {
		log.Fatal("[Project Init] alm_settings: "+p.Config.PullrequestBase, err)
	} else {
		log.Println("[Project Init] alm_settings: " + p.Config.PullrequestBase + "Create...")
	}

	// api/qualitygates/select?
	// projectKey=$${SAST_SONARQUBE_PROJECT_KEY}
	// &gateName=${DRONE_REPO_NAME}
	resp, _ = sonar.Post("api/qualitygates/select", url.Values{
		"projectKey": {SonarProjectName},
		"gateName":   {p.Config.DroneRepoName},
	})
	if _, err := io.ReadAll(resp.Body); err != nil {
		log.Fatal("[Project Init] select quality gate: "+p.Config.DroneRepoName, err)
	} else {
		log.Println("[Project Init] select quality gate: " + p.Config.DroneRepoName + "Create...")
	}

	log.Println("=== ARGS ===")
	for _, e := range args {
		log.Printf("\t - %v \n", e)
	}
	log.Println("===========")

	cmd := exec.Command("sonar-scanner", args...)
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("==> Code Analysis Result:\n")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
