# Charly
![charly](https://github.com/cooperspencer/charly/blob/main/charly.png)

With `charly` you can do specific tasks when a repository gets a new commit.

You can monitor repositories from:
- Github
- Gitlab
- Gitea
- Gogs

VCS I want to add:
- BitBucket
- GitBucket
- OneDev

## How to make a configuration file
[Here is an example](https://github.com/cooperspencer/charly/blob/main/conf.example.yaml)

## How to run the binary version
`./charly path-to-conf.yml`

## How to compile
`go build .`

## Scripts
For every script run the following variables are set as environment variables:
- COMMIT
- BRANCH 
- URL 
- SSHURL
- USER
- REPO