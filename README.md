<p align="center"><a href="https://saferwall.com" target="_blank" rel="noopener noreferrer"><img width="100" src="https://i.imgur.com/zjCOKPo.png" alt="Saferwall logo"></a></p>
<h4 align="center">Saferwall is an open source malware analysis platform</a>.</h4>

<p align="center"> 
  <a href="https://gitter.im/saferwall/community"><img src="https://img.shields.io/gitter/room/saferwall/community?style=flat-square"></a>
</p>

It aims for the following goals:
- Provide a collaborative platform to share samples among malware researchers.
- Acts as a system expert, to help researchers generates an automated malware analysis report.
- Hunting platform to find new malwares.
- Quality ensurance for signature before releasing.

<p align="center"><img src="https://i.imgur.com/lYv1B4S.png" width="700px" height="auto"></p>

## Features

- Static analysis:
    - Crypto hashes, packer identification
    - Strings extraction
- Multiple AV scanner which includes major antivirus vendors:

    Vendors | status | Vendors | status
    --- | --- | --- | ---
    Avast | :heavy_check_mark: | FSecure | :heavy_check_mark: 
    Avira | :heavy_check_mark: | Kaspersky | :heavy_check_mark: 
    Bitdefender | :heavy_check_mark: | McAfee | :heavy_check_mark: 
    ClamAV | :heavy_check_mark: | Sophos | :heavy_check_mark: 
    Comodo | :heavy_check_mark: | Symantec | :heavy_check_mark: 
    ESET | :heavy_check_mark: | Windows Defender | :heavy_check_mark: 

## Installation

Saferwall take advantage of [kubernetes](https://kubernetes.io/) for its high availability, scalibility and the huge ecosystem behind it. 

Everything runs inside Kubernetes. You can either deploy it in the cloud or have it self hosted.

 To make it easy to get a production grade Kubernetes cluster up and running, we use [kops](https://github.com/kubernetes/kops). It automatically provisions a kubernetes cluster hosted on AWS, GCE, DigitalOcean or OpenStack and also on bare metal. For the time being, only AWS is officially supported.

Steps:
(This still needs to be improved)
1. Clone the project: `git clone https://github.com/saferwall/saferwall`
2. Using a debian linux, make sure `build-essential` are installed: `sudo apt-get install build-essential`.
3. Install it: `make saferwall`.
4. The AVs requires licenses, that is why the containers are not publicly avaibable. Rename the `example.env` to `.env` and fill the credentials according to which AVs you want to have. Also, put licenses under `build/data` and run: `make build`.
5. Edit the deployments/values.yaml to match your needs.

## Built with:

- Golang mostly.
- Backend: [Echo](https://echo.labstack.com/)
- Frontend: [VueJS](https://vuejs.org/) + [Bulma](https://bulma.io/)
- Messaging: [NSQ](https://nsq.io/)
- Database: [Couchbase](https://www.couchbase.com/)
- Logging: [FileBeat](https://www.elastic.co/beats/filebeat) + [ElasticSearch](https://www.elastic.co/) + [Kibanna](https://www.elastic.co/)
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
