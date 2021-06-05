# Contributing

When contributing to this repository, please first discuss the change you wish to make via issue, email, or any other method with the owners of this repository before making a change. 

Please note we have a code of conduct, please follow it in all your interactions with the project.

# Table of Contents

- [Repostiory Layout](#Repostiory-Layout)
- [Requirements](#Requirements)
- [Developing on the backend](#Developing-on-the-backend)
- [Developing on the frontend](#Developing-on-the-frontend)

## Repostiory Layout
* __api__ : proto buffer specs, swagger manifests.
* __build__ : docker files, makefiles, packer scripts.
* __cmd__: main applications for this project.
* __configs__: Configuration file templates or default configs.
* __docs__: design and user documents (in addition to your godoc generated documentation).
* __deployments__: helm chart.
* __pkg__ : core package implementations and consumer instance.
* __scripts__: scripts to perform any build, install, analysis, etc operations.
* __test__: test data, (the tests are found on the location as the go code).

## Requirements

- Copy the `example.env` to `.env`. This file stores the project configuration.
- Nearly all operations that we work with daily in this project, is automated using __makefiles__. You can print thefull list of supported commands by typing `make` on the root directory:

```shell
 make
help                           This help.
awscli-install                 Install aws cli tool and configure it
docker-build                   Build the container
docker-build-nc                Build the container without caching
docker-run                     Run container on port configured in `config.env`
...
protobuf-generate-api          Generates protocol buffers definitions files. 
protobuf-generate-api-all      Generates protocol buffers definitions files for all AVs.
protobuf-install-compiler      Install protobuf compiler
protobuf-protoc-gen-go         Install protoc plugin for Go
certbot-install                Install Certbot for nginx
elastic-drop-db                Delete all indexes.
```

Your commits **should** follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.
