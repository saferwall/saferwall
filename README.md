<p align="center"><a href="https://saferwall.com" target="_blank" rel="noopener noreferrer"><img width="100" src="https://i.imgur.com/zjCOKPo.png" alt="Saferwall logo"></a></p>

<p align="center">
<b>Collaborative and Streamlined <ins>Threat Analysis</ins> at Scale</b>
</p>

<p align="center"> 
  <a href="https://gitter.im/saferwall/community"><img src="https://img.shields.io/gitter/room/saferwall/community?style=flat-square"></a>
  <img alt="Coverage" src="https://img.shields.io/codecov/c/github/saferwall/saferwall?style=flat-square">
  <img alt="Discord" src="https://img.shields.io/discord/803411418854064148?label=Discord&style=flat-square">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/saferwall/saferwall/Test%20Helm%20Charts?style=flat-square">
  <img alt="Downloads" src="https://img.shields.io/github/downloads/saferwall/saferwall/v0.1.0/total?style=flat-square">
  <img alt="Report Card" src="https://goreportcard.com/badge/github.com/saferwall/saferwall">
  <img alt="GitHub" src="https://img.shields.io/github/license/saferwall/saferwall?style=flat-square">
  </p>

<!-- start elevator-pitch -->

Saferwall allows you to analyze, triage and classify threats in just minutes.

<!-- end elevator-pitch -->

:star: **Collaborative** - Built for _security teams_ and _researchers_ to streamline analysis, identification and sharing malware samples.

:cloud: **Fast & cloud-native** - Scalable and cloud-native by design, deploy in minutes to bare metal or in the cloud.

:zap: **Save time** - Automate cumbersome analysis tasks, generate IoC's and reports with **zero friction**.

:package: **Batteries included** - All your favorite tools included, build intelligence feeds for hunting threats or generating signatures.

:heart: **Open source first** - We are _open-source_, _developer friendly_ and _user driven._

<p align="center"><img src="https://i.imgur.com/lYv1B4S.png" width="auto" height="auto"></p>

---

## Batteries Included

- Static Analysis:

  - File metadata, packer identification and crypto hashes.
  - String (ASCII/Unicode and ASM) extraction.
  - [PE (Portable Executable) file parser](https://github.com/saferwall/pe)

- Dynamic Analysis:

  - Automated Malware Analysis using a Hypervisor based VM.
  - Intercepting OS System Calls to build an exeuction trace of executable files.
  - Generate detailed reports and gain insight into malware behavior.
  - Choose which API's to trace, grab _screenshots_ and file changes aswell as memory dumps.

- Multiple AV scanner supporting major vendors:

  | Vendors     | status             | Vendors          | status             |
  | ----------- | ------------------ | ---------------- | ------------------ |
  | Avast       | :heavy_check_mark: | FSecure          | :heavy_check_mark: |
  | Avira       | :heavy_check_mark: | Kaspersky        | :heavy_check_mark: |
  | Bitdefender | :heavy_check_mark: | McAfee           | :heavy_check_mark: |
  | ClamAV      | :heavy_check_mark: | Sophos           | :heavy_check_mark: |
  | Comodo      | :heavy_check_mark: | Symantec         | :heavy_check_mark: |
  | ESET        | :heavy_check_mark: | Windows Defender | :heavy_check_mark: |
  | TrendMicro  | :heavy_check_mark: | DrWeb            | :heavy_check_mark: |

- Integrations with your own data processing pipeline.

## Get Started

Saferwall takes advantage of [Kubernetes](https://kubernetes.io/) for its high availability, scalability and ecosystem behind it.

Everything runs inside Kubernetes. You can either deploy it in the cloud or have it self hosted.

Here are the different deployment options available depending on how you are planning to use it:

- _"I want to try it first"_ : Use the cloud instance in [https://saferwall.com](https://saferwall.com).

- _"I want to run it locally"_ : A [Vagrant](https://www.vagrantup.com/) box is available, the only requirements are VirtualBox and Vagrant with full support
  of Windows, Linux and OSX, see [the guide](docs/DEPLOYING-TEST.md) for detailed steps.

- _"I want to make a PR or make changes"_ : When you intend to make changes to the code or make PR's, see [this guide](docs/DEPLOYING-DEV.md) for detailed steps.

- _"I love it ! I want to run it in prod"_ : First get you a [kops](https://github.com/kubernetes/kops) and check [this guide](docs/DEPLOYING-PROD.md).

The _production_ deployment using Kops automatically provisions a Kubernetes cluster hosted on AWS, GCE, DigitalOcean or OpenStack and also on bare metal. For the time being, only [AWS](https://aws.amazon.com/) is officially supported. A [helm](https://helm.sh/) chart is also provided for fast deployement. This setup works well for compagnies or small teams planning to scan a massive amounts of file.

## Our Stack:

- Golang mostly.
- Backend: [Echo](https://echo.labstack.com/)
- Frontend: [VueJS](https://vuejs.org/) + [Tailwind.css](https://tailwindcss.com/)
- Messaging: [NSQ](https://nsq.io/)
- Database: [Couchbase](https://www.couchbase.com/)
- Logging: [FileBeat](https://www.elastic.co/beats/filebeat) + [ElasticSearch](https://www.elastic.co/) + [Kibanna](https://www.elastic.co/)
- Metrics: [Prometheus](https://prometheus.io/)
- Minio: [Object storage](https://min.io/)
- Deployment: [Helm](https://helm.sh/) + [Kubernetes](https://kubernetes.io/)

## Current architecture / Workflow:

<p align="center"><img src="https://i.imgur.com/W0qXb5y.png" width="600px" height="auto"></p>

Here is a basic workflow of what happens when a new file is submitted:

- Frontend talks to the the backend via REST APIs.
- Backend uploads samples to the object storage.
- Backend pushes a message into the scanning queue.
- Consumer fetches the file and copies it to the nfs share avoiding to pull the sample on every container.
- Consumer starts scanning routines for static information such as (File metadata, File format details...)
- Consumer calls asynchronously scanning services (like AV scanners) via gRPC calls and waits for results.

## Acknowledgements

- [horsicq](https://github.com/horsicq) for [Detect It Easy](https://github.com/horsicq/Detect-It-Easy).

## Contributing

Please read [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.
