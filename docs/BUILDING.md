# Building and deployements (WIP)

- Clone the project: `git clone https://github.com/saferwall/saferwall`
- Using a debian linux (preferrably Ubuntu 18.04), make sure `build-essential` are installed: `sudo apt-get install build-essential curl`.
- Rename the `example.env` to `.env`. This file stores the project configuration.
- To deploy a __local instance for testing__, see [Deploying in Minikube](#Deploying-in-Minikube).
- To deploy a __prod instance on-premise__, see [Deploying on-promise](#Deploying-on-promise).
- To deploy __prod instance on the cloud__ (aws), see [Deploying on cloud](#Deploying-on-cloud).

## Deploying in Minikube

1. Install docker: `make docker-install`.
2. Install Minikube: `make minikube-install`
3. Install Kubectl: `make kubectl-install`
4. If you do not already have a hypervisor installed, install one of these now:
    - KVM, which also uses QEMU: `make kvm-install`
    - VirtualBox: `make vbox-install`
    - No Driver: minikube also supports a `--driver=none` option that runs the Kubernetes components on the host and not in a VM. Using this driver requires Docker and a Linux environment but not a hypervisor.
5. Start Minikube cluster: edit the `.env` to specify which driver to use for minikube, number of cpus, ram and disk size:
    ```c
    # supported values ['virtualbox', 'kvm2', 'none']
    export MINIKUBE_DRIVER=none
    # skip those if you set the driver to `none`.
    export MINIKUBE_CPU=2
    export MINIKUBE_MEMORY=4096
    export MINIKUBE_DISK_SIZE=40GB
    ```
    Then run `make minikube-start`.
6. Building the containers:
    - Before running any of the builds below, if you are not the _none_ driver, __make sure to eval__ the minikube env variables into your shell using the command: `eval $(minikube docker-env)`.
    - Build the backend: `make backend-build`
    - Build the frontend: `make ui-build`
    - Build the consumers: `make consumer-build`
    - Build the multiav:
        - Some AVs are not free and requires a license, you need to supply the licenses keys to be able to build the images. See [Building AV Images](#Building-AV-Images) on how to configure them.
        - multiav: `make multiav-build`
        - multiav-go: `make multiav-build-go`
7. Install Helm: `make helm-install`.
8. Add the required Helm Charts repositories: `make helm-add-repos`.
9. Fetch Helm dependecies: `make helm-update-dependency`.
10. Edit the `deployements/saferwall/values.yaml`
    - Set `nfs-server-provisioner.enabled` to true.
11. Install helm chart:
    - `cd deployement && helm install saferwall --generate-name`.
12. Wait (~1 min) till the output of `kubectl get pods` show all pods are running fine.
13. Edit again the `deployements/saferwall/values.yaml`
    - Set `couchbase-cluster.enabled` to true.
    - Set `backend.enabled` to true.
    - Set `frontend.enabled` to true.
    - Set `consumer.enabled` to true.
    - If you are  interested to see the logs in EFK:
        - Set `elasticsearch.enabled` to true.
        - Set `kibana.enabled` to true. 
        - Set `filebeat.enabled` to true.
    - Set `prometheus-operator.enabled` to true if you want to get metrics.
14. Install the chart:
    - `cd deployement` if you're not inside the deployement folder.
    - `helm upgrade saferwall <release-name generated before>`

## Deploying on cloud (AWS)

1. The first step is to build those containers. You need a repository to store them:
    - If you're deploying on the cloud, you will be probably using [docker hub](https://hub.docker.com/) or a similar service.
    - Create an account and put your credentials in the `.env` file like this:
    ```c
    # Docker hub
    export DOCKER_HUB_USR = your_docker_hub_username
    export DOCKER_HUB_PWD = your_docker_hub_password
    ```
2. Edit the file `deployments/saferwall/chsrts/nsq/values.yaml` so that:
```yaml
persistence:
    storageClass: "standard"
```


- Install it: `make saferwall`.
- Edit the deployments/values.yaml to match your needs.
- Logs are found elasticsearch:
<p align="center"><img src="https://i.imgur.com/6TnK2jR.png" width="500px" height="auto"></p>

## Deploying on-promise 

- WIP

## Building AV Images

- Edit the `.env` and fill the secrets according to which AVs you want to have.
    - Eset: copy the license to `./build/data/ERA-Endpoint.lic`, and also inside the `.env`:
        ```c
        export ESET_USER = EAV-KEYHERE
        export ESET_PWD = passwordhere
        ```
    - Avast: copy the license to `./build/data/license.avastlic`
    - Kaspersky: copy the license to `./build/data/kaspersky.license.key`
- Run `make multiav`: this will build and push to docker hub all the AVs.
- Run `make multiav-go`: this will take the images created before as a base, and build on top of them a gRPC daemon listening for files to scan.
- Every conponent of this project is running inside a docker container. This include the backend, the frontend, the AV scanners, the workers, the database, etc ...