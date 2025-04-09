package ip

import (
	"context"

	"github.com/looplab/fsm"
)

const (
	ACTIVATED = "ACTIVATED"
	STANDBY   = "STANDBY"
	FAILING   = "FAILING"
	BOOTING   = "BOOTING"
)

const (
	FaultEvent              = "fault"
	ElectedEvent            = "elected"
	DemotedEvent            = "demoted"
	HealthCheckFailEvent    = "health_check_fail"
	HealthCheckSuccessEvent = "health_check_success"
)

type NewStateMachineOpts struct {
	ActivatedCallback func(ctx context.Context, e *fsm.Event)
	StandbyCallback   func(ctx context.Context, e *fsm.Event)
	FailingCallback   func(ctx context.Context, e *fsm.Event)
}

func NewStateMachine(_ context.Context, opts NewStateMachineOpts) *fsm.FSM {
	callbacks := fsm.Callbacks{}

	if opts.ActivatedCallback != nil {
		callbacks["enter_"+ACTIVATED] = func(ctx context.Context, e *fsm.Event) {
			opts.ActivatedCallback(ctx, e)
		}
	}

	if opts.StandbyCallback != nil {
		callbacks["enter_"+STANDBY] = func(ctx context.Context, e *fsm.Event) {
			opts.StandbyCallback(ctx, e)
		}
	}

	if opts.FailingCallback != nil {
		callbacks["enter_"+FAILING] = func(ctx context.Context, e *fsm.Event) {
			opts.FailingCallback(ctx, e)
		}
	}

	return fsm.NewFSM(
		BOOTING,
		fsm.Events{
			{Name: FaultEvent, Src: []string{ACTIVATED}, Dst: ACTIVATED},
			{Name: FaultEvent, Src: []string{STANDBY}, Dst: ACTIVATED},
			{Name: FaultEvent, Src: []string{FAILING}, Dst: FAILING},
			{Name: FaultEvent, Src: []string{BOOTING}, Dst: ACTIVATED},
			{Name: ElectedEvent, Src: []string{STANDBY}, Dst: ACTIVATED},
			{Name: ElectedEvent, Src: []string{ACTIVATED}, Dst: ACTIVATED},
			{Name: ElectedEvent, Src: []string{FAILING}, Dst: FAILING},
			{Name: ElectedEvent, Src: []string{BOOTING}, Dst: ACTIVATED},
			{Name: DemotedEvent, Src: []string{ACTIVATED}, Dst: STANDBY},
			{Name: DemotedEvent, Src: []string{STANDBY}, Dst: STANDBY},
			{Name: DemotedEvent, Src: []string{FAILING}, Dst: FAILING},
			{Name: DemotedEvent, Src: []string{BOOTING}, Dst: STANDBY},
			{Name: HealthCheckFailEvent, Src: []string{ACTIVATED}, Dst: FAILING},
			{Name: HealthCheckFailEvent, Src: []string{STANDBY}, Dst: FAILING},
			{Name: HealthCheckFailEvent, Src: []string{FAILING}, Dst: FAILING},
			{Name: HealthCheckFailEvent, Src: []string{BOOTING}, Dst: FAILING},
			{Name: HealthCheckSuccessEvent, Src: []string{FAILING}, Dst: STANDBY},
			{Name: HealthCheckSuccessEvent, Src: []string{STANDBY}, Dst: STANDBY},
			{Name: HealthCheckSuccessEvent, Src: []string{ACTIVATED}, Dst: ACTIVATED},
			{Name: HealthCheckSuccessEvent, Src: []string{BOOTING}, Dst: BOOTING},
		},
		callbacks,
	)
}
