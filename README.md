# GitLab Adapter

Adapter between Learning Management System and GitLab

## Features

* Create repository
* Delete repository
* Verify that the latest build in the repository was successful

## Ways to work

* HTTP server
* Command Line Interface

## Configuration

All necessary configs you can put to `config/application.yaml` file

Or put on command line parameters

## Work

### Create repository

### CLI

```shell
gitlab-adapter create-repo --namespace 1234567890 --name test-name --path test-path --description "Description repo"
```

### Delete repository

### CLI

```shell
gitlab-adapter delete-repo --project-id 1000000000
```

### Verify that the latest build in the repository was successful

### CLI

```shell
gitlab-adapter verify-pipeline-status --project-id 1000000000
```

* Return status 0 - success
* Return status 1 - failed