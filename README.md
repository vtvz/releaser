# Releaser 

This tool helps to create release branches and tags in Gitlab multi-repository projects.

## Problem

Imagine you have a large project with multiple repositories:

- frontend
- backend
- service
- landing
- etc...

Your release flow is the following:

- create release candidate branch (`rc-1.0.x`) in each repository
- stabilize all components
- create release tag for release branches in each repository (`v1.0.0`)
- deploy all components

There already are 12 manual actions for 4 repositories!

And after some time you find a little critical bug on production.

- fix this bug
- create incremented release tag for release branches in each repository (`v1.0.1`)
- deploy all components

Another 8 manual actions. There are too many manual jobs that can be forgotten or done wrong.

This little tool helps to create release branches and tags in all project repositories in one command (simultaneous deployment will be in the future).

Just run `releaser rc projectName 1.0` to create `rc-1.0.x` release candidate branches in all repositories.

And when you feel this candidate is stable enough to be deployed to production, run `releaser r projectName 1.0` and `1.0.0` release tags will be created in each repository on `rc-1.0.x` bracnhes.

If current release has a problem which needs to be fixed? just apply changes to `rc-1.0.x` branch and after testing and further stabilization run `releaser r projectName 1.0` command to create brand new tags with increased patch version `1.0.1`.

## Download

Download bin-file for your OS:
- [Linux amd64](./bin/releaser-linux-amd64)
- [MacOS amd64](./bin/releaser-darwin-amd64)

## Install
### Linux
```bash
sudo curl -L https://raw.githubusercontent.com/vtvz/releaser/master/bin/releaser-linux-amd64 -o /usr/local/bin/releaser
sudo chmod +x /usr/local/bin/releaser
```

### MacOS
```bash
curl -L https://raw.githubusercontent.com/vtvz/releaser/master/bin/releaser-darwin-amd64 -o /usr/local/bin/releaser
chmod +x /usr/local/bin/releaser
```

## Usage

Make `.releaser.yaml` file in home directory with similar content:

```yaml
# Access token can be obtained here (api scope)
# https://gitlab.axmit.com/profile/personal_access_tokens
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
