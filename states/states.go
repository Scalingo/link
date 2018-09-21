package states

import (
	"github.com/looplab/fsm"
)

const (
	ACTIVATED = "ACTIVATED"
	STANDBY   = "STANDBY"
)

const (
	FaultEvent = "fault"
	Elected    = "elected"
	Demoted    = "demoted"
)

func NewStateMachine() *fsm.FSM {
	machine := fsm.NewFSM(
		STANDBY,
		fsm.Events{
			{Name: FaultEvent, Src: []string{ACTIVATED}, Dst: ACTIVATED},
			{Name: FaultEvent, Src: []string{STANDBY}, Dst: ACTIVATED},
			{Name: Elected, Src: []string{STANDBY}, Dst: ACTIVATED},
			{Name: Elected, Src: []string{ACTIVATED}, Dst: ACTIVATED},
			{Name: Demoted, Src: []string{ACTIVATED}, Dst: STANDBY},
		},
		fsm.Callbacks{
			"enter_" + ACTIVATED: func(e *fsm.Event) {
			},
			"enter_" + STANDBY: func(e *fsm.Event) {
			},
		},
	)
	return machine
}
