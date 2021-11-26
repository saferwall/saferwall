# helm-chart

A Saferwall Chart for Kubernetes

## Introduction

This chart bootstraps saferwall on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites Details

* Kubernetes 1.20+

## Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add saferwall https://saferwall.github.io/helm-charts

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
saferwall` to see the charts.

To install the saferwall chart:

    helm install --namespace saferwall --create-namespace my-<chart-name> saferwall/saferwall

To uninstall the chart:

    helm delete my-<chart-name>


## Configuration

The following tables lists the configurable parameters of the saferwall chart and their default values.

|         Parameter         |           Description             |                         Default                          |
|---------------------------|-----------------------------------|----------------------------------------------------------|
| `Image`                   | Container image name              | `docker.bintray.io/jfrog/artifactory-oss`                |
| `ImageTag`                | Container image tag               | `5.2.0`                                                 |
| `ImagePullPolicy`         | Container pull policy             | `Always`                                                 |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.


