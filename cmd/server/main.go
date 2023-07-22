package main

import (
	"log"
	"net/rpc"

	"github.com/yapatta/server_client_test/domain"
	"github.com/yapatta/server_client_test/server"
)

func main() {
	stateMachine := new(domain.StateMachine)
	rpcServer := rpc.NewServer()
	rpcServer.Register(stateMachine)

	c, err := rpc.Dial("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	stateMachine.Client = c
	defer c.Close()

	serverCtx, ln, rpcServer, hwc, wg, serverShutdown := server.ServerSetup(rpcServer, ":8000")

	go server.HandleListner(serverCtx, ln, rpcServer, hwc, wg)

	serverShutdown()
}
