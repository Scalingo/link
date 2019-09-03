package main

import (
	"context"
	"fmt"

	"github.com/Scalingo/link/ip"
	"github.com/looplab/fsm"
)

func main() {
	machine := ip.NewStateMachine(context.Background(), ip.NewStateMachineOpts{})

	fmt.Println(fsm.Visualize(machine))
}
