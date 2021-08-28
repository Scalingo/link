package ip

import (
	"context"
	"testing"
	"time"

	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/locker/lockermock"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/network/networkmock"
	"github.com/golang/mock/gomock"
	"github.com/looplab/fsm"
)

func TestSetActivated(t *testing.T) {
	ip := models.IP{
		IP: "10.0.0.1/32",
		ID: "test-1234",
	}

	examples := map[string]struct {
		noNetwork              bool
		expectNetworkInterface func(mock *networkmock.MockNetworkInterface)
	}{
		"when no_netwok is set to false, add the IP": {
			noNetwork: false,
			expectNetworkInterface: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().EnsureIP(ip.IP).Return(nil)
			},
		},
		"when no_network is set to true, do nothing": {
			noNetwork: true,
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			networkMock := networkmock.NewMockNetworkInterface(ctrl)

			if example.expectNetworkInterface != nil {
				example.expectNetworkInterface(networkMock)
			}

			ip.NoNetwork = example.noNetwork

			manager := &manager{
				networkInterface: networkMock,
				ip:               ip,
			}

			manager.setActivated(context.Background(), &fsm.Event{})
		})
	}
}

func TestSetStandBy(t *testing.T) {
	ip := models.IP{
		IP: "10.0.0.1/32",
		ID: "test-1234",
	}

	examples := map[string]struct {
		noNetwork              bool
		expectNetworkInterface func(mock *networkmock.MockNetworkInterface)
	}{
		"when no_netwok is set to false, remove the IP": {
			noNetwork: false,
			expectNetworkInterface: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().RemoveIP(ip.IP).Return(nil)
			},
		},
		"when no_network is set to true, do nothing": {
			noNetwork: true,
		},
	}
	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			networkMock := networkmock.NewMockNetworkInterface(ctrl)

			if example.expectNetworkInterface != nil {
				example.expectNetworkInterface(networkMock)
			}

			ip.NoNetwork = example.noNetwork

			manager := &manager{
				networkInterface: networkMock,
				ip:               ip,
			}

			manager.setStandBy(context.Background(), &fsm.Event{})
		})
	}
}

func TestSetFailing(t *testing.T) {
	ip := models.IP{
		IP: "10.0.0.1/32",
		ID: "test-1234",
	}

	examples := map[string]struct {
		noNetwork              bool
		expectNetworkInterface func(mock *networkmock.MockNetworkInterface)
	}{
		"when no_netwok is set to false, remove the IP and stop the locker": {
			noNetwork: false,
			expectNetworkInterface: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().RemoveIP(ip.IP).Return(nil)
			},
		},
		"when no_network is set to true, only stop the locker": {
			noNetwork: true,
		},
	}
	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			networkMock := networkmock.NewMockNetworkInterface(ctrl)
			lockerMock := lockermock.NewMockLocker(ctrl)
			lockerMock.EXPECT().Unlock(gomock.Any()).Return(nil)

			if example.expectNetworkInterface != nil {
				example.expectNetworkInterface(networkMock)
			}

			ip.NoNetwork = example.noNetwork

			manager := &manager{
				networkInterface: networkMock,
				locker:           lockerMock,
				ip:               ip,
			}

			manager.setFailing(context.Background(), &fsm.Event{})
		})
	}

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

	examples := map[string]struct {
		state                  string
		noNetwork              bool
		expectNetworkInterface func(mock *networkmock.MockNetworkInterface)
	}{
		"if the IP is activated": {
			state: ACTIVATED,
			expectNetworkInterface: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().EnsureIP(ip.IP).Return(nil).MaxTimes(config.ARPGratuitousCount)
			},
		},
		"if the IP is activated but no_network is set to true": {
			state:     ACTIVATED,
			noNetwork: true,
		},
		"if the IP is not activated": {
			state: STANDBY,
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			networkMock := networkmock.NewMockNetworkInterface(ctrl)
			if example.expectNetworkInterface != nil {
				example.expectNetworkInterface(networkMock)
			}

			sm := NewStateMachine(ctx, NewStateMachineOpts{})
			sm.SetState(example.state)

			ip.NoNetwork = example.noNetwork

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
	}
}
