# Saferwall CLI

A CLI tool to use Saferwall to download samples, scan or re-scan new samples.

## Usage

To use the CLI tool you need a [Saferwall](https://saferwall.com) account in order to authenticate
yourself.

The CLI tool reads your username and password from a local ```.env``` file.

```sh

SAFERWALL_AUTH_USERNAME=username
SAFERWALL_AUTH_PASSWORD=password

```

### Download

You can download files using their SHA256 hash and specify an output folder, you can also download a batch of samples by copying their SHA256 hash to the clipboard.

```sh

cli download 0001cb47c8277e44a09543291d95559886b9c2da195bd78fdf108775ac91ac53 -o tmp/

```

### Scan or Rescan

You can scan or rescan files using the scan command.

```sh

cli [re]scan /samples/putty.exe

```

For the rest of the commands you can use ```help``` to show a usage guide.

```sh

A cli tool to interfact with saferwall APIs (scan, rescan, upload, ...)

Usage:
  sfwcli [flags]
  sfwcli [command]

Available Commands:
  download    Download file
  help        Help about any command
  rescan      Resccan file
  s3upload    S3 upload
  scan        Scan file
  version     Vesion number

Flags:
  -h, --help   help for sfwcli

```
