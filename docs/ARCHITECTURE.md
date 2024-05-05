# Architecture

This document describes the high-level architecture of saferwall. If you want to familiarize yourself with the architecture base, you are just in the right place!

![Architecture](./assets/architecture.svg)

## Communication between services.

- `webapis` -> json(fileScanConfig) -> `orchestrator`
- `orchestrator` -> string(sha256) -> `multiav`, `pe`, `meta`, `post-processor`, `sandbox` -> protobuf -> `aggregator`
- `post-processor` -> http -> `ml-classifier`.

## Sandbox service:

- The malware detonation service (`sandbox`) runs on **bare metal** linux nodes that hosts KVM. We call them `kvm-nodes`.
- This service is deployed as a **daemonset** that manages the virtual machines running inside KVM.
- The worker nodes that don't have support for running KVM will run other type of workloads (`multiav` , `webapi` server, etc ..).
- The `sandbox` service is using libvirt APIs to manage the state of the virtual machines. It uses the **local** `qemu:///system` connection.
  - This design is easier to work with than having to control multiple instances via a **remote** protocol which requires authentication configuration. Also dealing with service discovery for `kvm-nodes`.
- Workflow:
  1. Upon the startup of the `sandbox` service, it discovers the list of running **domains** and it stores them in its state. All VMs are marked as free and ready to receive jobs.
  2. Using a round robin approach, it assigns the job to a VM.
  3. Once the job is done, the VM is freed and restored to a clean state.
