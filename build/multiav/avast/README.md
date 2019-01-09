# Avast

The Avast Security for Linux are a set of components distributed in the form of standard software packages: avast, avast-proxy and avast-fss. This docker images contain the `avast` deb package which provides the core scanner service (avast) and a command line scan utility (scan). For more imformation, please refer to the original [vendor](https://support.avast.com/en-eu/article/131/).

## Getting Started

These instructions will cover usage information and for the docker container.

### Prerequisities


In order to run this container you'll need docker installed.

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)

### Usage

The container is running a gRPC server (port:50051) which exposes a set of APIs (like ScanFile / ScanURL /  ...). You can download a client which talks to the gRPC server here.

```
client -h

Usage:
  client [OPTIONS]

Application Options:
  -s, --scan=            Scan a file|directory|url
  -u, --update           Update VPS database
  -v, --vps-version      Get VPS version
  -p, --program-version  Get program version
  -l, --license=         Activate a license from a license key [license.avastlic]
  -e, --status-license   Check expirate of license

Help Options:
  -h, --help             Show this help message
```


You first need to acquire a [license](https://www.avast.com/linux-server-antivirus) to be able to scan files with Avast. To apply the license, run:

```
client --license=/path/to/license.avastlic
```

Scanning a directory:
```
client -s /malwares
2019/01/05 12:10:42 Scanning: /malwares/e6d87a6672ddf7bddebf582d33aa5ec5f975bd456843466282092dda37fff15b
2019/01/05 12:10:42 Scan Result: Win32:Trojan-gen
2019/01/05 12:10:42 Scanning: /malwares/96e4468c284fd9f067061b79ac3a360a0ab5d3357e5ea50609dd01aa805628c1
2019/01/05 12:10:42 Scan Result: Win32:Obfuscator-H [Trj]
2019/01/05 12:10:42 Scanning: /malwares/putty
2019/01/05 12:10:42 Scan Result: [CLEAN]
```

You can also consume the service programmatically by writing your own client in any language you wish by using the protobuf spec available [here](https://github.com/saferwall/saferwall/blob/master/api/protobuf-spec/avast.proto).





## Find Us

* [GitHub](https://github.com/saferwall)
* [Twitter](https://twitter.com/saferwall)

## Contributing

Please read [CONTRIBUTING.md](https://github.com/saferwall/saferwall/blob/master/docs/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.


## License

This project is licensed under the Apache v2 License - see the [LICENSE.md](https://github.com/saferwall/saferwall/blob/master/LICENSE) file for details.
