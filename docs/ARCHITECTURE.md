# Architecture

This document describes the high-level architecture of saferwall. If you want to familiarize yourself with the architecture base, you are just in the right place!

![Architecture](./assets/architecture.svg)

## KVM Nodes

- The malware detonation service runs on bare metal linux nodes that hosts KVM. We call them `kvm-nodes`.
- Each kvm-node is part of the kubernetes cluster, they belong to the same kubernetes `instance group`.
- Each kvm-node runs a daemonset that manages the virtual machines running inside KVM.

## Communication between services.

- `webapis` -> json(fileScanConfig) -> `orchestrator`
- `orchestrator` -> string(sha256) -> `multiav`, `pe`, `meta`, `post-processor` -> protobuf -> `aggregator`
- `post-processor` -> http -> `ml-classifier`.
- `orchestrator` -> json(fileScanConfig) -> `sandbox`
