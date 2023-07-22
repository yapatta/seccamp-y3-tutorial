# Replicated State Machine Tutorial for SecurityCamp 2023 Y3 Track

This repository is the tutorial for Seccamp 2023 pre-learning.
For now, successful execution of task 4 is confirmed.
The server can receive RPC requests from clients to change server-local state, and replicate those commands to the replication server.
Details of tasks are written [here](https://github.com/secamp-y3/tutorial.go).

## Execution

1. start replication server

```bash
go run cmd/replica/main.go
```

2. start rpc server

```bash
go run cmd/server/main.go
```

3. start rpc client

```bash
go run cmd/client/main.go
```

## Valid Command

- `Add <val>`: add val to the internal state of the server
- `Sub <val>`: subtract val to the internal state of the server
- `Mul <val>`: multiply val to the internal state of the server
- `div <val>`: divide the internal state of the server by val (receive an error when `val = 0`)
