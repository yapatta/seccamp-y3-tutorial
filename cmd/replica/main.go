package main

import (
	"net/rpc"

	"github.com/yapatta/server_client_test/domain"
	"github.com/yapatta/server_client_test/server"
)

func main() {
	replicator := new(domain.StateMachine)
	rpcServer := rpc.NewServer()
	rpcServer.Register(replicator)

	serverCtx, ln, rpcServer, hwc, wg, serverShutdown := server.ServerSetup(rpcServer, ":9000")

	go server.HandleListner(serverCtx, ln, rpcServer, hwc, wg)

	serverShutdown()
}
