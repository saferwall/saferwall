# Sandbox

## Agent client-server communication

- gRPC client/server app to bootstrap a binary sample for the malware sandbox analysis.

```mermaid
sequenceDiagram
rect rgb(255, 150, 150)
Client->>Server: /Deploy saferwall sandbox.zip
Server-->>Client: Deployed: {version: 1.0.0}
end
Client->>Server: /Analyze {binary, config}
Note over Client,Server: When timeout is reached ...
Server-->>Client: Analysis finished, here are your results {apitrace.jsonl}
```

## Workflow

- Find a free VM
- Make an RPC `/deploy` packages.
- Make an RPC `/analyze` sample.
- Wait for the given timeout until results come back.
- Process the results:
    - Convert the APIs trace log and the sandbox log from jsonl to json.
    - Go over screenshots, generate thumbs and upload them.
    - Go over artifacts+dumps, yara scan them and upload the results.

## References

- https://access.redhat.com/solutions/732773
- https://www.redhat.com/en/blog/introduction-virtio-networking-and-vhost-net
