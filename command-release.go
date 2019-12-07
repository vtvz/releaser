package main

import (
	"github.com/xanzy/go-gitlab"
	"fmt"
	"net/http"
)

func hasTag(possibleTag string, tags []*gitlab.Tag) (bool) {
	for _, tag := range tags {
		if possibleTag == tag.Name {
			return true
		}
	}

	return false
}

func CommandRelease(git *gitlab.Client, config *ReleaseConfig) (error) {
	branch, _, err := git.Branches.GetBranch(config.CommonConfig.GitlabProject.ID, config.CommonConfig.RcBranchName)

	if err != nil {
		return err
	}

	search := func(req *http.Request) (error) {
		q := req.URL.Query()
		q.Add("search", "^v" + config.CommonConfig.Version)

		req.URL.RawQuery = q.Encode()

		return nil
	}

	tags, _, err := git.Tags.ListTags(config.CommonConfig.GitlabProject.ID,  &gitlab.ListTagsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}, search)

	if err != nil {
    	return err
	}

	createTagVersion := func(patchVersion int) (string) {
		return fmt.Sprintf("v%s.%d", config.CommonConfig.Version, patchVersion);
	}

	patchVersion := 0
	for hasTag(createTagVersion(patchVersion), tags) {
		patchVersion++
    }

    tagVersion := createTagVersion(patchVersion)

    tag, _, err := git.Tags.CreateTag(config.CommonConfig.GitlabProject.ID, &gitlab.CreateTagOptions{
    	TagName: &tagVersion,
    	Ref: &branch.Name,
    	ReleaseDescription: &config.ReleaseNotes,
    })

    if err != nil {
    	return err
    }

    fmt.Printf("New release tag has been created: %s \n", tag.Name)

    return nil
} 
