package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func HandleListner(serverCtx context.Context, ln net.Listener, rpcServer *rpc.Server, hlwg chan interface{}, wg *sync.WaitGroup) {
	defer close(hlwg)

	for {
		conn, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Timeout() {
					log.Print("Connection timeout")
					continue
				}
			}
			if errors.Is(err, net.ErrClosed) {
				<-serverCtx.Done()
				log.Print("replicator listner closing...")
				return
			}
		}

		wg.Add(1)
		go func() {
			defer func() {
				conn.Close()
				log.Printf("disconnected: %v", conn.RemoteAddr().String())
				wg.Done()
			}()

			sc := make(chan struct{}, 1)
			go func() {
				defer close(sc)
				log.Printf("connected: %v", conn.RemoteAddr().String())
				rpcServer.ServeConn(conn)
			}()

			select {
			case <-serverCtx.Done():
			case <-sc:
			}
		}()
	}
}

func ServerSetup(rpcServer *rpc.Server, port string) (context.Context, net.Listener, *rpc.Server, chan interface{}, *sync.WaitGroup, func()) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Ignore()
	signal.Notify(sigChan, syscall.SIGINT)

	wg := sync.WaitGroup{}
	hwc := make(chan interface{}, 1)

	serverCtx, shutdown := context.WithCancel(context.Background())

	return serverCtx, ln, rpcServer, hwc, &wg, func() {
		s := <-sigChan

		switch s {
		case syscall.SIGINT:
			log.Println("server shutdown...")
			shutdown()
			ln.Close()
			<-hwc
			wg.Wait()
		default:
			panic("unexpected signal received")
		}
	}
}
