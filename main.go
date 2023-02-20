package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

var build = "1" // build number set at compile time

func main() {
	app := cli.NewApp()
	app.Name = "Drone-Sonar-Plugin"
	app.Usage = "Drone plugin to integrate with SonarQube."
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:   "key",
			Usage:  "project key",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "name",
			Usage:  "project name",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "host",
			Usage:  "SonarQube host",
			EnvVar: "PLUGIN_SONAR_HOST",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "SonarQube token",
			EnvVar: "PLUGIN_SONAR_TOKEN",
		},

		// advanced parameters
		cli.StringFlag{
			Name:   "ver",
			Usage:  "Project version",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "branch",
			Usage:  "Project branch",
			EnvVar: "DRONE_BRANCH",
		},
		cli.StringFlag{
			Name:   "timeout",
			Usage:  "Web request timeout",
			Value:  "60",
			EnvVar: "PLUGIN_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "sources",
			Usage:  "analysis sources",
			Value:  ".",
			EnvVar: "PLUGIN_SOURCES",
		},
		cli.StringFlag{
			Name:   "inclusions",
			Usage:  "code inclusions",
			EnvVar: "PLUGIN_INCLUSIONS",
		},
		cli.StringFlag{
			Name:   "exclusions",
			Usage:  "code exclusions",
			EnvVar: "PLUGIN_EXCLUSIONS",
		},
		cli.StringFlag{
			Name:   "level",
			Usage:  "log level",
			Value:  "INFO",
			EnvVar: "PLUGIN_LEVEL",
		},
		cli.StringFlag{
			Name:   "showProfiling",
			Usage:  "showProfiling during analysis",
			Value:  "false",
			EnvVar: "PLUGIN_SHOWPROFILING",
		},
		cli.BoolFlag{
			Name:   "branchAnalysis",
			Usage:  "execute branchAnalysis",
			EnvVar: "PLUGIN_BRANCHANALYSIS",
		},
		cli.BoolFlag{
			Name:   "usingProperties",
			Usage:  "using sonar-project.properties",
			EnvVar: "PLUGIN_USINGPROPERTIES",
		},
		// sonar pr
		cli.StringFlag{
			Name:   "pullrequestKey",
			Usage:  "sonar.pullrequest.key",
			EnvVar: "DRONE_PULL_REQUEST",
		},
		cli.StringFlag{
			Name:   "pullrequestBranch",
			Usage:  "sonar.pullrequest.branch",
			EnvVar: "DRONE_SOURCE_BRANCH",
		},
		cli.StringFlag{
			Name:   "pullrequestBase",
			Usage:  "sonar.pullrequest.base",
			EnvVar: "DRONE_TARGET_BRANCH",
		},
		cli.StringFlag{
			Name:   "droneRepoName",
			Usage:  "drone repo name",
			EnvVar: "DRONE_REPO_NAME",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func run(c *cli.Context) {
	plugin := Plugin{
		Config: Config{
			Key:   c.String("key"),
			Name:  c.String("name"),
			Host:  c.String("host"),
			Token: c.String("token"),

			Version:         c.String("ver"),
			Branch:          c.String("branch"),
			Timeout:         c.String("timeout"),
			Sources:         c.String("sources"),
			Inclusions:      c.String("inclusions"),
			Exclusions:      c.String("exclusions"),
			Level:           c.String("level"),
			ShowProfiling:   c.String("showProfiling"),
			BranchAnalysis:  c.Bool("branchAnalysis"),
			UsingProperties: c.Bool("usingProperties"),

			PullrequestKey:    c.String("pullrequestKey"),
			PullrequestBranch: c.String("pullrequestBranch"),
			PullrequestBase:   c.String("pullrequestBase"),

			DroneRepoName: c.String("droneRepoName"),
		},
	}

	log.Println("=== plugin struct ===")
	log.Printf("%+v\n", plugin)
	log.Println("=== ENV ===")
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		log.Printf("\t * %v \n", pair)
	}
	log.Println("===========")

	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
