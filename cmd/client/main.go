package main

import "github.com/yapatta/server_client_test/client"

func main() {
	c, stop, clientShutdown := client.ClientSetup(":8000")

	go client.RpcRead(c, stop)
	clientShutdown()
}
