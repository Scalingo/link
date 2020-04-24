package ip

import (
	"context"
	"testing"
	"time"

	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/locker/lockermock"
	"github.com/golang/mock/gomock"
	"github.com/looplab/fsm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_TryToGetIP(t *testing.T) {
	examples := []struct {
		Name           string
		Locker         func(*lockermock.MockLocker)
		ExpectedEvents []string
		KeepAliveRetry int
		CurrentState   string
	}{
		{
			Name: "When refresh fails, fault event",
			Locker: func(mock *lockermock.MockLocker) {
				mock.EXPECT().Refresh(gomock.Any()).Return(errors.New("NOP"))
			},
			KeepAliveRetry: 0,
			ExpectedEvents: []string{FaultEvent},
			CurrentState:   STANDBY,
		}, {
			Name: "When refresh fails with retry, no fault",
			Locker: func(mock *lockermock.MockLocker) {
				mock.EXPECT().Refresh(gomock.Any()).Return(errors.New("NOP"))
			},
			KeepAliveRetry: 1,
			ExpectedEvents: []string{},
			CurrentState:   STANDBY,
		}, {
			Name: "When IsMaster fails just one time, no fault",
			Locker: func(mock *lockermock.MockLocker) {
				mock.EXPECT().Refresh(gomock.Any()).Return(nil)
				mock.EXPECT().IsMaster(gomock.Any()).Return(false, errors.New("NOP"))
			},
			ExpectedEvents: []string{},
			CurrentState:   STANDBY,
		}, {
			Name: "When we are not master",
			Locker: func(mock *lockermock.MockLocker) {
				mock.EXPECT().Refresh(gomock.Any()).Return(nil)
				mock.EXPECT().IsMaster(gomock.Any()).Return(false, nil)
			},
			ExpectedEvents: []string{DemotedEvent},
			CurrentState:   STANDBY,
		}, {
			Name: "When we are master",
			Locker: func(mock *lockermock.MockLocker) {
				mock.EXPECT().Refresh(gomock.Any()).Return(nil)
				mock.EXPECT().IsMaster(gomock.Any()).Return(true, nil)
			},
			ExpectedEvents: []string{ElectedEvent},
			CurrentState:   STANDBY,
		}, {
			Name:           "When the current fsm state is FAILING it should not do anything",
			ExpectedEvents: []string{},
			CurrentState:   FAILING,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			locker := lockermock.NewMockLocker(ctrl)
			if example.Locker != nil {
				example.Locker(locker)
			}

			cfg, err := config.Build()
			require.NoError(t, err)

			cfg.KeepAliveRetry = example.KeepAliveRetry

			manager := &manager{
				locker:       locker,
				config:       cfg,
				stateMachine: fsm.NewFSM(example.CurrentState, fsm.Events{}, fsm.Callbacks{}),
			}

			eventChan := make(chan string, 10)
			doneChan := make(chan bool)
			manager.eventChan = eventChan
			go func() {
				manager.tryToGetIP(ctx)
				// Wait for the eventChan to be processed
				time.Sleep(100 * time.Millisecond)
				doneChan <- true
			}()
			timer := time.NewTimer(500 * time.Millisecond)
			var events []string

			cont := true
			for cont {
				select {
				case <-timer.C:
					t.Fatal("Method did not return")
					break
				case event := <-eventChan:
					events = append(events, event)
				case <-doneChan:
					cont = false
				}
			}

			require.Equal(t, len(example.ExpectedEvents), len(events))

			for i := 0; i < len(example.ExpectedEvents); i++ {
				assert.Equal(t, example.ExpectedEvents[i], events[i])
			}
		})
	}
}
