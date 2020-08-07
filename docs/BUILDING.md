# Building and deployements (WIP)

- Clone the project: `git clone https://github.com/saferwall/saferwall`
- Using a debian linux (preferrably Ubuntu 18.04), make sure `build-essential` are installed: `sudo apt-get install build-essential curl`.
- Copy the `example.env` to `.env`. This file stores the project configuration.
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
    ```mk
    # supported values ['virtualbox', 'kvm2', 'none']
    export MINIKUBE_DRIVER=none
    # skip those if you set the driver to `none`.
    export MINIKUBE_CPU=2
    export MINIKUBE_MEMORY=4096
    export MINIKUBE_DISK_SIZE=40GB
    ```
    Then run `make minikube-up`.
6. Building the containers:
    - Before running any of the builds below, if you are not using the _none_ driver, __make sure to eval__ the minikube environment variables into your shell using the command by running: `eval $(minikube docker-env)`.
    - Those are __optional__, run them only if you wish to not to use the public containers. 
        - Build the frontend: `make ui-build`.
        - Build the consumers: `make consumer-build`.
        - Build the backend `make backend-build`.
    - Build the __multiav__:
        - Some AVs are not free and requires a license, you need to supply the licenses keys to be able to build the images. See [Building AV Images](#Building-AV-Images) on how to configure them.
        - By default, saferwall will use only the free ones.
7. Install Helm: `make helm-install`.
8. Add the required Helm Charts repositories: `make helm-add-repos`.
9. Fetch Helm dependecies: `make helm-update-dependency`.
10. Edit the `deployments/saferwall/values.yaml`
    - If you are interested to see the logs in EFK:
        - Set `elasticsearch.enabled` to true.
        - Set `kibana.enabled` to true. 
        - Set `filebeat.enabled` to true.
    - Set `prometheus-operator.enabled` to true if you want to get metrics.
11. Init cert manager: `make helm-init-cert-manager`.
12. Install helm chart: `make helm-release`.
13. Wait until the output of `kubectl get pods` show all pods are running fine.
14. Edit your host file to setup a dns entry for the minikube ip:
    - `echo "$(minikube ip) mysaferwall.com api.mysaferwall.com" | sudo tee -a /etc/hosts`
15. Open the browser and naviguate to `https://mysaferwall.com`.
16. Browser will complain about self signed certificate, you will need to import the cert into your favorite browser:
    - Get the tls details from kubernetes and convert them to pkcs#12:
    ```bash
    kubectl get secret <SECRET_NAME> -o jsonpath="{.data['ca\.crt']}" | base64 --decode  >> ca.crt
    kubectl get secret <SECRET_NAME> -o jsonpath="{.data['tls\.crt']}" | base64 --decode  >> tls.crt
    kubectl get secret <SECRET_NAME> -o jsonpath="{.data['tls\.key']}" | base64 --decode  >> tls.key
    openssl pkcs12 -export -out saferwall.p12 -inkey tls.key -in tls.crt -certfile ca.crt
    ```
    - Go then to the browser certificate settings and import `saferwall.p12`.


## Deploying on cloud (AWS)

1. Install aws cli tool:
    - If you have already `aws cli` installed and credentials configured, skip to step 2.
    - Otherwise, run `make awscli-install`.
2. Install kops:
    - If you have already `kops` 1.17+ installed, skip to step 3.
    - Otherwise, run `make kops-install`.
3. Edit the file `.env` and rename the `KOPS_CLUSTER_NAME` and `KOPS_STATE_STORE` to your name of choice.
```mk
# ======================== KOPS ========================
export KOPS_CLUSTER_NAME=example.k8s.local
export KOPS_STATE_STORE=s3://kops-example-com-state-store
```
4. Edit the file `.env` again and choose the size of the cluster, the nodes count, the node size, region, etc:
```mk
export AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id)
export AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key)
export AWS_REGION = us-east-1
export AWS_NODE_COUNT = 1
export AWS_NODE_SIZE = t2.medium
export AWS_MASTER_SIZE = t2.small
export AWS_EFS_TOKEN = example-efs
```
4. Create the cluster: `make kops-cluster`.
5. The next step consists of building the docker containers used in the project:
 - Install docker: `make docker-install`.
 - You need a repository to store those images:
    - If you're deploying on the cloud, you will be probably using [docker hub](https://hub.docker.com/) or a similar service.
    - Create an account and put your credentials in the `.env` file like this: 
    ```mk
    # Docker hub
    export DOCKER_HUB_USR = your_docker_hub_username
    export DOCKER_HUB_PWD = your_docker_hub_password
    ```
    - Build and release: 
        - `make backend-release`
        - `make frontend-release`
        - `make consumer-release`
        - `make multiav-release`
        - `make multiav-release-go`
6. Install Helm: `make helm-install`.
7. Add the required Helm Charts repositories: `make helm-add-repos`.
8. Fetch Helm dependecies: `make helm-update-dependency`.
9. Init cert manager: `make helm-init-cert-manager`.
10. Edit the `deployments/saferwall/values.yaml`
    - Set `efs-provisioner.enabled` to true.
    - Set `couchbase-operator.cluster.volumeClaimTemplates.spec.storageClassName` to `default`.
    - If you are interested to see the logs in EFK:
        - Set `elasticsearch.enabled` to true.
        - Set `kibana.enabled` to true. 
        - Set `filebeat.enabled` to true.
    - Set `prometheus-operator.enabled` to true if you want to get metrics.
11. Install helm chart: `make helm-release`.

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
- Run `make multiav-build`: this will build and push to docker hub all the AVs.
- Run `make multiav-build-go`: this will take the images created before as a base, and build on top of them a gRPC daemon listening for files to scan.

 Logs are found elasticsearch:
<p align="center"><img src="https://i.imgur.com/6TnK2jR.png" width="500px" height="auto"></p>
