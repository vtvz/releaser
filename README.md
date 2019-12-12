# Releaser 

This tool helps to create release branches and tags in Gitlab multi-repository projects.

## Problem

Imagine you have a large project with multiple repositories:

- frontend
- backend
- service
- landing
- etc...

Your release flow is folowing:

- create release candidate branch (`rc-1.0.x`) in each repository
- stabilize all components
- create release tag on release branches in each repository (`v1.0.0`)
- deploy all components

There is already 12 manual actions for 4 repository!

And after a some time later you find a little critical bug on production.

- fix this bug
- create incremented release tag on release branches in each repository (`v1.0.1`)
- deploy all components

Another 8 manual actions. There is too many manual jobs to be forgotten or done wrong.

This little tool helps to create release branches and tags in all project repositories in one command (simultaneous deployment will be in future).

Just run `releaser rc projectName 1.0` to make `rc-1.0.x` release candidate branches in all repositories.

And when you feel this candidate is stable enough to be deployed to production run `releaser r projectName 1.0` and `1.0.0` release tags will be created in each repository on `rc-1.0.x` bracnhes.

If current release has a problem needed to be fixed just apply changes to `rc-1.0.x` branch and after testing and further stabilization run `releaser r projectName 1.0` command to create brand new tags with increased patch version `1.0.1`.

## Install

Coming soon. You can clone repository and `go build`. Or ask me to send you binary file.

## Usage

Make `.releaser.yaml` file in home directory with similar content:

```yaml
accessToken: "<default_secret>"
baseUrl: https://gitlab.local.com/

projects:
  project1:
    accessToken: "<another_secret>"
    baseUrl: https://gitlab.com/
    repos:
      - project/client
      - project/api
      - project/service
  project2:
    mainBranch: developer
    repos:
      - project2/client
      - project2/api
      - project2/service

```

Run `releaser rc project1 1.0` to create `rc-1.0.x` release candidate branch and merge request to `mainBranch` in all three `project1` repositories.

Run `releaser r project1 1.0` to create `v1.0.0` release tag on `rc-1.0.x` branch in all three `project1` repositories.

Run `releaser r project1 1.0` to create next release tag `v1.0.1`
