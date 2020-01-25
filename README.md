<p align="center"><a href="https://saferwall.com" target="_blank" rel="noopener noreferrer"><img width="100" src="https://i.imgur.com/zjCOKPo.png" alt="Saferwall logo"></a></p>

(WIP) Saferwall is an open source malware sandbox. In the first hand, it aims to provide a collaborative platform to share samples among malware researchers. On the other hand, being a tool which acts as a system expert which generates an automated malware analysis for humans. 

## Features

- Static analysis:
    - Calculate common crypto hashes, packer identification, strings extraction, etc...).
- Multiple AV scanner which includes major antivirus vendors:

    Vendors | status 
    --- | ---
    Avast | :heavy_check_mark: 
    Avira | :heavy_check_mark: 
    Bitdefender | :heavy_check_mark: 
    ClamAV | :heavy_check_mark: 
    Comodo | :heavy_check_mark: 
    ESET | :heavy_check_mark: 
    FSecure | :heavy_check_mark: 
    Kaspersky | :heavy_check_mark: 
    MCAfee | :heavy_check_mark: 
    Sophos | :heavy_check_mark: 
    Symantec | :heavy_check_mark: 
    Windows Defender | :heavy_check_mark: 

## Installation

Saferwall take advantage of [kubernetes](https://kubernetes.io/) for its high availability, scalibility and the huge ecosystem behind it. Installing this app means deploying a kubernetes cluster. You can either deploy it on cloud or on bare metal. To make it easy to get a production grade Kubernetes cluster up and running, we use [kops](https://github.com/kubernetes/kops). Kops automated provisionning of kubernetes cluster hosted on AWS, GCE, DigitalOcean or OpenStack. For the time being, only AWS is officially supported.

1. Clone the project: `git clone https://github.com/saferwall/saferwall`
2. Using a debian linux, make sure `build-essential` are installed: `sudo apt-get install build-essential`.
3. Install it: `make saferwall` (one time)
4. The AVs requires licenses, that is why the containers are not publicly avaibable. Put licenses under `build/data ` and run: `make build`
5. Edit the deployments/values.yaml to match your needs.


## Acknowledgements

- [Fangdun Cai](https://github.com/fundon) for the awesome vue-admin dashboard.
- [horsicq](https://github.com/horsicq) for his tool [Detect It Easy](https://github.com/horsicq/Detect-It-Easy).

## Contributing

Please read [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.
