<p align="center"><a href="https://saferwall.com" target="_blank" rel="noopener noreferrer"><img width="100" src="https://i.imgur.com/zjCOKPo.png" alt="Saferwall logo"></a></p>
<h4 align="center">Saferwall is an open source malware analysis platform</a>.</h4>

<p align="center"> 
  <a href="https://gitter.im/saferwall/community"><img src="https://img.shields.io/gitter/room/saferwall/community?style=flat-square"></a>
  <img alt="Discord" src="https://img.shields.io/discord/803411418854064148?label=Discord&style=flat-square"> 
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/saferwall/saferwall/Test%20Helm%20Charts?style=flat-square">
  <img alt="Report Card" src="https://goreportcard.com/badge/github.com/saferwall/saferwall">
  <img alt="GitHub" src="https://img.shields.io/github/license/saferwall/saferwall?style=flat-square">
  </p>

It aims for the following goals:
- Provide a collaborative platform to share samples among malware researchers.
- Acts as a system expert, to help researchers generates an automated malware analysis report.
- Hunting platform to find new malwares.
- Quality ensurance for signatures before releasing.

<p align="center"><img src="https://i.imgur.com/lYv1B4S.png" width="auto" height="auto"></p>

## Features

- Static analysis:
    - Crypto hashes, packer identification,
    - Strings extraction
    - [PE (Portable Executable) file parser](https://github.com/saferwall/pe)
- Multiple AV scanner which includes major antivirus vendors:

    Vendors | status | Vendors | status
    --- | --- | --- | ---
    Avast | :heavy_check_mark: | FSecure | :heavy_check_mark: 
    Avira | :heavy_check_mark: | Kaspersky | :heavy_check_mark: 
    Bitdefender | :heavy_check_mark: | McAfee | :heavy_check_mark: 
    ClamAV | :heavy_check_mark: | Sophos | :heavy_check_mark: 
    Comodo | :heavy_check_mark: | Symantec | :heavy_check_mark: 
    ESET | :heavy_check_mark: | Windows Defender | :heavy_check_mark: 
    TrendMicro | :heavy_check_mark: | DrWeb | :heavy_check_mark: 

## Installation

Saferwall take advantage of [kubernetes](https://kubernetes.io/) for its high availability, scalability and the huge ecosystem behind it. 

Everything runs inside Kubernetes. You can either deploy it in the cloud or have it self hosted.

Here are the different deployment options available depending on how you are planning to use it:
- Just to get a feeling of the app, you can use the already hosted instance in [https://saferwall.com](https://saferwall.com).
- For __local testing__ purposes or __individual__ usage, a [Vagrant](https://www.vagrantup.com/) box is available, the only requirements is virtualbox and vagrant. This setup runs on Windows, Linux and OSX. Please refer to this [link](docs/DEPLOYING-TEST.md) for detailed steps.
- In __development scenarios__, when you intend to make changes to the code and add features, please refer to this [link](docs/DEPLOYING-DEV.md) for detailed steps.
- For __production grade deployment__, we use [kops](https://github.com/kubernetes/kops). It automatically provisions a kubernetes cluster hosted on AWS, GCE, DigitalOcean or OpenStack and also on bare metal. For the time being, only [AWS](https://aws.amazon.com/) is officially supported. A [helm](https://helm.sh/) chart is also provided for fast deployement. This work well for compagnies or small teams planning to scan a massive amounts of file. Please refer to this [link](docs/DEPLOYING-PROD.md) for detailed steps.


## Built with:

- Golang mostly.
- Backend: [Echo](https://echo.labstack.com/)
- Frontend: [VueJS](https://vuejs.org/) + [Bulma](https://bulma.io/)
- Messaging: [NSQ](https://nsq.io/)
- Database: [Couchbase](https://www.couchbase.com/)
- Logging: [FileBeat](https://www.elastic.co/beats/filebeat) + [ElasticSearch](https://www.elastic.co/) + [Kibanna](https://www.elastic.co/)
- Metrics: [Prometheus](https://prometheus.io/)
- Minio: [Object storage](https://min.io/)
- Deployment: [Helm](https://helm.sh/) + [Kubernetes](https://kubernetes.io/)

## Current architecture / Workflow:

<p align="center"><img src="https://i.imgur.com/W0qXb5y.png" width="600px" height="auto"></p>

Here is a basic workflow which happens during a file scan:
- Frontend talks to the the backend via REST APIs.
- Backend uploads samples to the object storage.
- Backend pushes a message into the scanning queue.
- Consumer fetches the file and copy it into to the nfs share avoiding to pull the sample on every container.
- Consumer calls asynchronously scanning services (like AV scanners) via gRPC calls and waits for results.

## Acknowledgements

- [Fangdun Cai](https://github.com/fundon) for the awesome vue-admin dashboard.
- [horsicq](https://github.com/horsicq) for his tool [Detect It Easy](https://github.com/horsicq/Detect-It-Easy).

## Contributing

Please read [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.
