package ip

import (
	"context"

	"github.com/looplab/fsm"
)

const (
	ACTIVATED = "ACTIVATED"
	STANDBY   = "STANDBY"
)

const (
	FaultEvent   = "fault"
	ElectedEvent = "elected"
	DemotedEvent = "demoted"
)

func (m *manager) newStateMachine(ctx context.Context) {
	machine := fsm.NewFSM(
		STANDBY,
		fsm.Events{
			{Name: FaultEvent, Src: []string{ACTIVATED}, Dst: ACTIVATED},
			{Name: FaultEvent, Src: []string{STANDBY}, Dst: ACTIVATED},
			{Name: ElectedEvent, Src: []string{STANDBY}, Dst: ACTIVATED},
			{Name: ElectedEvent, Src: []string{ACTIVATED}, Dst: ACTIVATED},
			{Name: DemotedEvent, Src: []string{ACTIVATED}, Dst: STANDBY},
			{Name: DemotedEvent, Src: []string{STANDBY}, Dst: STANDBY},
		},
		fsm.Callbacks{
			"enter_" + ACTIVATED: func(e *fsm.Event) {
				m.setActivated(ctx)
			},
			"enter_" + STANDBY: func(e *fsm.Event) {
				m.setStandBy(ctx)
			},
		},
	)
	m.stateMachine = machine
}
