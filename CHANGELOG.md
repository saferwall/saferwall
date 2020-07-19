# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.2] - 2020-08-01
### Added
- Add a cmd tool to batch upload files.
- Add a pkg for parsing portable executable files.
- Add a UI for displaying PE parser results.
- Add s3upload pkg to simplify mass-uploading of files into s3.
- Add upload pkg to simplify uploading a local db into saferwall.

### Changed
- Add nfs-server-provisionner for local testing in minikube.
- Improve the building process documentation thanks to [Jameel Haffejee](https://github.com/RC114).
- Reworked to file tags schema.
- Improve the rendering of the landing page.
- Fix phrasing in README from [@bf](https://github.com/bf).
- Fix recover from panic routine in parse-pe in consumer.
- Add exception catching in strings pkg.
- Enable Kibana / ElasticSearch / FileBeat in helm deployments.
- Add ConextLogger in consumer to always log sha256.


## [0.0.1] - 2020-03-09
### Added
- Initiale release includes a multi-av scanner + strings + file metadata.
- UI with options to download, rescan, like a sample and share comments.
- User profile to track submissions, followers and see activities.
