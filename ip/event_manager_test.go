package ip

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/looplab/fsm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/locker/lockermock"
	"github.com/Scalingo/link/v2/models/modelsmock"
	"github.com/Scalingo/link/v2/plugin/pluginmock"
	"github.com/Scalingo/link/v2/watcher/watchermock"
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

			cfg, err := config.Build(ctx)
			require.NoError(t, err)

			cfg.KeepAliveRetry = example.KeepAliveRetry

			manager := &EndpointManager{
				locker:       locker,
				config:       cfg,
				stateMachine: fsm.NewFSM(example.CurrentState, fsm.Events{}, fsm.Callbacks{}),
			}

			eventChan := make(chan string, 10)
			doneChan := make(chan bool)
			manager.eventChan = eventChan
			go func() {
				manager.tryToGetEndpoint(ctx)
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

func TestManager_Stop(t *testing.T) {
	examples := []struct {
		Name                      string
		Locker                    func(*lockermock.MockLocker)
		HostCount                 int
		ShouldWaitForReallocation bool
		CurrentState              string
		Events                    []string
	}{
		{
			Name: "When there are no other hosts",
			Locker: func(l *lockermock.MockLocker) {
				l.EXPECT().IsMaster(gomock.Any()).Return(true, nil)
			},
			HostCount:                 1,
			ShouldWaitForReallocation: false,
			CurrentState:              ACTIVATED,
			Events:                    []string{DemotedEvent},
		}, {
			Name: "When there are no other host and we are not failing",
			Locker: func(l *lockermock.MockLocker) {
				l.EXPECT().IsMaster(gomock.Any()).Return(true, nil)
			},
			HostCount:                 1,
			ShouldWaitForReallocation: false,
			CurrentState:              FAILING,
			Events:                    []string{},
		}, {
			Name: "When there are other hosts trying to take the lock",
			Locker: func(l *lockermock.MockLocker) {
				l.EXPECT().IsMaster(gomock.Any()).Return(true, nil)
				l.EXPECT().IsMaster(gomock.Any()).Return(false, nil)
			},
			HostCount:                 2,
			ShouldWaitForReallocation: true,
			CurrentState:              STANDBY,
			Events:                    []string{DemotedEvent},
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			lockerMock := lockermock.NewMockLocker(ctrl)
			storageMock := modelsmock.NewMockStorage(ctrl)
			watcherMock := watchermock.NewMockWatcher(ctrl)
			pluginMock := pluginmock.NewMockPlugin(ctrl)

			pluginMock.EXPECT().ElectionKey(gomock.Any()).Return("test-election-key").AnyTimes()

			hosts := make([]string, example.HostCount)
			storageMock.EXPECT().GetEndpointHosts(gomock.Any(), "test-election-key").Return(hosts, nil)
			if example.Locker != nil {
				example.Locker(lockerMock)
			}
			watcherMock.EXPECT().Stop(gomock.Any())
			lockerMock.EXPECT().Stop(gomock.Any()).Return(nil)
			storageMock.EXPECT().UnlinkEndpointFromCurrentHost(gomock.Any(), "test-election-key").Return(nil)
			eventChan := make(chan string, 2)
			events := make([]string, 0)

			manager := &EndpointManager{
				stateMachine: fsm.NewFSM(example.CurrentState, fsm.Events{}, fsm.Callbacks{}),
				locker:       lockerMock,
				storage:      storageMock,
				watcher:      watcherMock,
				eventChan:    eventChan,
				plugin:       pluginMock,
			}

			err := manager.Stop(context.Background())
			require.NoError(t, err)

			timer := time.NewTimer(1 * time.Second)
			stop := false
			for !stop {
				select {
				case <-timer.C:
					t.Fatal("eventChan was never closed")
					break
				case event, ok := <-eventChan:
					if !ok {
						stop = true
						break
					}
					events = append(events, event)
				}
			}

			assert.True(t, manager.stopped)
			assert.Equal(t, example.Events, events)
		})
	}
}
