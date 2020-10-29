
# Deploying on AWS

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
9. Init cert manager: `make k8s-init-cert-manager`.
10. Edit the `deployments/saferwall/values.yaml`
    - Set `efs-provisioner.enabled` to true.
    - Set `couchbase-operator.cluster.volumeClaimTemplates.spec.storageClassName` to `default`.
    - If you are interested to see the logs in EFK:
        - Set `elasticsearch.enabled` to true.
        - Set `kibana.enabled` to true. 
        - Set `filebeat.enabled` to true.
    - Set `prometheus-operator.enabled` to true if you want to get metrics.
11. Install helm chart: `make helm-release`.

## Tips for deploying a production cluster

- To have a HA cluster we need at least more than one master and several workers, in different availability zones.
- With multiple master nodes, you will be able both to do graceful (zero-downtime) upgrades and you will be able to survive AZ failures.
- Harden Kubernetes API to be accessible only from alllowed IPs via firewall rules or change the kops cluster config with the `kubernetesApiAccess`.

## Billing Tips

- Opt for EC2 Spot Instances and reserved instances.
- Support kubernetes cluster with on-demand instances, which can take up the slack in the event of any interruptions to spot instances. This will improve availability and reliability.
- Size of the master node:
	- 1-5 nodes: m3.medium
	- 6-10 nodes: m3.large
	- 11-100 nodes: m3.xlarge
	- 101-250 nodes: m3.2xlarge
	- 251-500 nodes: c4.4xlarge
	- more than 500 nodes: c4.8xlarge
