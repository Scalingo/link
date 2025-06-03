package ip

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/looplab/fsm"

	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/locker/lockermock"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/plugin/pluginmock"
)

func TestSetActivated(t *testing.T) {
	t.Run("It should call the plugin Activated Method", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pluginMock := pluginmock.NewMockPlugin(ctrl)

		endpoint := models.Endpoint{
			ID: "test-1234",
		}

		pluginMock.EXPECT().Activate(gomock.Any()).Return(nil)

		manager := &EndpointManager{
			plugin:   pluginMock,
			endpoint: endpoint,
		}

		manager.setActivated(context.Background(), &fsm.Event{})
	})
}

func TestSetStandBy(t *testing.T) {
	t.Run("It should call the plugin Deactivate method", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := models.Endpoint{
			ID: "test-1234",
		}

		pluginMock := pluginmock.NewMockPlugin(ctrl)
		pluginMock.EXPECT().Deactivate(gomock.Any()).Return(nil)

		manager := &EndpointManager{
			plugin:   pluginMock,
			endpoint: endpoint,
		}

		manager.setStandBy(context.Background(), &fsm.Event{})
	})
}

func TestSetFailing(t *testing.T) {
	t.Run("It remove the endpoint and stop the locker", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lockerMock := lockermock.NewMockLocker(ctrl)

		endpoint := models.Endpoint{
			ID: "test-1234",
		}

		pluginMock := pluginmock.NewMockPlugin(ctrl)
		pluginMock.EXPECT().Deactivate(gomock.Any()).Return(nil)
		lockerMock.EXPECT().Unlock(gomock.Any()).Return(nil)

		manager := &EndpointManager{
			plugin:   pluginMock,
			locker:   lockerMock,
			endpoint: endpoint,
		}

		manager.setFailing(context.Background(), &fsm.Event{})
	})
}

func Test_startPluginEnsureLoop(t *testing.T) {
	endpoint := models.Endpoint{
		ID: "test-1234",
	}

	config := config.Config{
		PluginEnsureInterval: 10 * time.Millisecond,
	}

	t.Run("If the endpoint is activated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		pluginMock := pluginmock.NewMockPlugin(ctrl)
		pluginMock.EXPECT().Ensure(gomock.Any()).Return(nil).MinTimes(9)

		sm := NewStateMachine(ctx, NewStateMachineOpts{})
		sm.SetState(ACTIVATED)
		manager := &EndpointManager{
			stateMachine: sm,
			config:       config,
			endpoint:     endpoint,
			plugin:       pluginMock,
		}

		doneChan := make(chan bool)
		go func() {
			manager.startPluginEnsureLoop(ctx)
			doneChan <- true
		}()
		time.Sleep(100 * time.Millisecond)

		manager.stopMutex.Lock()
		manager.stopped = true
		manager.stopMutex.Unlock()

		timer := time.NewTimer(500 * time.Millisecond)
		select {
		case <-timer.C:
			t.Fatal("NOT RESPONDING")
		case <-doneChan:
		}
	})

	t.Run("If the endpoint is not activated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		pluginMock := pluginmock.NewMockPlugin(ctrl)
		sm := NewStateMachine(ctx, NewStateMachineOpts{})
		sm.SetState(FAILING)
		manager := &EndpointManager{
			plugin:       pluginMock,
			stateMachine: sm,
			config:       config,
			endpoint:     endpoint,
		}

		doneChan := make(chan bool)
		go func() {
			manager.startPluginEnsureLoop(ctx)
			doneChan <- true
		}()
		time.Sleep(50 * time.Millisecond)

		manager.stopMutex.Lock()
		manager.stopped = true
		manager.stopMutex.Unlock()

		timer := time.NewTimer(500 * time.Millisecond)
		select {
		case <-timer.C:
			t.Fatal("NOT RESPONDING")
		case <-doneChan:
		}
	})
}
