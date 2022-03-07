# Charly
![charly](https://github.com/cooperspencer/charly/blob/main/charly.png)

With `charly` you can do specific tasks when a repository gets a new commit.

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
- USER
- REPO
- TOKEN
- GIT_PWD
- SSHKEYFILE
- SSHKEYPWD
- USERNAME
- REPO