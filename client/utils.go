package client

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/yapatta/server_client_test/domain"
)

func ReadUnderlying(lines chan interface{}) {
	s := bufio.NewReader(os.Stdin)
	for {
		tStr, err := s.ReadString('\n')
		tStr = strings.TrimSuffix(tStr, "\n")
		if err != nil {
			log.Print("invalid input")
			continue
		}
		lines <- tStr
	}
}

func RpcRead(client *rpc.Client, stop chan struct{}) {
	input := make(chan interface{})
	go ReadUnderlying(input)

	for {
		lineOrErr := <-input
		line, ok := lineOrErr.(string)
		if !ok {
			close(stop)
			return
		}

		calc, err := domain.ParseStr(line)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		args := domain.RequestCalcArgs{Calc: calc}

		var reply domain.ResponseCalcArgs
		err = client.Call("StateMachine.Calc", &args, &reply)
		if err != nil {
			if strings.Contains(err.Error(), "0 division error") {
				log.Println("Calc error: ", err.Error())
				continue
			}

			if strings.Contains(err.Error(), "connection is shut down") {
				log.Print("disconnected from server")
				close(stop)
				return
			}
			log.Fatal("Calc error: ", err.Error())
		}
		fmt.Println(reply.Message)
	}
}

func ClientSetup(port string) (*rpc.Client, chan struct{}, func()) {
	client, err := rpc.Dial("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Ignore()
	signal.Notify(sigChan, syscall.SIGINT)

	stop := make(chan struct{})

	return client, stop, func() {
		for {
			select {
			case s := <-sigChan:
				if s == syscall.SIGINT {
					close(stop)
				} else {
					panic("unexpected signal received")
				}
			case <-stop:
				log.Println("client shutdown...")
				client.Close()
				return
			}
		}
	}

}
