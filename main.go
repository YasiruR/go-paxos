package main

import (
	"fmt"
	"github.com/go-paxos/domain"
	"github.com/go-paxos/logger"
	"github.com/go-paxos/roles"
	"github.com/go-paxos/server"
	"github.com/google/uuid"
	traceableContext "github.com/tryfix/traceable-context"
	"os"
	"strconv"
	"strings"
)

const (
	typeReplica = `replica`
	typeLeader  = `leader`
)

func main() {
	args := os.Args
	if len(args) != 5 {
		fmt.Println(`command should be in the form of ./<program> <replica or leader> <host:port> <list of leaders> <list of replicas \n
			eg: /run leader localhost:2022 localhost:2023,localhost:2024 localhost:2025,localhost:2026`)
	}

	ctx := traceableContext.WithUUID(uuid.New())
	logger.SetConfigs(ctx)
	logger.Init(ctx)
	domain.SetConfigs(ctx)

	p := port(args[2])
	leaders, replicas := hosts(args[3]), hosts(args[4])

	var replica *roles.Replica
	var leader *roles.Leader
	if args[1] == typeReplica {
		replica = roles.NewReplica(args[2], leaders)
	} else if args[1] == typeLeader {
		leader = roles.NewLeader(args[2], leaders, replicas)
	}

	server.Init(ctx, p, leader, replica)
}

func port(host string) int {
	list := strings.Split(host, `:`)
	if len(list) != 2 {
		panic(`hostname should be in the form of <hostname:port>`)
	}

	p, err := strconv.Atoi(list[1])
	if err != nil {
		panic(err)
	}

	return p
}

func hosts(arg string) []string {
	list := strings.Split(arg, `,`)
	if len(list) == 0 {
		panic(`leader/replica list is empty`)
	}

	return list
}
