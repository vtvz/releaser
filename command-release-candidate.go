package main

import (
	"github.com/xanzy/go-gitlab"
	"fmt"
)

func getOrCreateBranch(git *gitlab.Client, config *ReleaseCandidateConfig) (*gitlab.Branch, bool, error) {
	branch, response, err := git.Branches.GetBranch(config.CommonConfig.GitlabProject.ID, config.CommonConfig.RcBranchName)

	if branch != nil {
		return branch, true, nil
	}

	if (response.Response.StatusCode != 404) {
		return nil, false, err
	}

	branch, response, err = git.Branches.CreateBranch(
		config.CommonConfig.GitlabProject.ID,
		&gitlab.CreateBranchOptions{
			Branch: &config.CommonConfig.RcBranchName,
			Ref: &config.CommonConfig.ProjectConfig.MainBranch,
		},
	)

	return branch, false, err
}

func CommandReleaseCandidate(git *gitlab.Client, config *ReleaseCandidateConfig) (error) {
	branch, exists, err := getOrCreateBranch(git, config)

	if err != nil {
		return err
	}

	if exists {
    	fmt.Printf("Release candidate branch already exists: %s\n", branch.Name)
	} else {
    	fmt.Printf("New release candidate branch has been created: %s\n", branch.Name)
	}

	openned := "opened"
	mergeRequests, _, err := git.MergeRequests.ListProjectMergeRequests(
		config.CommonConfig.GitlabProject.ID,
		&gitlab.ListProjectMergeRequestsOptions{
			SourceBranch: &config.CommonConfig.RcBranchName,
			State: &openned,
		},
	)

	if err != nil {
		return err
	}

	if len(mergeRequests) > 0 {
    	fmt.Printf("Merge request for release candidate branch already exists: %s\n", mergeRequests[0].WebURL)

    	return nil
	}

	mergeRequestTitle := fmt.Sprintf("Release candidate %s.x", config.CommonConfig.Version)

	master := "master"
	mergeRequest, _, err := git.MergeRequests.CreateMergeRequest(config.CommonConfig.GitlabProject.ID, &gitlab.CreateMergeRequestOptions{
		Title: &mergeRequestTitle,
		SourceBranch: &config.CommonConfig.RcBranchName,
		TargetBranch: &master,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Merge request for release candidate branch has been created: %s\n", mergeRequest.WebURL)

	return nil
} 
