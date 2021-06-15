package ip

import (
	"context"
	"testing"
	"time"

	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/healthcheck/healthcheckmock"
	"github.com/golang/mock/gomock"
	"github.com/looplab/fsm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestManager_HealthChecker(t *testing.T) {
	examples := []struct {
		Name           string
		Checker        func(*healthcheckmock.MockChecker)
		ExpectedEvents []string
		CurrentState   string
	}{
		{
			Name: "With not enough failing events",
			Checker: func(mock *healthcheckmock.MockChecker) {
				mock.EXPECT().IsHealthy(gomock.Any()).Return(false, errors.New("failing")).MaxTimes(2)
				mock.EXPECT().IsHealthy(gomock.Any()).Return(true, nil).AnyTimes()
			},
			ExpectedEvents: []string{HealthCheckSuccessEvent},
			CurrentState:   STANDBY,
		}, {
			Name: "With enough failing events but we're not ACTIVATED",
			Checker: func(mock *healthcheckmock.MockChecker) {
				mock.EXPECT().IsHealthy(gomock.Any()).Return(false, errors.New("failing")).MaxTimes(3)
				mock.EXPECT().IsHealthy(gomock.Any()).Return(true, nil).AnyTimes()
			},
			ExpectedEvents: []string{HealthCheckFailEvent, HealthCheckSuccessEvent},
			CurrentState:   STANDBY,
		}, {
			Name: "With enough failing events",
			Checker: func(mock *healthcheckmock.MockChecker) {
				mock.EXPECT().IsHealthy(gomock.Any()).Return(false, errors.New("failing")).MaxTimes(3)
				mock.EXPECT().IsHealthy(gomock.Any()).Return(true, nil).AnyTimes()
			},
			ExpectedEvents: []string{HealthCheckFailEvent, HealthCheckSuccessEvent},
			CurrentState:   ACTIVATED,
		}, {
			Name: "With a success event and a stop",
			Checker: func(mock *healthcheckmock.MockChecker) {
				mock.EXPECT().IsHealthy(gomock.Any()).Return(true, nil).MaxTimes(2)
			},
			ExpectedEvents: []string{HealthCheckSuccessEvent},
			CurrentState:   STANDBY,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			checker := healthcheckmock.NewMockChecker(ctrl)
			example.Checker(checker)

			manager := &manager{
				checker:      checker,
				stateMachine: fsm.NewFSM(example.CurrentState, fsm.Events{}, fsm.Callbacks{}),
				config: config.Config{
					HealthcheckInterval:     10 * time.Millisecond,
					FailCountBeforeFailover: 3,
					KeepAliveInterval:       10 * time.Millisecond,
				},
			}

			eventChan := make(chan string, 1)
			doneChan := make(chan bool)
			manager.eventChan = eventChan
			go func() {
				manager.healthChecker(ctx)
				doneChan <- true
			}()
			timer := time.NewTimer(500 * time.Millisecond)
			var events []string

			cont := true
			i := 0
			for cont {
				select {
				case <-timer.C:
					t.Fatal("Method did not return")
					break
				case event := <-eventChan:
					events = append(events, event)
					i++
				case <-doneChan:
					cont = false
				}
				if i >= len(example.ExpectedEvents) {
					cont = false
				}
			}

			manager.stopMutex.Lock()
			manager.stopped = true
			manager.stopMutex.Unlock()

			for i := 0; i < len(example.ExpectedEvents); i++ {
				assert.Equal(t, example.ExpectedEvents[i], events[i])
			}
		})
	}
}
