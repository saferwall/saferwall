# Saferwall CLI

A CLI tool to download or scan files using the command line.

## Usage

CLI tool requires authentification through environment variables during first execution.

### Authentification

Before being able to download or scan files you need to authenticate yourself.

```sh

cli auth "username" "password"

```

### Download

You can download files using their SHA256 hash.

```sh

cli download 0001cb47c8277e44a09543291d95559886b9c2da195bd78fdf108775ac91ac53

```