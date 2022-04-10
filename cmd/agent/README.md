```
sequenceDiagram
rect rgb(255, 150, 150)
Client->>Server: /Deploy saferwall sandbox.zip
Server-->>Client: Deployed: {version: 1.0.1}
end
Client->>Server: /Analyze {binary, config}
Note over Client,Server: When timeout is reached ...
Server-->>Client: Analysis finished, here are your results {apitrace.jsonl}
```