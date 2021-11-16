package main

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"

	"github.com/Scalingo/link/v2/ip"
)

func main() {
	machine := ip.NewStateMachine(context.Background(), ip.NewStateMachineOpts{})

	fmt.Println(fsm.Visualize(machine))
}
