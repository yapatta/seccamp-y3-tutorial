package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type Op string

const (
	OpAdd = Op("Add")
	OpSub = Op("Sub")
	OpMul = Op("Mul")
	OpDiv = Op("Div")
)

type Calc struct {
	Operator Op
	Operand  int
}

func parseError(format string, args ...any) error {
	return fmt.Errorf("parse error: %v", fmt.Sprintf(format, args...))
}

func ParseStr(input string) (*Calc, error) {
	arr := strings.Split(input, " ")
	if len(arr) != 2 {
		return nil, parseError("wrong input length %d", len(arr))
	}

	op := Op(arr[0])

	switch op {
	case OpAdd, OpSub, OpMul, OpDiv:
	default:
		return nil, parseError("undefined operator %v", op)
	}

	num, err := strconv.Atoi(arr[1])
	if err != nil {
		return nil, err
	}

	return &Calc{Operator: op, Operand: num}, nil
}
