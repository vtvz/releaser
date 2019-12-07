# Releaser 

This tool helps to create release branches and tags in Gitlab multi-repository projects.

## Install

Coming soon. You can clone repository and `go build`.

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

Run `releaser project1 rc 1.0` to create `rc-1.0.x` release candidate branch and merge request to `mainBranch` in all three `project1` repositories.

Run `releaser project1 r 1.0` to create `v1.0.0` release tag on `rc-1.0.x` branch in all three `project1` repositories.

Run `releaser project1 r 1.0` to create next release tag `v1.0.1`