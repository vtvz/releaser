package main

import (
	"github.com/xanzy/go-gitlab"
	"fmt"
	"os"
	"regexp"
)

type CommonConfig struct {
	GitlabProject *gitlab.Project
	ProjectConfig *Project
	RcBranchName string
	Version string
}

type ReleaseConfig struct {
	*CommonConfig
	ReleaseNotes string
}

type ReleaseCandidateConfig struct {
	*CommonConfig
}

func getArgs() (project, command, version string) {
	if len(os.Args) != 4 {
        fmt.Println("Required 3 args:")
        fmt.Println("<project> <action rc|r> <version in 1.0 format>")

        os.Exit(1)
	}

	project = os.Args[1]
	command = os.Args[2]
	version = os.Args[3]

	if command != "r" && command != "rc" {
        fmt.Println("Second argument should be equal to 'r' or 'rc'")
        
        os.Exit(1)
	}

	match, _ := regexp.MatchString("^[\\d]+\\.[\\d]+$", version)
	if !match {
        fmt.Println("Version should be in format MAJ.MIN")
        
        os.Exit(1)
	}

	return
}

func main() {	
	projectKey, command, version := getArgs()

	config := ResolveConfig()

	project, ok := config.Projects[projectKey]

	if !ok {
		fmt.Printf("Project '%s' is not configured yet\n", projectKey)

        os.Exit(1)
	}

    fmt.Printf("Selected project is %s\n", projectKey)

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
	if command == "r" {
	    releaseNotesBytes, err := CaptureInputFromEditor()

	    if err != nil {
	    	panic(err)
	    }
	    
	    releaseNotes = string(releaseNotesBytes)
	}

	for _, gitlabProject := range gitlabProjects {
		rcBranchName := "rc-" + version + ".x"

		config := &CommonConfig{
			GitlabProject: gitlabProject,
			ProjectConfig: project,
			RcBranchName: rcBranchName,
			Version: version,
		}

		fmt.Printf("Manage project: %s\n", gitlabProject.PathWithNamespace)

		var err error
	    switch command {
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
