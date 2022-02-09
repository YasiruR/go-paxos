# Go-Paxos

Go-Paxos is a golang implementation of the paxos algorithm which is introduced by famous 
computer scientist Leslie Lamport to reach consensus among a set individual nodes in 
distributed systems. This implementation has extended the original proposed solution into 
what is known as multi-paxos thus supporting multiple rounds of decisions.

## Usage

#### To build from the source code
1. Run ``go build -o run`` in the parent directory
2. Move the executable file (run) and configs.yaml to relevant nodes in the cluster
3. Update configurations if required
   1. `leader_http_timeout`: Timeout of proposer waiting for responses from acceptors (in seconds)
   2. `replica_http_timeout`: Timeout of replica waiting for the requested leader (in seconds)

#### To execute

`./<program> <replica or leader> <host:port> <list of leaders> <list of replicas>`

Eg: 
1. As a leader: ./run leader localhost:2022 localhost:2023,localhost:2024 localhost:2025,localhost:2026
2. As a replica: ./run replica localhost:2025 localhost:2022,localhost:2023,localhost:2024 localhost:2026

## Tester

Testing scripts are included in the `scripts` directory

## Logging

A logger library is integrated for debugging purposes with following
hierarchical log levels. Required level can be enabled via `log_level` in configs.yaml
by setting the relevant value.

1. `ERROR`
2. `INFO`
3. `DEBUG`
4. `TRACE`