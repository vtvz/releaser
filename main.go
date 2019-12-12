package main

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"
	"github.com/xanzy/go-gitlab"
	"os"
	"regexp"
)

type CommonConfig struct {
	GitlabProject *gitlab.Project
	ProjectConfig *Project
	RcBranchName  string
	Version       string
}

type ReleaseConfig struct {
	*CommonConfig
	ReleaseNotes string
}

type ReleaseCandidateConfig struct {
	*CommonConfig
}

type Args struct {
	ConfigFile string
	Command    string
	Project    string
	Version    string
}

func getArgs() (*Args, error) {
	var args Args

	home, _ := homedir.Dir()
	flag.StringVarP(&args.ConfigFile, "config", "c", home+"/.releaser.yaml", "Path to config")
	flag.Parse()

	if flag.NArg() != 3 {
		return nil, fmt.Errorf("Required 3 args:, <action rc|r> <project> <version in 1.0 format>")
	}

	args.Command = flag.Arg(0)
	args.Project = flag.Arg(1)
	args.Version = flag.Arg(2)

	if args.Command != "r" && args.Command != "rc" {
		return nil, fmt.Errorf("Second argument should be equal to 'r' or 'rc'")
	}

	match, _ := regexp.MatchString("^[\\d]+\\.[\\d]+$", args.Version)
	if !match {
		return nil, fmt.Errorf("Version should be in format MAJ.MIN")
	}

	fmt.Printf("%+v\n", args)

	return &args, nil
}

func main() {
	args, err := getArgs()
	if err != nil {
		panic(err)
	}

	config, err := ResolveConfig(args)
	if err != nil {
		panic(err)
	}

	project, ok := config.Projects[args.Project]

	if !ok {
		fmt.Printf("Project '%s' is not configured yet\n", args.Project)

		os.Exit(1)
	}

	fmt.Printf("Selected project is %s\n", args.Project)

	git := gitlab.NewClient(nil, project.AccessToken)

	if project.BaseUrl != "" {
		git.SetBaseURL(project.BaseUrl)
	}

	repos := project.Repos

	var gitlabProjects []*gitlab.Project
	for _, projectName := range repos {
		gitlabProject, _, err := git.Projects.GetProject(projectName, &gitlab.GetProjectOptions{})

		if err != nil {
			fmt.Printf("Error occurred on retrieving '%s' gitlab project info\n", projectName)
			panic(err)
		}

		gitlabProjects = append(gitlabProjects, gitlabProject)
	}

	var releaseNotes string
	if args.Command == "r" {
		releaseNotesBytes, err := CaptureInputFromEditor()

		if err != nil {
			panic(err)
		}

		releaseNotes = string(releaseNotesBytes)
	}

	for _, gitlabProject := range gitlabProjects {
		rcBranchName := "rc-" + args.Version + ".x"

		config := &CommonConfig{
			GitlabProject: gitlabProject,
			ProjectConfig: project,
			RcBranchName:  rcBranchName,
			Version:       args.Version,
		}

		fmt.Printf("Manage project: %s\n", gitlabProject.PathWithNamespace)

		var err error
		switch args.Command {
		case "rc":
			err = CommandReleaseCandidate(git, &ReleaseCandidateConfig{CommonConfig: config})
		case "r":
			err = CommandRelease(git, &ReleaseConfig{CommonConfig: config, ReleaseNotes: releaseNotes})
		}

		if err != nil {
			panic(err)
		}
	}
}
