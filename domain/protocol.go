package domain

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
)

type RequestCalcArgs struct {
	*Calc
}

type ResponseCalcArgs struct {
	Message string
}

type StateMachine struct {
	Val    int
	Client *rpc.Client
	Mu     sync.Mutex
}

func (sm *StateMachine) Calc(args *RequestCalcArgs, reply *ResponseCalcArgs) error {
	calc := args.Calc
	oldVal := sm.Val

	sm.Mu.Lock()
	defer sm.Mu.Unlock()
	switch calc.Operator {
	case OpAdd:
		sm.Val += calc.Operand
	case OpSub:
		sm.Val -= calc.Operand
	case OpMul:
		sm.Val *= calc.Operand
	case OpDiv:
		if calc.Operand == 0 {
			return fmt.Errorf("0 division error")
		}
		sm.Val /= calc.Operand
	default:
		return fmt.Errorf("undefined operator: %v", calc.Operator)
	}

	if sm.Client != nil {
		var replicatorReply ResponseCalcArgs
		if err := sm.Client.Call("StateMachine.Calc", &args, &replicatorReply); err != nil {
			log.Println(err.Error())
		}
	}
	message := fmt.Sprintf("%d = %v %d to %v", sm.Val, calc.Operator, calc.Operand, oldVal)
	log.Println(message)
	reply.Message = message

	return nil
}

func (ia *StateMachine) IdentityMapping(args *string, reply *string) error {
	*reply = *args
	return nil
}
