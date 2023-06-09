# GitLab Adapter

[![Build status](https://github.com/ai-lenok/gitlab-adapter/actions/workflows/workflow.yml/badge.svg?branch=main)](https://github.com/ai-lenok/gitlab-adapter/actions/workflows/workflow.yml)

Adapter between Learning Management System and GitLab

## Features

* Create repository
* Delete repository
* Verify that the latest build in the repository was successful

## Ways to work

* HTTP server
* Command Line Interface

### HTTP examples

[RESTful client](test-http-client/client.http)

### Health check

```shell
GET http://{{host}}/health
```

## Configuration

All necessary configs you can put to `config/application.yaml` file

Or put like environment variables with prefix `GA_`

Or put on command line parameters

### Example

#### Config file

```yaml
gitlab:
  host: "https://gitlab.com"
  token: "change-me"
```

#### Environment variables

```shell
export GA_GITLAB_HOST="https://gitlab.com"
export GA_GITLAB_TOKEN="change-me"
```

#### Command line parameters

```shell
gitlab-adapter start-server --gitlab.host "https://gitlab.com" --gitlab.token "change-me"
```

## Work

### Create repository

### CLI

```shell
gitlab-adapter create-repo --namespace 1234567890 --name test-name --path test-path --description "Description repo"
```

### HTTP

```shell
POST http://{{host}}/api/v1/project
```

### Delete repository

### CLI

```shell
gitlab-adapter delete-repo --project-id 1000000000
```

### HTTP

```shell
DELETE http://{{host}}/api/v1/project
```

### Verify that the latest build in the repository was successful

### CLI

```shell
gitlab-adapter verify-pipeline --project-id 1000000000
```

* Return status 0 - success
* Return status 1 - failed

### HTTP

```shell
POST http://{{host}}/api/v1/project/verify-pipeline
```

* Return status 204 - success
* Return status 409 - failed

## App

### Build

```shell
make build
```

### Remove all artefacts

```shell
make clean
```

### Run server

```shell
make start-server
```

## Docker

### Build

```shell
docker image build --tag IMAGE_NAME .
```

### Run

```shell
docker container run --publish 8080:8080 --volume ./config/:/app/config/ ghcr.io/ai-lenok/gitlab-adapter:main
```

#### Docker-compose

```shell
docker-compose up --detach
```

## Kubernetes

### Run

```shell
kubectl apply --filename manifest/
```

### Clean

```shell
kubectl delete --filename manifest/
```