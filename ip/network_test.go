package ip

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/looplab/fsm"

	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/locker/lockermock"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/network/networkmock"
)

func TestSetActivated(t *testing.T) {
	t.Run("It should call EnsureIP", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		networkMock := networkmock.NewMockNetworkInterface(ctrl)

		ip := models.IP{
			IP: "10.0.0.1/32",
			ID: "test-1234",
		}
		networkMock.EXPECT().EnsureIP(ip.IP).Return(nil)

		manager := &manager{
			networkInterface: networkMock,
			ip:               ip,
		}

		manager.setActivated(context.Background(), &fsm.Event{})
	})
}

func TestSetStandBy(t *testing.T) {
	t.Run("It should call RemoveIP", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ip := models.IP{
			IP: "10.0.0.1/32",
			ID: "test-1234",
		}

		networkMock := networkmock.NewMockNetworkInterface(ctrl)
		networkMock.EXPECT().RemoveIP(ip.IP).Return(nil)

		manager := &manager{
			networkInterface: networkMock,
			ip:               ip,
		}

		manager.setStandBy(context.Background(), &fsm.Event{})
	})
}

func TestSetFailing(t *testing.T) {
	t.Run("It remove the IP and stop the locker", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		networkMock := networkmock.NewMockNetworkInterface(ctrl)
		lockerMock := lockermock.NewMockLocker(ctrl)

		ip := models.IP{
			IP: "10.0.0.1/32",
			ID: "test-1234",
		}
		networkMock.EXPECT().RemoveIP(ip.IP).Return(nil)
		lockerMock.EXPECT().Unlock(gomock.Any()).Return(nil)

		manager := &manager{
			networkInterface: networkMock,
			locker:           lockerMock,
			ip:               ip,
		}

		manager.setFailing(context.Background(), &fsm.Event{})
	})
}

func TestStartARPEnsure(t *testing.T) {
	ip := models.IP{
		IP: "10.0.0.1/32",
		ID: "test-1234",
	}

	config := config.Config{
		ARPGratuitousInterval: 10 * time.Millisecond,
		ARPGratuitousCount:    3,
	}

	t.Run("If the IP is activated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		networkMock := networkmock.NewMockNetworkInterface(ctrl)
		networkMock.EXPECT().EnsureIP(ip.IP).Return(nil).MaxTimes(config.ARPGratuitousCount)

		sm := NewStateMachine(ctx, NewStateMachineOpts{})
		sm.SetState(ACTIVATED)
		manager := &manager{
			networkInterface: networkMock,
			stateMachine:     sm,
			config:           config,
			ip:               ip,
		}

		doneChan := make(chan bool)
		go func() {
			manager.startArpEnsure(ctx)
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

	t.Run("If the IP is not activated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		networkMock := networkmock.NewMockNetworkInterface(ctrl)
		sm := NewStateMachine(ctx, NewStateMachineOpts{})
		sm.SetState(FAILING)
		manager := &manager{
			networkInterface: networkMock,
			stateMachine:     sm,
			config:           config,
			ip:               ip,
		}

		doneChan := make(chan bool)
		go func() {
			manager.startArpEnsure(ctx)
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
