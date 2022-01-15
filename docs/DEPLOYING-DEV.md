# Deploying documentation for developers

Ubuntu `18 or 20` users will benefit from `make` commands available for quickly bootstrapping this deployment.
- Install few dependencies: `sudo apt-get install make curl git -y`.
- Clone the project: `git clone https://github.com/saferwall/saferwall`
- Copy the `example.env` to `.env`, this file stores the project configuration.

# Deploying in a VM vs on a physical machine

- There are 2 options:
  - If you don't want to run a VM, as it takes more resources and not ideal for a single developer machine, use [kind](https://kind.sigs.k8s.io/).
  - If you don't mind spinning a VM, or you want the most full-featured local Kubernetes solution, then you can go ahead with [minikube](https://minikube.sigs.k8s.io/docs/).

It is __recommanded__ to go with `kind` if you don't know know which one to choose.

# Deploying in Kind or Minikube

1. Install `Docker`: `make docker-install`.
2. Install `Kind`: `make kind-install` or Minikube: `make minikube-install`
3. Install `Kubectl`: `make kubectl-install`
4. __Minikube__ users only:
    - A hypervisor like `QEMU/KVM` or `Virtualbox` is required: For KVM/QEMU: `make kvm-install`, for VirtualBox: `make vbox-install`.
    - Edit the `.env` to specify which driver to use and number of cpus, ram and disk size:
        ```mk
        # supported values ['virtualbox', 'kvm2']
        export MINIKUBE_DRIVER=virtualbox
        export MINIKUBE_CPU=4
        export MINIKUBE_MEMORY=8192
        export MINIKUBE_DISK_SIZE=40GB
        ```
5. _Optional_ step: build the paid anti-virus docker images.
    - Some AVs are not free and requires a license, you need to supply the licenses keys to be able to build the images. See [Building AV Images](#Building-AV-Images) on how to configure them.
    - By default, saferwall will use only the free ones.
6. Create kind cluster: `make kind-up` or minikube cluster: `make minikube-up`.
7. Install Helm: `make helm-install`.
8. Edit the `deployments/saferwall/values.yaml`
    - If you are interested to see the logs in EFK:
        - Set `elasticsearch.enabled` to true.
        - Set `kibana.enabled` to true.
        - Set `filebeat.enabled` to true.
    - Set `prometheus-operator.enabled` to true if you want to get metrics.
9. Install helm chart: `make helm-release`.
10. Wait until the output of `watch | kubectl get pods` show all pods are running fine.
11. Edit your host file to setup a dns entry for for the services running inside the cluster:
    - Minikube: `echo "$(minikube ip) mysaferwall.com api.mysaferwall.com" | sudo tee -a /etc/hosts`
    - Kind: `echo "127.0.0.1 mysaferwall.com api.mysaferwall.com" | sudo tee -a /etc/hosts`
12. Open the browser and naviguate to `mysaferwall.com` and `api.mysaferwall.com` and add an certificate exception for both domains.
13. The credentials to access the dashboard are:
    ```
    username: Administrator
    password: password
    ```

# Building AV Images

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
